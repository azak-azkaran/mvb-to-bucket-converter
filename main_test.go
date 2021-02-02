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
	assert.True(t,len(*content) == 2)
}
