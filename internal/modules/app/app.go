package app

import (
	"os"
	"path/filepath"

	"fmt"
	"strconv"
	"strings"

	"github.com/tabortao/gocron/internal/modules/logger"
	"github.com/tabortao/gocron/internal/modules/setting"
	"github.com/tabortao/gocron/internal/modules/utils"
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

	setSqliteEnv := func(dbPath string) {
		if os.Getenv("GOCRON_DB_ENGINE") == "" {
			os.Setenv("GOCRON_DB_ENGINE", "sqlite")
		}
		if dbPath != "" {
			os.Setenv("GOCRON_DB_DATABASE", dbPath)
		}
	}

	if dbEnv := os.Getenv("GOCRON_DB_DATABASE"); dbEnv != "" {
		if utils.FileExist(dbEnv) && (strings.HasSuffix(dbEnv, ".db") || strings.HasSuffix(dbEnv, ".sqlite") || strings.HasSuffix(dbEnv, ".sqlite3")) {
			setSqliteEnv(dbEnv)
			return true
		}
	}

	// 检查 install.lock
	_, err := os.Stat(filepath.Join(ConfDir, "install.lock"))
	if !os.IsNotExist(err) {
		return true
	}

	detectExistingSqliteDB := func() (string, bool) {
		possibleDbPaths := []string{
			filepath.Join(AppDir, "data", "gocron.db"),
			filepath.Join(AppDir, "gocron.db"),
			"data/gocron.db",
			"gocron.db",
		}

		for _, path := range possibleDbPaths {
			if utils.FileExist(path) {
				return path, true
			}
		}
		return "", false
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

				if detected, ok := detectExistingSqliteDB(); ok {
					setSqliteEnv(detected)
					return true
				}
			} else {
				// 对于 MySQL/Postgres，如果配置文件存在，我们假设它是配置好的
				// 因为我们很难在不连接的情况下验证。
				return true
			}
		}
	}

	if detected, ok := detectExistingSqliteDB(); ok {
		setSqliteEnv(detected)
		return true
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
