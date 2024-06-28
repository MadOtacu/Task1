package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var src = flag.String("src", "", "Ссылка на файл с сcылками")

var dst = flag.String("dst", "", "Путь для записи файлов")

func connect(element string, i int, wg sync.WaitGroup) {
	defer wg.Done()

	if !strings.HasPrefix(element, "https://") {
		fmt.Println("Элемент " + element + " не является ссылкой")
		return
	}

	data := bytes.Buffer{}

	resp, errServ := http.Get(element)
	if errServ != nil {
		fmt.Println("Хост " + element + " не отвечает")
		return
	}
	defer resp.Body.Close()

	copyErr := resp.Write(&data)
	if copyErr != nil {
		panic(copyErr)
	}

	num := strconv.Itoa(i)
	fmt.Println("Запись данных из " + element + " в файл " + num + ".html")
	writeErr := os.WriteFile(*dst+num+".html", data.Bytes(), 0777)
	if writeErr != nil {
		panic(writeErr)
	}

	fmt.Println("В файл " + num + ".html сохранен шаблон сайта " + element)
}

func main() {
	var wg sync.WaitGroup
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
	textTotal := textArr[:len(textArr)-1]
	fmt.Println(len(textTotal))
	for i, element := range textTotal {
		wg.Add(1)
		fmt.Println("ELEMENT", element)
		go connect(element, i, wg)
	}
	wg.Wait()
	fmt.Println("WG dONE")
	elapsed := time.Since(start)
	fmt.Println("Время выполнения программы:", elapsed)
}
