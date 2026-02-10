package app

import (
	"os"
	"path/filepath"

	"fmt"
	"strconv"
	"strings"

	"github.com/gocronx-team/gocron/internal/modules/logger"
	"github.com/gocronx-team/gocron/internal/modules/setting"
	"github.com/gocronx-team/gocron/internal/modules/utils"
)

var (
	// AppDir 应用根目录
	AppDir string // 应用根目录
	// ConfDir 配置文件目录
	ConfDir string // 配置目录
	// LogDir 日志目录
	LogDir string // 日志目录
	// AppConfig 配置文件
	AppConfig string // 应用配置文件
	// Installed 应用是否已安装
	Installed bool // 应用是否安装过
	// Setting 应用配置
	Setting *setting.Setting // 应用配置
	// VersionId 版本号
	VersionId int // 版本号
	// VersionFile 版本文件
	VersionFile string // 版本号文件
)

// InitEnv 初始化
func InitEnv(versionString string) {
	logger.InitLogger()
	var err error
	// 开发环境使用当前目录，生产环境使用可执行文件目录
	execPath, err := os.Executable()
	if err != nil {
		logger.Fatal(err)
	}
	execDir := filepath.Dir(execPath)

	// 如果可执行文件在 tmp 目录（开发环境），使用项目根目录
	if filepath.Base(execDir) == "tmp" {
		AppDir = filepath.Join(filepath.Dir(execDir), ".gocron")
	} else {
		AppDir = filepath.Join(execDir, ".gocron")
	}
	fmt.Printf("AppDir: %s\n", AppDir)
	ConfDir = filepath.Join(AppDir, "conf")
	LogDir = filepath.Join(AppDir, "log")
	AppConfig = filepath.Join(ConfDir, "app.ini")
	VersionFile = filepath.Join(ConfDir, ".version")
	fmt.Printf("ConfDir: %s, LogDir: %s\n", ConfDir, LogDir)
	createDirIfNotExists(AppDir, ConfDir, LogDir)
	Installed = IsInstalled()
	VersionId = ToNumberVersion(versionString)
}

