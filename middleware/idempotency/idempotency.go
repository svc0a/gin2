package idempotency

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Response struct {
	header map[string]string
	body   string
}

type Store interface {
	Store(key string, value *Response)
	Load(key string) (*Response, error)
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
			for k1, v1 := range result.header {
				c.Writer.Header().Set(k1, v1)
			}
			_, err := c.Writer.Write([]byte(result.body))
			if err != nil {
				logrus.Error(err)
				return
			}
			c.Abort()
			return
		}
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		// 否则，继续处理请求
		c.Next()
		// 将请求的响应结果存储起来，作为幂等性响应
		if c.Writer.Status() != http.StatusOK {
			return
		}
		responseBody := blw.body.Bytes()
		header := map[string]string{}
		for k1, v1 := range c.Writer.Header() {
			header[k1] = v1[0]
		}
		response := &Response{
			header: header,
			body:   string(responseBody),
		}
		store.Store(k, response)
	}
}
