package utils

import "strings"

// CleanHTMLEntities 清理命令中的 HTML 实体编码
// 这个函数用于修复前端可能传递过来的 HTML 实体编码问题
// 例如: &quot; -> ", &apos; -> ', &lt; -> <, &gt; -> >, &amp; -> &
func CleanHTMLEntities(command string) string {
	// 如果命令中不包含 HTML 实体,直接返回
	if !strings.Contains(command, "&") {
		return command
	}

	// 定义 HTML 实体替换映射
	replacements := map[string]string{
		"&quot;": "\"",
		"&apos;": "'",
		"&#39;":  "'",
		"&lt;":   "<",
		"&gt;":   ">",
		"&amp;":  "&", // 注意: &amp; 必须最后替换,避免重复替换
	}

	result := command
	// 先替换除 &amp; 之外的所有实体
	for entity, char := range replacements {
		if entity != "&amp;" {
			result = strings.ReplaceAll(result, entity, char)
		}
	}
	// 最后替换 &amp;
	result = strings.ReplaceAll(result, "&amp;", "&")

	return result
}

// ContainsHTMLEntity 检测命令中是否包含 HTML 实体
func ContainsHTMLEntity(command string) bool {
	entities := []string{"&quot;", "&apos;", "&#39;", "&lt;", "&gt;", "&amp;"}
	for _, entity := range entities {
		if strings.Contains(command, entity) {
			return true
		}
	}
	return false
}
