package base

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gocronx-team/gocron/internal/models"
)

// ParsePageAndPageSize 解析查询参数中的页数和每页数量
func ParsePageAndPageSize(c *gin.Context, params models.CommonMap) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = models.PageSize
	}

	params["Page"] = page
	params["PageSize"] = pageSize
}
