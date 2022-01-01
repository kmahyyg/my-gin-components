package logger

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// LoggerBuild returns zap-based Gin request logging middleware
// @param none
// @return logger: []gin.HandlerFunc
func LoggerBuild() []gin.HandlerFunc {
	logger, _ := zap.NewProduction()
	panicGinLog := ginzap.RecoveryWithZap(logger, true)
	normalGinLog := ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
	})
	finLoggers := []gin.HandlerFunc{panicGinLog,normalGinLog}
	return finLoggers
}