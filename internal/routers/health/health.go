package health

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tabortao/gocron/internal/models"
	"github.com/tabortao/gocron/internal/modules/app"
	"github.com/tabortao/gocron/internal/modules/rpc/grpcpool"
	"github.com/tabortao/gocron/internal/modules/utils"
	"github.com/tabortao/gocron/internal/service"
)

func Healthz(c *gin.Context) {
	ok := true
	data := map[string]interface{}{
		"installed": app.Installed,
		"now":       time.Now().Format(time.RFC3339),
		"runtime": map[string]interface{}{
			"goroutines": runtime.NumGoroutine(),
		},
		"rpc": map[string]interface{}{
			"poolSize": grpcpool.Pool.Size(),
		},
	}

	dbInfo := map[string]interface{}{
		"ok": true,
	}
	if app.Installed {
		if models.Db == nil {
			dbInfo["ok"] = false
			dbInfo["error"] = "db is nil"
			ok = false
		} else {
			sqlDB, err := models.Db.DB()
			if err != nil {
				dbInfo["ok"] = false
				dbInfo["error"] = err.Error()
				ok = false
			} else {
				ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
				defer cancel()
				if err := sqlDB.PingContext(ctx); err != nil {
					dbInfo["ok"] = false
					dbInfo["error"] = err.Error()
					ok = false
				}
			}
		}
	}
	data["db"] = dbInfo

	scheduler := service.GetSchedulerStatus()
	data["scheduler"] = scheduler
	if app.Installed && !scheduler.Running {
		ok = false
	}

	jsonResp := utils.JsonResponse{}
	if ok {
		c.String(http.StatusOK, jsonResp.Success("ok", data))
		return
	}
	c.String(http.StatusServiceUnavailable, jsonResp.Success("unhealthy", data))
}
