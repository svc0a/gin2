package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Store interface {
	Store(key string, value []byte)
	Load(key string) ([]byte, error)
}

// bodyLogWriter 用于捕获响应主体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func IdempotencyMiddleware(store Store, idempotencyKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 中获取 Idempotency-Key
		k := c.GetHeader(idempotencyKey)
		if k == "" {
			// 如果不存在幂等性 Key，则直接继续执行
			c.Next()
			return
		}
		result, err := store.Load(k)
		if err == nil {
			// 如果已经处理过，返回存储的响应
			c.Status(http.StatusOK)
			c.Data(http.StatusOK, "application/json", result)
			c.Abort()
			return
		}
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		// 否则，继续处理请求
		c.Next()
		// 将请求的响应结果存储起来，作为幂等性响应
		if c.Writer.Status() != http.StatusOK {
			return
		}
		responseBody := blw.body.Bytes()
		store.Store(k, responseBody)
	}
}
