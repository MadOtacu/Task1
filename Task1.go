package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var src = flag.String("src", "", "Ссылка на файл с сcылками")

var dst = flag.String("dst", "", "Путь для записи файлов")

func connect(element string) (*bytes.Buffer, error) {
	if !strings.HasPrefix(element, "https://") {
		return nil, errors.New("Invalid link")
	}

	data := bytes.Buffer{}

	resp, errServ := http.Get(element)
	if errServ != nil {
		fmt.Println("Хост " + element + " не отвечает")
		return nil, errors.New("Host is not responding")
	}
	defer resp.Body.Close()

	copyErr := resp.Write(&data)
	if copyErr != nil {
		panic(copyErr)
	}
	return &data, nil
}

func main() {
	start := time.Now()
	flag.Parse()

	var byteText, err = os.ReadFile(*src)
	if err != nil {
		flag.PrintDefaults()
		os.Exit(2)
	}

	if _, errCheck := os.Stat(*dst); os.IsNotExist(errCheck) {
		errCreation := os.Mkdir(*dst, 0777)
		if errCreation != nil {
			flag.PrintDefaults()
			os.Exit(2)
		}
	}

	text := string(byteText)
	textArr := strings.Split(text, "\n")

	for i, element := range textArr {
		temp, connectErr := connect(element)
		if connectErr != nil {
			continue
		}

		num := strconv.Itoa(i)
		fmt.Println("Запись данных из " + element + " в файл " + num + ".html")
		writeErr := os.WriteFile(*dst+num+".html", temp.Bytes(), 0777)
		if writeErr != nil {
			panic(writeErr)
		}

		fmt.Println("В файл " + num + ".html сохранен шаблон сайта " + element)
	}
	elapsed := time.Since(start)
	fmt.Println("Время выполнения программы:", elapsed)
}
