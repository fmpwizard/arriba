package main

//This is a sample application that shows how to use arriba
//cd into the example folder and run go run main.go

import (
	"github.com/fmpwizard/arriba"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.HandleFunc("/index", home)
	//FunctionMap holds a map of function names as they appear on the html and maps to the real function to call
	arriba.FunctionMap.Lock()
	arriba.FunctionMap.M["ChangeName"] = ChangeName
	arriba.FunctionMap.M["ChangeLastName"] = ChangeLastName
	arriba.FunctionMap.Unlock()
	http.ListenAndServe(":7070", nil)
}

func home(rw http.ResponseWriter, req *http.Request) {
	t, err := ioutil.ReadFile("index.html")
	if err != nil {
		panic(err)
	}

	value := arriba.MarshallElem(string(t))
	rw.Header().Add("Content-Type", "text/html; charset=UTF-8")
	rw.Write([]byte(value))

}

func ChangeName(html string) string {
	return strings.Replace(html, "Diego", "Gabriel", 1)
}

func ChangeLastName(html string) string {
	return strings.Replace(html, "Medina", "Bauman", 1)
}
