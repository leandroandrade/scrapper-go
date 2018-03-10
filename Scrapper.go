package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const site = "https://httpstatusdogs.com/"

const output = "/tmp/images-scrapper/"

func main() {
	fmt.Println("Inciando scrapper!")
	request()
	fmt.Println("Scrapper finalizado!")
}

func request() {
	response, err := http.Get(site)
	check(err)

	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	check(err)

	regex := regexp.MustCompile("<* src=[^>]*\\salt=[^>]*")
	values := regex.FindAll(content, -1)

	for _, value := range values {
		data := formatHTMLTag(value)
		image, description := getValues(data)

		downloadImage(image, description)
	}

}

func downloadImage(image string, description string) {
	url := site + image

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

func getValues(data []string) (image string, description string) {
	image = strings.Join(data[:1], "")
	description = strings.Join(data[1:], "")
	return
}

func formatHTMLTag(value []byte) []string {
	line := strings.TrimSpace(string(value))
	replacer := strings.NewReplacer("\"", "", "src=", "", "alt=", "")
	return strings.Split(replacer.Replace(line), " ")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