// IsInstalled 判断应用是否已安装
func IsInstalled() bool {
	// 如果配置了 DB 引擎环境变量，视为已安装（自动配置模式）
	if os.Getenv("GOCRON_DB_ENGINE") != "" {
		return true
	}
	
	// 检查 install.lock
	_, err := os.Stat(filepath.Join(ConfDir, "install.lock"))
	if !os.IsNotExist(err) {
		return true
	}

	// 增强检查：如果 app.ini 存在，尝试读取其中的 DB 配置
	if utils.FileExist(AppConfig) {
		// 这里我们不能直接调用 setting.Read，因为 setting 包可能还未完全准备好（虽然它没有外部依赖），
		// 而且我们不想在这里引入太复杂的逻辑。
		// 简单的检查：如果 app.ini 存在，我们假设用户可能手动配置过。
		// 但为了保险，我们可以尝试读取它看看是否有 sqlite 的配置且文件存在。
		// 由于 setting.Read 是安全的，我们可以尝试调用它。
		// 注意：setting 包的 Read 函数现在会从环境变量覆盖，但我们这里只关心文件配置。
		// 实际上，如果 AppConfig 存在，通常意味着已经安装过。
		// 唯一的例外是用户拷贝了一个空的 app.ini 模板。
		// 让我们检查一下 AppConfig 是否有效。
		cfg, err := setting.Read(AppConfig)
		if err == nil {
			// 如果是 SQLite，检查 DB 文件是否存在
			if cfg.Db.Engine == "sqlite" {
				dbPath := cfg.Db.Database
				if !filepath.IsAbs(dbPath) {
					// 如果是相对路径，它是相对于执行目录的（通常是 AppDir 或当前目录）
					// 在 setting.go 中没有处理相对路径的基准，通常由 GORM 处理。
					// GORM/SQLite 驱动通常相对于当前工作目录。
					// 但这里我们简单检查一下常见位置。
					if utils.FileExist(dbPath) || utils.FileExist(filepath.Join(AppDir, dbPath)) {
						return true
					}
				} else {
					if utils.FileExist(dbPath) {
						return true
					}
				}
			} else {
				// 对于 MySQL/Postgres，如果配置文件存在，我们假设它是配置好的
				// 因为我们很难在不连接的情况下验证。
				return true
			}
		}
	}

	// 最后的兜底：如果默认的 sqlite 文件存在（./data/gocron.db 或 ./gocron.db），
	// 且用户没有明确配置 app.ini，我们是否应该视为已安装？
	// 用户的需求是 "已经存在 ./data/gocron.db ... 却进入初始化界面"
	// 这说明用户可能期望自动发现。
	possibleDbPaths := []string{
		"gocron.db",
		"data/gocron.db",
		filepath.Join(AppDir, "gocron.db"),
		filepath.Join(AppDir, "data/gocron.db"),
	}

	for _, path := range possibleDbPaths {
		if utils.FileExist(path) {
			// 只有当没有 app.ini 时，我们才敢假设这个 DB 是用户想要的
			if !utils.FileExist(AppConfig) {
				// 我们需要一种方式告诉系统使用这个 DB。
				// 但 IsInstalled 只是返回 bool。
				// 如果我们返回 true，后续 setting.Read 会读取默认值 (sqlite, gocron.db)。
				// 如果实际 DB 在 data/gocron.db，而默认是 gocron.db，那还是连不上（会创建新的空 DB）。
				// 所以这里不能简单返回 true，除非我们能修改配置。
				// 但 app 包不应该修改 setting。
				
				// 这种情况下，最好的办法是告诉用户：请配置环境变量或 app.ini 指向你的 DB。
				// 或者，我们修改 setting.go 的默认逻辑来寻找 DB？不，那太隐晦了。
				
				// 回到用户的场景：用户说 "已经存在 ./data/gocron.db"。
				// 如果用户是通过 docker 挂载进来的，或者手动放的。
				// 如果没有 install.lock，系统认为未安装。
				// 如果我们在这里返回 true，系统会用默认配置 (gocron.db) 启动，
				// 结果是：系统启动了，但看不到原来的数据（因为连的是新创建的 gocron.db），
				// 或者如果默认就是 gocron.db 且文件就在那，那就完美了。
				
				// 但如果文件在 ./data/gocron.db，而默认配置是 ./gocron.db，
				// 返回 true 会导致连接错误位置。
				
				// 除非我们能检测到 ./data/gocron.db 并通过某种方式传递给 setting。
				// 比如设置环境变量？
				if path == "data/gocron.db" || strings.HasSuffix(path, "data/gocron.db") {
					os.Setenv("GOCRON_DB_DATABASE", path)
					return true
				}
				if path == "gocron.db" || strings.HasSuffix(path, "gocron.db") {
					return true
				}
			}
		}
	}

	return false
}

// CreateInstallLock 创建安装锁文件
func CreateInstallLock() error {
	lockFile := filepath.Join(ConfDir, "install.lock")
	err := os.WriteFile(lockFile, []byte(""), 0600)
	if err != nil {
		logger.Error("创建安装锁文件conf/install.lock失败", err)
		fmt.Printf("Error creating install.lock: %v\n", err)
	} else {
		fmt.Printf("Successfully created install.lock at %s\n", lockFile)
	}

	return err
}

// UpdateVersionFile 更新应用版本号文件
func UpdateVersionFile() {
	err := os.WriteFile(VersionFile,
		[]byte(strconv.Itoa(VersionId)),
		0600,
	)

	if err != nil {
		logger.Fatal(err)
	}
}

// GetCurrentVersionId 获取应用当前版本号, 从版本号文件中读取
func GetCurrentVersionId() int {
	if !utils.FileExist(VersionFile) {
		return 0
	}

	bytes, err := os.ReadFile(VersionFile)
	if err != nil {
		logger.Fatal(err)
	}

	versionId, err := strconv.Atoi(strings.TrimSpace(string(bytes)))
	if err != nil {
		logger.Fatal(err)
	}

	return versionId
}

// ToNumberVersion 把字符串版本号a.b.c转换为整数版本号abc
func ToNumberVersion(versionString string) int {
	versionString = strings.TrimPrefix(versionString, "v")
	v := strings.Replace(versionString, ".", "", -1)
	if len(v) < 3 {
		v += "0"
	}

	versionId, err := strconv.Atoi(v)
	if err != nil {
		logger.Fatal(err)
	}

	return versionId
}

// 检测目录是否存在
func createDirIfNotExists(path ...string) {
	for _, value := range path {
		if utils.FileExist(value) {
			continue
		}
		err := os.MkdirAll(value, 0755)
		if err != nil {
			logger.Fatal(fmt.Sprintf("创建目录失败:%s", err.Error()))
		}
	}
}
