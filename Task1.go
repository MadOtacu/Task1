package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Функция вызова get-запроса и записи результата в файл
func connect(URLtoGet string, iter int, dst *string, wg *sync.WaitGroup) {
	// Завершаем поток при выходе из функции
	defer wg.Done()

	// Если строка содержит корректный префикс пропускает ее дальше
	// Если нет, то выводит предупреждение и завершает функцию
	if !strings.HasPrefix(URLtoGet, "https://") {
		fmt.Println("Элемент " + URLtoGet + " не является ссылкой")
		return
	}

	// Массив байтов для копирования в файл
	bytesToCopy := bytes.Buffer{}

	// Get запрос к сайту и получение ответа
	resp, errServ := http.Get(URLtoGet)
	// Обработка ошибки ответа сервера
	// При возникновении выводит предупреждение и завершает функцию
	if errServ != nil {
		fmt.Println("Хост " + URLtoGet + " не отвечает")
		return
	}
	// Закрытие подключенного Get запроса
	defer resp.Body.Close()

	// Запись ответа сайта в массив байтов
	errWrite := resp.Write(&bytesToCopy)
	// Проверка записи ответа в массив байтов, при ошибке завершает функцию
	if errWrite != nil {
		log.Println(errWrite)
		return
	}

	// Конвертация номера итерации в строку
	numToCreate := strconv.Itoa(iter)

	// Попытка записи данных из массива байтов в генерируемый файл, и вывод уведомления об этом
	fmt.Println("Запись данных из " + URLtoGet + " в файл " + numToCreate + ".html")
	errFileGenerating := os.WriteFile(*dst+numToCreate+".html", bytesToCopy.Bytes(), 0777)
	// Обработка ошибки записи в файл, при ошибке завершает функцию
	if errFileGenerating != nil {
		log.Println(errFileGenerating)
		return
	}

	fmt.Println("В файл " + numToCreate + ".html сохранен шаблон сайта " + URLtoGet)
}

func main() {
	// Создание WaitGroup
	var wg sync.WaitGroup

	// Создание флагов
	var src = flag.String("src", "", "Ссылка на файл с сcылками")
	var dst = flag.String("dst", "", "Путь для записи файлов")

	// Запуск таймера выполнения программы
	start := time.Now()

	// Парсинг флагов
	flag.Parse()

	// Чтение данных из переданного файла
	var byteFileContent, errReadingFile = os.ReadFile(*src)
	// При некорректной передаче флага возвращает данные о флагах и завершает программу
	if errReadingFile != nil {
		flag.PrintDefaults()
		os.Exit(2)
	}

	// Проверка на наличие директории
	if _, errDirCheck := os.Stat(*dst); errDirCheck != nil {
		// Создание директории при необходимости
		errDirCreation := os.MkdirAll(*dst, 0777)
		// Проверка создания директории, при некорректной передаче флага возвращает данные о флагах и завершает программу
		if errDirCreation != nil {
			flag.PrintDefaults()
			os.Exit(2)
		}
	}

	// Преобразование массива битов в строку
	textFileContent := string(byteFileContent)
	// Разбиение сплошной строки на массив строк
	textArrFileContent := strings.Split(textFileContent, "\n")
	// Удаление последнего ненужного символа
	textTotalFileContent := textArrFileContent[:len(textArrFileContent)-1]
	// Перебор строк, индексация WaitGroup и применение горутин метода connect
	for i, stringToGet := range textTotalFileContent {
		wg.Add(1)
		go connect(stringToGet, i, dst, &wg)
	}
	// Ожидание завершения всех горутин
	wg.Wait()
	// Запись завершения времени программы на счетчик и вывод значения
	elapsed := time.Since(start)
	fmt.Println("Время выполнения программы:", elapsed)
}
