package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func Logger(targetPackage string) *zap.SugaredLogger {
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	writerSyncer := os.Stdout
	levelEnabler := zapcore.DebugLevel

	core := zapcore.NewCore(encoder, writerSyncer, levelEnabler)

	logger := zap.New(core, zap.Fields(zap.String("package", targetPackage)))
	return logger.Sugar()
}
