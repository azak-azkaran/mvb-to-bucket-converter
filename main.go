package main

import (
	"encoding/csv"
	"fmt"
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

func Convert(content [][]string) [][]string{
	var result [][]string
	for _,line := range content{
		result = append(result, line)
	}
	return content
}

func ReadFile(filename string) ( error, [][]string ){
	file, err := os.Open(filename)
    if err != nil {
        Sugar.Error(err)
				return err, nil
			}

	  reader := csv.NewReader(file)
		reader.Comma = ';'
		var content [][]string

		content, err = reader.ReadAll()
    if err != nil {
        Sugar.Error(err)
				return err, nil
			}

		content = content[12:]
		return nil, content
}

func main() {
	Sugar.Info("Hello World")
	filename := ReadArgs(os.Environ())
	ReadFile(filename)
}
