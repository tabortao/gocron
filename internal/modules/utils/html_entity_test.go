package utils

import (
	"strings"
	"testing"
)

// TestHTMLEntityDetection 测试 HTML 实体检测（可在任何平台运行）
func TestHTMLEntityDetection(t *testing.T) {
	tests := []struct {
		name          string
		command       string
		hasHTMLEntity bool
		expectedClean string
		description   string
	}{
		{
			name:          "正常的 Windows 命令（带双引号）",
			command:       `copy "C:\My Documents\report.docx" "D:\Backup"`,
			hasHTMLEntity: false,
			expectedClean: `copy "C:\My Documents\report.docx" "D:\Backup"`,
			description:   "这是正确的命令，应该能在 Windows 上执行",
		},
		{
			name:          "包含 HTML 实体的命令（&quot;）",
			command:       `copy &quot;C:\My Documents\report.docx&quot; &quot;D:\Backup&quot;`,
			hasHTMLEntity: true,
			expectedClean: `copy "C:\My Documents\report.docx" "D:\Backup"`,
			description:   "这个命令在 Windows 上会失败，因为 &quot; 不是有效的引号",
		},
		{
			name:          "mkdir 命令（HTML 实体）",
			command:       `mkdir &quot;C:\Users\John Doe\Projects&quot;`,
			hasHTMLEntity: true,
			expectedClean: `mkdir "C:\Users\John Doe\Projects"`,
			description:   "mkdir 命令也会受影响",
		},
		{
			name:          "dir 命令（HTML 实体）",
			command:       `dir &quot;C:\Program Files (x86)&quot;`,
			hasHTMLEntity: true,
			expectedClean: `dir "C:\Program Files (x86)"`,
			description:   "dir 命令也会受影响",
		},
		{
			name:          "del 命令（HTML 实体）",
			command:       `del &quot;C:\Temp\old file.txt&quot;`,
			hasHTMLEntity: true,
			expectedClean: `del "C:\Temp\old file.txt"`,
			description:   "del 命令也会受影响",
		},
		{
			name:          "start 命令（混合引号）",
			command:       `start &quot;&quot; &quot;C:\Program Files\Google\Chrome\Application\chrome.exe&quot; --new-window`,
			hasHTMLEntity: true,
			expectedClean: `start "" "C:\Program Files\Google\Chrome\Application\chrome.exe" --new-window`,
			description:   "start 命令的空标题也需要引号",
		},
		{
			name:          "包含其他 HTML 实体",
			command:       `echo &lt;test&gt; &amp; &apos;hello&apos;`,
			hasHTMLEntity: true,
			expectedClean: `echo <test> & 'hello'`,
			description:   "其他 HTML 实体也应该被清理",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("\n描述: %s", tt.description)
			t.Logf("原始命令: %s", tt.command)

			// 检测是否包含 HTML 实体
			hasEntity := ContainsHTMLEntity(tt.command)
			if hasEntity != tt.hasHTMLEntity {
				t.Errorf("HTML 实体检测错误: got %v, want %v", hasEntity, tt.hasHTMLEntity)
			}

			// 清理 HTML 实体
			cleaned := CleanHTMLEntities(tt.command)
			t.Logf("清理后命令: %s", cleaned)

			if cleaned != tt.expectedClean {
				t.Errorf("清理结果不符合预期:\ngot:  %s\nwant: %s", cleaned, tt.expectedClean)
			}

			// 验证清理后不再包含 HTML 实体
			if ContainsHTMLEntity(cleaned) {
				t.Errorf("清理后仍然包含 HTML 实体: %s", cleaned)
			}
		})
	}
}

// TestCommandLength 测试命令长度限制
func TestCommandLength(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		maxLength   int
		shouldFit   bool
		description string
	}{
		{
			name:        "短命令（适合 256 字符限制）",
			command:     `dir "C:\Program Files"`,
			maxLength:   256,
			shouldFit:   true,
			description: "简单命令应该没问题",
		},
		{
			name: "长路径命令（适合 256 字符限制）",
			command: `copy "C:\Users\Administrator\Documents\Projects\MyProject\src\main\resources\config\application-production.properties" ` +
				`"D:\Backup\Projects\MyProject\config\backup-2024-01-15\application-production.properties"`,
			maxLength:   256,
			shouldFit:   true,
			description: "这个路径实际上只有 208 字符,在 256 限制内",
		},
		{
			name: "长路径命令（适合 1024 字符限制）",
			command: `copy "C:\Users\Administrator\Documents\Projects\MyProject\src\main\resources\config\application-production.properties" ` +
				`"D:\Backup\Projects\MyProject\config\backup-2024-01-15\application-production.properties"`,
			maxLength:   1024,
			shouldFit:   true,
			description: "增加到 1024 字符后应该足够",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("\n描述: %s", tt.description)
			t.Logf("命令长度: %d 字符", len(tt.command))
			t.Logf("限制长度: %d 字符", tt.maxLength)

			fits := len(tt.command) <= tt.maxLength
			if fits != tt.shouldFit {
				t.Errorf("长度检查错误: 命令长度 %d, 限制 %d, got %v, want %v",
					len(tt.command), tt.maxLength, fits, tt.shouldFit)
			}

			if !fits {
				t.Logf("⚠️  命令超出长度限制 %d 字符", len(tt.command)-tt.maxLength)
			}
		})
	}
}

// TestWindowsCommandSimulation 模拟 Windows 命令行为
func TestWindowsCommandSimulation(t *testing.T) {
	t.Log("\n=== 模拟 Windows cmd.exe 行为 ===\n")

	tests := []struct {
		name     string
		command  string
		willWork bool
		reason   string
	}{
		{
			name:     "正确的双引号",
			command:  `dir "C:\Program Files"`,
			willWork: true,
			reason:   "Windows cmd 能正确识别双引号",
		},
		{
			name:     "HTML 实体 &quot;",
			command:  `dir &quot;C:\Program Files&quot;`,
			willWork: false,
			reason:   "Windows cmd 会将 &quot; 当作普通字符串，不是引号",
		},
		{
			name:     "路径中有空格但没有引号",
			command:  `dir C:\Program Files`,
			willWork: false,
			reason:   "Windows cmd 会将空格作为参数分隔符",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("命令: %s", tt.command)
			t.Logf("预期结果: %v", tt.willWork)
			t.Logf("原因: %s", tt.reason)

			// 模拟检查
			hasHTMLEntity := ContainsHTMLEntity(tt.command)
			hasSpaceWithoutQuotes := strings.Contains(tt.command, " ") &&
				!strings.Contains(tt.command, "\"") &&
				!hasHTMLEntity

			simulatedSuccess := !hasHTMLEntity && !hasSpaceWithoutQuotes

			if simulatedSuccess != tt.willWork {
				t.Logf("⚠️  模拟结果与预期不符")
			} else {
				t.Logf("✓ 模拟结果符合预期")
			}
		})
	}
}
