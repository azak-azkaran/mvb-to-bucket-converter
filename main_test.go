package main

import (
	"fmt"
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
		{"01.02.21", "01.02.21", "AAAAAAA","AMAZON PAYMENTS EUROPE S.C.A.","SEPA-BASISLASTSCHR.\n111-2222222-3333333 AMZN Mk","","EUR 200","S"},
	}
	result := Convert(content)
	assert.NotEmpty(t,result)
	assert.True(t,len(result) ==1)
	assert.True(t,len(result[0]) ==4)
	assert.True(t,result[0][0] == "01.02.21")
	assert.True(t,result[0][1] == "MVB")
	assert.True(t,result[0][1] == "111-2222222-3333333 AMZN Mk")
	assert.True(t,result[0][3] == "-200")
}
