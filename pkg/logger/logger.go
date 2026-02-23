package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(envMode string) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder // Set the log level to InfoLevel for production and DebugLevel for development

	fileEncoder := zapcore.NewJSONEncoder(config)

	writer := zapcore.AddSync(os.Stdout)

	defaultLoglevel := zapcore.InfoLevel

	if envMode == "development" {
		defaultLoglevel = zapcore.DebugLevel
	}

	core := zapcore.NewCore(fileEncoder, writer, defaultLoglevel)

	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// if envMode == "production" {
	// 	config = zap.NewProductionConfig()
	// } else {
	// 	config = zap.NewDevelopmentConfig()
	// 	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// }

	// config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// var err error
	// Log, err = config.Build()
	// if err != nil {
	// 	panic(err)
	// }
}
