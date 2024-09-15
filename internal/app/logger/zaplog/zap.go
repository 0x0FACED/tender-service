package zaplog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	log *zap.Logger
}

func New() *ZapLogger {
	dirName := "logs"
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		log.Fatalln("cant make dir: ", err)
		return nil
	}

	filename := time.Now().Format("2006-01-02") + ".log"
	filePath := filepath.Join(dirName, filename)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("cant open file: ", err)
		return nil
	}

	config := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	cEnc := zapcore.NewConsoleEncoder(config)
	fEnc := zapcore.NewConsoleEncoder(config)

	core := zapcore.NewTee(
		zapcore.NewCore(cEnc, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(fEnc, zapcore.AddSync(file), zapcore.DebugLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	return &ZapLogger{
		log: logger,
	}
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("[2006-01-02 | 15:04:05]"))
}

func convertAnyToZapFields(fields ...any) ([]zapcore.Field, error) {
	flds := make([]zapcore.Field, 0, len(fields))
	for _, f := range fields {
		switch v := f.(type) {
		case zapcore.Field:
			flds = append(flds, v)
		default:
			return nil, fmt.Errorf("unexpected type: %T, expected zapcore.Field", v)
		}
	}
	return flds, nil
}

func (z *ZapLogger) Info(wrappedMsg string, fields ...any) {
	flds, err := convertAnyToZapFields(fields...)
	if err != nil {
		z.log.Panic("cannot convert", zap.Error(err))
	}
	z.log.Info("[MSG]: "+wrappedMsg, flds...)
}

func (z *ZapLogger) Debug(wrappedMsg string, fields ...any) {
	flds, err := convertAnyToZapFields(fields...)
	if err != nil {
		z.log.Panic("cannot convert", zap.Error(err))
	}
	z.log.Debug("[MSG]: "+wrappedMsg, flds...)
}

func (z *ZapLogger) Error(wrappedMsg string, fields ...any) {
	flds, err := convertAnyToZapFields(fields...)
	if err != nil {
		z.log.Panic("cannot convert", zap.Error(err))
	}
	z.log.Error("[MSG]: "+wrappedMsg, flds...)
}

func (z *ZapLogger) Fatal(wrappedMsg string, fields ...any) {
	flds, err := convertAnyToZapFields(fields...)
	if err != nil {
		z.log.Panic("cannot convert", zap.Error(err))
	}
	z.log.Fatal("[MSG]: "+wrappedMsg, flds...)
}
