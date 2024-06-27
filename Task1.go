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
)

var src = flag.String("src", "", "Input file path")
var dst = flag.String("dst", "", "Output file path")

func errHandler(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	flag.Parse()
	var byteText, err = os.ReadFile(*src)
	errHandler(err)
	text := string(byteText[:])
	fmt.Println(text)
	textArr := strings.Split(text, "\n")

	for i, element := range textArr {
		if strings.HasPrefix(element, "https://") {
			var temp bytes.Buffer

			resp, errServ := http.Get(element)
			errHandler(errServ)

			copyErr := resp.Write(&temp)
			errHandler(copyErr)

			num := strconv.Itoa(i)
			fmt.Println(num)
			writeErr := ioutil.WriteFile(*dst+num+".html", temp.Bytes(), 0666)
			errHandler(writeErr)

			fmt.Println(element)
		}
	}
}
