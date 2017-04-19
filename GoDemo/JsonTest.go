package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"bytes"
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


	type App struct {
		Id string `json:"id"`
		Title string `json:"title"`
	}

	data := []byte(`
    {
        "id": "k34rAT4",
        "title": "My Awesome App"
    }
	`)

	var app App
	var errApp error

	errApp = json.Unmarshal(data, &app)
	if (errApp != nil) {
		fmt.Println("Error Json ")
	}

	fmt.Println("app.title: " + app.Title)
	fmt.Println(bytes.NewBuffer(data).String())
	fmt.Println(string(data))
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
