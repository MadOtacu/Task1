package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var src = flag.String("src", "", "Ссылка на файл с сылками")

var dst = flag.String("dst", "", "Путь для записи с файлов")

func main() {
	start := time.Now()
	flag.PrintDefaults()
	flag.Parse()

	var byteText, err = os.ReadFile(*src)
	if err != nil {
		panic(err)
	}

	if _, errCheck := os.Stat(*dst); os.IsNotExist(errCheck) {
		errCreation := os.Mkdir(*dst, 0777)
		if errCreation != nil {
			panic(errCreation)
		}
	}

	text := string(byteText)
	textArr := strings.Split(text, "\n")

	for i, element := range textArr {
		if !strings.HasPrefix(element, "https://") {
			continue
		}
		var temp bytes.Buffer

		resp, errServ := http.Get(element)
		if errServ != nil {
			fmt.Println("Хост " + element + " не отвечает")
			continue
		}

		copyErr := resp.Write(&temp)
		if copyErr != nil {
			panic(copyErr)
		}

		num := strconv.Itoa(i)
		fmt.Println("Запись данных из " + element + " в файл " + num + ".html")
		writeErr := ioutil.WriteFile(*dst+num+".html", temp.Bytes(), 0777)
		if writeErr != nil {
			panic(writeErr)
		}

		fmt.Println("В файл " + num + ".html сохранен шаблон сайта " + element)
		resp.Body.Close()
	}
	elapsed := time.Since(start)
	fmt.Println("Время выполнения программы:", elapsed)
}
