package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Student struct {
	Name    string
	Age     int
	Guake   bool
	Classes []string
	Price   float32
}

func main() {
	st := &Student {
		"Xiao Ming",
		16,
		true,
		[]string{"Math", "English", "Chinese"},
		9.99,
	}

	b, err := json.Marshal(st)

	println(b)
	println(err)

	var isvalid bool

	isvalid = false

	isValidString := strconv.FormatBool(isvalid)


	fmt.Println("------" + isValidString)

	fmt.Println(GetMockData())

	timestamp := time.Now().Unix()

	fmt.Println(timestamp)

	fmt.Println(time.Now().String())

	fmt.Println(time.Unix(timestamp, 0).String())

	fmt.Println("Int convert")
	var sumVal int64
	var sumByte []byte
	sumVal = 1000000000

	sumByte = []byte(strconv.Itoa(sumVal))
}

func GetMockData() ([]string) {
	s := []string{
		"1",
		"2",
		"hello",
		"hello",
	}
	return s
}
