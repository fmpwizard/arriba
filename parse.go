package arriba

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fmpwizard/arriba/vendor/code.google.com/p/go-html-transform/css/selector"
	"github.com/fmpwizard/arriba/vendor/code.google.com/p/go-html-transform/h5"
	"github.com/fmpwizard/arriba/vendor/code.google.com/p/go-html-transform/html/transform"
	"github.com/fmpwizard/arriba/vendor/code.google.com/p/go.net/html"
	"sync"
)

type HTMLTransform func(*html.Node) *html.Node

var FunctionMap = struct {
	sync.RWMutex
	M map[string]HTMLTransform
}{M: make(map[string]HTMLTransform)}

func Process2(in []byte) []byte {
	tree, _ := h5.New(bytes.NewReader(in))
	t := transform.New(tree)
	functionsInScope, _ := selector.Selector("[data-lift]")
	snippetNodess := functionsInScope.Find(tree.Top())
	for _, snippet := range snippetNodess {
		for _, function := range snippet.Attr {
			if function.Key == "data-lift" {
				replacement, err := do(function.Val, snippet)
				if err == nil {
					buf := bytes.NewBufferString("")
					html.Render(buf, replacement)
					t.Apply(transform.Replace(replacement), "[data-lift="+function.Val+"]")
				} else {
					fmt.Println("ERROR: " + err.Error())
				}
			}
		}
	}
	buf := bytes.NewBufferString("")
	html.Render(buf, t.Doc())
	return buf.Bytes()
}

func do(scopeFunction string, snippetHTML *html.Node) (*html.Node, error) {
	FunctionMap.RLock()
	f, found := FunctionMap.M[scopeFunction]
	FunctionMap.RUnlock()
	if found {
		return f(snippetHTML), nil
	} else {
		return &html.Node{}, errors.New("Did not find function: '" + scopeFunction + "'")
	}
}
