package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadArgs(t *testing.T){
	fmt.Println("Testing ReadArgs...")
	args := []string{"filename", "./test/testfile.csv"}

	filename := ReadArgs(args)
	assert.NotEmpty(t, filename)
	assert.Equal(t, "./test/testfile.csv", filename)
}

func TestReadFile(t *testing.T){
	fmt.Println("Testing ReadFiles...")

	err,content := ReadFile("./test/testfile.csv")
	assert.NoError(t,err)
	assert.NotNil(t,content)
	assert.True(t,len(content) >= 2)
	fmt.Println(content)
}

func TestConvert(t *testing.T){
	fmt.Println("Testing Convert...")

	content := [][]string{
		{"01.02.21", "01.02.21", "AAAAAAA","AMAZON PAYMENTS EUROPE S.C.A.","SEPA-BASISLASTSCHR.\n111-2222222-3333333 AMZN Mk","","EUR","200","S"},
		{"02.02.21", "02.02.21", "AAAAAAA","AMAZON PAYMENTS EUROPE S.C.A.","SEPA-BASISLASTSCHR.\n111-2222222-3333333 AMZN Mk","","EUR","200","H"},
	}
	result := Convert(content)
	fmt.Println(result)
	assert.NotEmpty(t,result)
	assert.True(t,len(result) ==3)
	assert.True(t,len(result[0]) ==4)
	assert.Equal(t, "01.02.21",result[1][0])
	assert.Equal(t , "MVB",result[1][1])
	assert.Equal(t , "111-2222222-3333333 AMZN Mk",result[1][2])
	assert.Equal(t , "-200",result[1][3])
	assert.Equal(t , "200",result[2][3])
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
		{"01.02.21", "01.02.21", "AAAAAAA","AMAZON PAYMENTS EUROPE S.C.A.","SEPA-BASISLASTSCHR.\n111-2222222-3333333 AMZN Mk","","EUR","200","S"},
	}
	WriteFile(filename, content)
	_, err = os.Stat(filename)
	assert.NoError(t,err)
}
