package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
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

func ReadArgs(args []string) string {
	if args == nil || len(args) == 1 {
		return ""
	}
	return args[len(args)-1]
}

func HandleMemo(line string) string {
	memo := strings.ReplaceAll(line, "\n", " ")

	memo = strings.ReplaceAll(memo, "SEPA-BASISLASTSCHR. ", "")
	memo = strings.ReplaceAll(memo, "Lastschrift ", "")
	memo = strings.ReplaceAll(memo, "IBAN: DE27510900000044145103", "book-n-drive mobilitaetssysteme GmbH")

	return memo
}

func Convert(content [][]string, memo_column int) [][]string {
	var result [][]string
	result = append(result, []string{"Date", "Payee", "Memo", "Amount"})
	for _, line := range content {

		size := len(line)
		var amount string
		if line[size-1] == "S" {
			amount = "-"
		}
		amount += line[size-2]

		test := []string{
			line[0],
			"MVB",
			HandleMemo(line[size-memo_column]),
			amount,
		}
		result = append(result, test)
	}
	return result
}

func ReadFile(filename string) (error, [][]string, int) {
	file, err := os.Open(filename)
	if err != nil {
		Sugar.Error(err)
		return err, nil, 0
	}

	bytes, err := ioutil.ReadAll(file)
	str := strings.Split(string(bytes), "\n")
	var memo_column int

	if strings.Index(str[12], "Buchungstag") == 0 {
		str = str[13 : len(str)-4]
		memo_column = 5
	} else {
		str = str[16 : len(str)-4]
		memo_column = 5
	}

	joined := strings.Join(str, "\n")

	stringreader := strings.NewReader(joined)
	reader := csv.NewReader(stringreader)
	reader.Comma = ';'
	var content [][]string

	content, err = reader.ReadAll()
	if err != nil {
		Sugar.Error(err)
		return err, nil, 0
	}

	return nil, content, memo_column
}

func WriteFile(filename string, content [][]string) error {
	Sugar.Info("Writing file to: ", filename)
	file, err := os.Create(filename)
	if err != nil {
		Sugar.Error(err)
		return err
	}

	w := csv.NewWriter(file)
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
	err, content, memo_column := ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = WriteFile(filename+"_converted.csv", Convert(content, memo_column))
	if err != nil {
		panic(err)
	}
	err = WriteFile(filename+"_striped.csv", content)
	if err != nil {
		panic(err)
	}
	Sugar.Info("Happy Ending")
}
