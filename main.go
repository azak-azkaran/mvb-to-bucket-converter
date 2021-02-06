package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

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
	return args[len(args)-1]
}

func HandleMemo(line string) string{
	var memo string
	split := strings.Split(line, "\n")
	if len(split) == 0 {
		return line
	}

	for _, l := range(split[1:]){
		memo = memo + l
	}
	return memo
}

func Convert(content [][]string) [][]string{
	var result [][]string
	content = content[1:]
	result = append(result, []string{"Date","Payee","Memo","Amount"})
	for _,line := range content{
		var amount string
		if line[8] == "S"{
			amount = "-"
		}
		amount += line[7]

		test := []string{
			line[0],
			"MVB",
			HandleMemo(line[4]),
			amount,
		}
		result = append(result, test)
	}
	return result
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
		content = content[:len(content)-3]
		return nil, content
}

func WriteFile(filename string, content [][]string) error{
	Sugar.Info("Writing file to: ", filename)
	file, err := os.Create(filename)
	if err != nil {
		Sugar.Error(err)
		return err
	}

	w:= csv.NewWriter(file)
	w.WriteAll(content)
	if w.Error() != nil {
		Sugar.Error(err)
		return err
	}
	return nil
}

func main() {
	Sugar.Info("Hello World")
	filename := ReadArgs(os.Args)
	err, content := ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = WriteFile(filename + "_converted.csv", Convert( content ))
	if err != nil {
		panic(err)
	}
	err = WriteFile(filename + "_striped.csv", content )
	if err != nil {
		panic(err)
	}
	Sugar.Info("Happy Ending")
}