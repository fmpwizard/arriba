package main

//This is a sample application that shows how to use arriba
//cd into the example folder and run go run main.go

import (
	"github.com/fmpwizard/arriba"
	"github.com/fmpwizard/arriba/vendor/code.google.com/p/go-html-transform/h5"
	"github.com/fmpwizard/arriba/vendor/code.google.com/p/go-html-transform/html/transform"
	"github.com/fmpwizard/arriba/vendor/code.google.com/p/go.net/html"
	"io/ioutil"
	"net/http"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.HandleFunc("/index", home)
	//FunctionMap holds a map of function names as they appear on the html and maps to the real function to call
	arriba.FunctionMap.Lock()
	arriba.FunctionMap.M["ChangeName"] = ChangeName
	arriba.FunctionMap.Unlock()
	http.ListenAndServe(":7070", nil)

}

func home(rw http.ResponseWriter, req *http.Request) {
	t, err := ioutil.ReadFile("index.html")
	if err != nil {
		panic(err)
	}

	value := arriba.Process(string(t))
	rw.Header().Add("Content-Type", "text/html; charset=UTF-8")
	rw.Write([]byte(value))

}

func ChangeName(node *html.Node) *html.Node {
	tree := h5.NewTree(node)
	t := transform.New(&tree)
	replacement := h5.Text("Hayley")
	t.Apply(transform.Replace(replacement), "p")
	return t.Doc()
}
