package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"bytes"
)

const site = "https://httpstatusdogs.com/"
const output = "OUTPUT_DIR"

var regex = regexp.MustCompile(`([0-9]{3}.jpg)`)

func main() {
	fmt.Println("Inciando scrapper!")
	process()
	fmt.Println("Scrapper finalizado!")
}

func process() {
	response, err := http.Get(site)
	check(err)
	defer response.Body.Close()

	var write bytes.Buffer
	io.Copy(&write, response.Body)

	values := regex.FindAllString(write.String(), -1)

	for _, image := range values {
		downloadImage(image)
	}
}

func downloadImage(image string) {
	url := site + "img/" + image

	response, err := http.Get(url)
	check(err)
	defer response.Body.Close()

	replacer := strings.NewReplacer("/", "-")
	file, err := os.Create(output + replacer.Replace(image))
	check(err)

	_, err = io.Copy(file, response.Body)
	check(err)

	file.Close()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
