package logzap

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ConfigZapLogger() {

	cfgInfo := zap.NewProductionEncoderConfig()
	cfgInfo.EncodeLevel = zapcore.CapitalLevelEncoder
	cfgInfo.TimeKey = ""
	cfgInfo.CallerKey = ""
	cfgInfo.FunctionKey = ""

	cfgErr := zap.NewProductionEncoderConfig()
	cfgErr.EncodeLevel = zapcore.CapitalLevelEncoder
	cfgErr.TimeKey = "at"
	cfgErr.CallerKey = "call_in"
	cfgErr.FunctionKey = "func"
	cfgErr.StacktraceKey = "stack"
	cfgErr.EncodeCaller = zapcore.ShortCallerEncoder
	cfgErr.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")

	infoLvl := zap.LevelEnablerFunc(func(level zapcore.Level) bool { return level == zapcore.InfoLevel })
	warnLvl := zap.LevelEnablerFunc(func(level zapcore.Level) bool { return level >= zapcore.WarnLevel })
	fatalLvl := zap.LevelEnablerFunc(func(level zapcore.Level) bool { return level > zapcore.ErrorLevel })

	infoEncoder := zapcore.NewConsoleEncoder(cfgInfo)
	errsEncoder := zapcore.NewJSONEncoder(cfgErr)

	core := zapcore.NewTee(
		zapcore.NewCore(infoEncoder, zapcore.AddSync(os.Stdout), infoLvl),
		zapcore.NewCore(errsEncoder, zapcore.AddSync(os.Stderr), warnLvl),
	)

	logger := zap.New(core)
	logger = logger.WithOptions(zap.AddCaller(), zap.AddStacktrace(fatalLvl))

	zap.ReplaceGlobals(logger)
}
