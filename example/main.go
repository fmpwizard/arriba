package main

//This is a sample application that shows how to use arriba
//cd into the example folder and run go run main.go

import (
	"fmt"
	"github.com/fmpwizard/arriba"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var funcMap = make(map[string]HTMLTransform)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.HandleFunc("/index", home)
	//funcMap holds a map of function names as they appear on the html and maps to the real function to call
	funcMap["ChangeTime"] = ChangeTime
	http.ListenAndServe(":7070", nil)
}

func home(rw http.ResponseWriter, req *http.Request) {
	t, err := ioutil.ReadFile("index.html")
	if err != nil {
		panic(err)
	}

	for _, partialHTML := range arriba.GetFunctions(string(t)) {
		//This is a silly way to replace the old html with new one, because
		//it will fail if we have the same raw html more than once.
		//t = []byte(strings.Replace(string(t), html, funcMap[functionName](html), 1))
		fmt.Printf("Function name : %+v\n\n", partialHTML.FunctionName)
		fmt.Printf("Function html : %+v\n\n", partialHTML.FunctionHtml)
		fmt.Printf("Raw HTML : %+v\n\n", partialHTML.RawHTML)
	}

	rw.Header().Add("Content-Type", "text/html; charset=UTF-8")
	rw.Write(t)

}

//ChangeTime takes a portion of the complete html page and does a replacement
//Future versions will use css transformations
func ChangeTime(html string) string {
	return strings.Replace(html, "Time goes here", time.Now().Format("2006-01-02T15:04:05.999999999Z07:00"), 1)
}

/*type HTMLTransform interface {
  ServeHTTP(ResponseWriter, *Request)
}*/

//HTMLTransform is the type of the functions we allow to do html transformation.
//This is too generic, but works for now.
type HTMLTransform func(string) string
