package middlewares

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Configure a custom logger
func newLogger() *zap.Logger {
    config := zap.NewProductionEncoderConfig()
    config.TimeKey = "timestamp"  // Custom timestamp key
    config.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339) // Readable timestamp format
    config.CallerKey = ""  // Disable caller field

    encoder := zapcore.NewJSONEncoder(config)
    core := zapcore.NewCore(encoder, zapcore.AddSync(zapcore.Lock(os.Stdout)), zap.InfoLevel)
    
    return zap.New(core)
}

var logger = newLogger()

func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        serviceName := c.GetHeader("x-source-service")
        clientIP := c.ClientIP()

        logger.Info("Incoming request",
            zap.String("service", serviceName),
            zap.String("client_ip", clientIP),
            zap.String("path", c.Request.URL.Path),
        )

        c.Next()
    }
}
