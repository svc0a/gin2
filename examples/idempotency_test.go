package examples

import (
	"github.com/gin-gonic/gin"
	"github.com/svc0a/gin2/middleware"
	"net/http"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	r := gin.Default()
	// 使用幂等性中间件
	r.Use(middleware.IdempotencyMiddleware(middleware.NewMemoryStore(), "idempotency-key"))

	// 测试路由
	r.Any("/process", func(c *gin.Context) {
		// 模拟业务处理逻辑
		c.JSON(http.StatusOK, gin.H{
			"status":    "success",
			"message":   "Request processed",
			"timestamp": time.Now().Format(time.DateTime),
		})
	})

	r.Run() // 启动服务
}
