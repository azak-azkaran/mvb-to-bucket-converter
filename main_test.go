package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadArgs(t *testing.T) {
	fmt.Println("Testing ReadArgs...")
	args := []string{"filename", "./test/testfile.csv"}

	filename := ReadArgs(args)
	assert.NotEmpty(t, filename)
	assert.Equal(t, "./test/testfile.csv", filename)
}

func TestReadFile(t *testing.T) {
	fmt.Println("Testing ReadFiles...")

	err, content, _, memo, value := ReadFile("./test/testfile.csv")
	assert.NoError(t, err)
	assert.NotNil(t, content)
	assert.True(t, len(content) >= 2)
	assert.Equal(t, 5, memo)
	assert.Equal(t, 2, value)
	fmt.Println(content)

	err, content, _, memo, value = ReadFile("./test/testfile2.csv")
	assert.NoError(t, err)
	assert.NotNil(t, content)
	assert.True(t, len(content) >= 1)
	assert.Equal(t, 5, memo)
	assert.Equal(t, 2, value)
	fmt.Println(content)
}

func TestConvert(t *testing.T) {
	fmt.Println("Testing Convert...")

	content := [][]string{
		{"01.02.21", "01.02.21", "AAAAAAA", "AMAZON PAYMENTS EUROPE S.C.A.", "SEPA-BASISLASTSCHR.\n111-2222222-3333333 AMZN Mk", "", "EUR", "200", "S"},
		{"02.02.21", "02.02.21", "AAAAAAA", "AMAZON PAYMENTS EUROPE S.C.A.", "SEPA-BASISLASTSCHR.\n111-2222222-3333333 AMZN Mk", "", "EUR", "200", "H"},
	}
	result := Convert(content, 0, 5, 2)
	fmt.Println(result)
	assert.NotEmpty(t, result)
	assert.True(t, len(result) == 3)
	assert.True(t, len(result[0]) == 4)
	assert.Equal(t, "01/02/21", result[1][0])
	assert.Equal(t, "MVB", result[1][1])
	assert.Equal(t, "111-2222222-3333333 AMZN Mk", result[1][2])
	assert.Equal(t, "-200", result[1][3])
	assert.Equal(t, "200", result[2][3])
}

func TestWriteFile(t *testing.T) {
	fmt.Println("Testing WriteFile...")
	filename := "./test/newfile.csv"
	t.Cleanup(func() {
		os.Remove(filename)
	})

	_, err := os.Stat(filename)
	assert.True(t, os.IsNotExist(err))

	content := [][]string{
		{"01.02.21", "01.02.21", "AAAAAAA", "AMAZON PAYMENTS EUROPE S.C.A.", "SEPA-BASISLASTSCHR.\n111-2222222-3333333 AMZN Mk", "", "EUR", "200", "S"},
	}
	WriteFile(filename, content)
	assert.FileExists(t, filename)
}

func TestMainTest1(t *testing.T) {
	fmt.Println("Testing Main...")
	t.Cleanup(func() {
		os.Remove("./test/testfile.csv_converted.csv")
		os.Remove("./test/testfile.csv_striped.csv")
		os.Remove("./test/testfile2.csv_converted.csv")
		os.Remove("./test/testfile2.csv_striped.csv")
	})

	os.Args = []string{
		"mvb-to-bucket-converter",
		"./test/testfile.csv",
	}

	main()
	assert.FileExists(t, "./test/testfile.csv_converted.csv")
	assert.FileExists(t, "./test/testfile.csv_striped.csv")

	file, err := os.Open("./test/testfile.csv_converted.csv")
	require.NoError(t, err)
	reader := csv.NewReader(file)
	require.NotNil(t, reader)
	require.NoError(t, err)
	content, err := reader.ReadAll()
	require.NoError(t, err)
	assert.Len(t, content, 3)
	assert.Len(t, content[0], 4)
	assert.Len(t, content[1], 4)
	assert.Len(t, content[2], 4)

	assert.Equal(t, "Memo", content[0][2])
	assert.Equal(t, "Amount", content[0][3])
	assert.Equal(t, "-200", content[1][3])
	assert.Equal(t, "-39,05", content[2][3])

	file, err = os.Open("./test/testfile.csv_striped.csv")
	require.NoError(t, err)
	reader = csv.NewReader(file)
	require.NotNil(t, reader)
	require.NoError(t, err)
	content, err = reader.ReadAll()
	require.NoError(t, err)
	assert.Len(t, content, 2)
	assert.Len(t, content[0], 9)
	assert.Len(t, content[1], 9)
	assert.Equal(t, "200", content[0][7])
	assert.Equal(t, "S", content[0][8])
	assert.Equal(t, "39,05", content[1][7])
	assert.Equal(t, "S", content[1][8])
}

func TestMainTest2(t *testing.T) {
	fmt.Println("Testing2 Main...")
	t.Cleanup(func() {
		os.Remove("./test/testfile2.csv_converted.csv")
		os.Remove("./test/testfile2.csv_striped.csv")
	})

	os.Args = []string{
		"mvb-to-bucket-converter",
		"./test/testfile2.csv",
	}

	main()
	assert.FileExists(t, "./test/testfile2.csv_converted.csv")
	assert.FileExists(t, "./test/testfile2.csv_striped.csv")

	file, err := os.Open("./test/testfile2.csv_converted.csv")
	require.NoError(t, err)
	reader := csv.NewReader(file)
	require.NotNil(t, reader)
	require.NoError(t, err)
	content, err := reader.ReadAll()
	require.NoError(t, err)
	assert.Len(t, content, 2)
	assert.Len(t, content[0], 4)
	assert.Len(t, content[1], 4)

	assert.Equal(t, "Memo", content[0][2])
	assert.Equal(t, "Amount", content[0][3])

	assert.Equal(t, "-4,41", content[1][3])
	assert.True(t, strings.Contains(content[1][2], "Patreon"), "Message: "+content[1][2])
}

func TestMainTest3(t *testing.T) {
	fmt.Println("Testing3 Main...")
	t.Cleanup(func() {
		os.Remove("./test/testfile3.csv_converted.csv")
		os.Remove("./test/testfile3.csv_striped.csv")
	})

	os.Args = []string{
		"mvb-to-bucket-converter",
		"./test/testfile3.csv",
	}

	main()
	assert.FileExists(t, "./test/testfile3.csv_converted.csv")
	assert.FileExists(t, "./test/testfile3.csv_striped.csv")

	file, err := os.Open("./test/testfile3.csv_converted.csv")
	require.NoError(t, err)
	reader := csv.NewReader(file)
	require.NotNil(t, reader)
	require.NoError(t, err)
	content, err := reader.ReadAll()
	require.NoError(t, err)
	assert.Len(t, content, 2)
	assert.Len(t, content[0], 4)
	assert.Len(t, content[1], 4)

	assert.Equal(t, "Memo", content[0][2])
	assert.Equal(t, "Amount", content[0][3])

	assert.Equal(t, "-37,99", content[1][3])
	assert.True(t, strings.Contains(content[1][2], "AMZN"), "Message: "+content[1][2])
}
