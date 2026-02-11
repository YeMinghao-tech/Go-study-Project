package router

import (
	"bytes"
	"io"
	"mall/consts"
	"mall/utils/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetRequestBody(ctx *gin.Context) string {
	data, _ := io.ReadAll(ctx.Request.Body)
	return string(data)
}

type responseWriterWrapper struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func AccessLogMiddleware(filter func(*gin.Context) bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if filter != nil && !filter(ctx) {
			ctx.Next()
			return
		}
		body := GetRequestBody(ctx)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(body)))
		begin := time.Now()
		fields := []zap.Field{
			zap.String("ip", ctx.RemoteIP()),                        // 客户端IP
			zap.String("method", ctx.Request.Method),                // HTTP方法（GET/POST/PUT等）
			zap.String("path", ctx.Request.URL.Path),                // 请求路径（如/admin/user）
			zap.String("params", ctx.Request.URL.RawQuery),          // URL查询参数（如?page=1&size=10）
			zap.Any("body", body),                                   // 请求体内容
			zap.String("token", ctx.GetHeader(consts.UserTokenKey)), // 从请求头获取用户Token
		}
		var responseBody bytes.Buffer
		multiWriter := io.MultiWriter(ctx.Writer, &responseBody) // 多写器：同时写入原响应和缓冲区
		ctx.Writer = &responseWriterWrapper{
			ResponseWriter: ctx.Writer,  // 保留原生ResponseWriter
			Writer:         multiWriter, // 写入到多写器
		}

		ctx.Next()
		respBody := responseBody.String()
		if len(respBody) > 1024 {
			respBody = respBody[:1024]
		}
		fields = append(fields, zap.Int64("dur_ms", time.Since(begin).Milliseconds()))
		fields = append(fields, zap.Int("status", ctx.Writer.Status()))
		fields = append(fields, zap.String("resp", respBody))
		logger.Info("access_log", fields...)
	}
}
