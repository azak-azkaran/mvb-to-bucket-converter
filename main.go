package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	LogConfig                    = zap.NewProductionConfig()
	Sugar     *zap.SugaredLogger = LogInit()
)

func LogInit() *zap.SugaredLogger {
	var err error
	LogConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	LogConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	LogConfig.Encoding = "console"
	logger, err := LogConfig.Build()
	if err != nil {
		fmt.Println("Error building logger:", err)
	}
	defer logger.Sync() // flushes buffer, if any
	return logger.Sugar()
}

func ReadArgs(args []string) string{
	if args == nil || len(args)== 1{
		return ""
	}
	return args[1]
}

func ReadFile(filename string) ( error, *string ){
	  content, err := ioutil.ReadFile(filename)
    if err != nil {
        Sugar.Error(err)
				return err, nil
			}
		stuff := string(content)
		return nil, &stuff
}

func main() {
	Sugar.Info("Hello World")
	filename := ReadArgs(os.Environ())
	ReadFile(filename)
}
