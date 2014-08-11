package arriba

import (
	"github.com/fmpwizard/arriba/vendor/code.google.com/p/go-html-transform/h5"
	"github.com/fmpwizard/arriba/vendor/code.google.com/p/go-html-transform/html/transform"
	"github.com/fmpwizard/arriba/vendor/code.google.com/p/go.net/html"
	"testing"
)

func init() {
	FunctionMap.Lock()
	FunctionMap.M["ChangeName"] = ChangeName
	FunctionMap.M["ChangeLastName"] = ChangeLastName
	FunctionMap.Unlock()
}

func TestMarshallElemDifferentSnippets(t *testing.T) {
	res := Process([]byte(html1))
	if string(res) != html1Expected {
		t.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html1Expected, string(res))
	}
}

/*

func TestMarshallElemSingleSnippet(t *testing.T) {
	res := MarshallElem(html2)
	if res != html2Expected {
		t.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html2Expected, res)
	}
}

func TestMarshallElemNestedSnippet(t *testing.T) {
	res := MarshallElem(html3)
	if res != html3Expected {
		t.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html3Expected, res)
	}
}

func TestMarshallElemMultipleSnippetSameLevel(t *testing.T) {
	res := MarshallElem(html4)
	if res != html4Expected {
		t.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html4Expected, res)
	}
}

func TestMarshallElemSnippetNotFound(t *testing.T) {
	res := MarshallElem(html5)
	if res != html5Expected {
		t.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html5Expected, res)
	}
}

func TestMarshallUntouchedStrings(t *testing.T) {
	res := MarshallElem(html6)
	if res != html6Expected {
		t.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html6Expected, res)
	}
}

func TestMarshallUntouchedStringsAfterFunction(t *testing.T) {
	res := MarshallElem(html7)
	if res != html7Expected {
		t.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html7Expected, res)
	}
}

func TestMarshallMultipleComplexAttributes(t *testing.T) {
	res := MarshallElem(html8)
	if res != html8Expected {
		t.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html8Expected, res)
	}
}
*/
/*func TestMarshallHtml5Transform(t *testing.T) {
	res := ReplaceInnerSpan(html9)
	if res != html9Expected {
		t.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html9Expected, res)
	}
}
*/
/*func BenchmarkMarshallElemDifferentSnippets(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := MarshallElem(html1)
		if res != html1Expected {
			b.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html1Expected, res)
		}
	}
}

func BenchmarkMarshallElemSingleSnippet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := MarshallElem(html2)
		if res != html2Expected {
			b.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html2Expected, res)
		}
	}
}

func BenchmarkTestMarshallElemNestedSnippet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := MarshallElem(html3)
		if res != html3Expected {
			b.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html3Expected, res)
		}
	}
}

func BenchmarkMarshallElemMultipleSnippetSameLevel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := MarshallElem(html4)
		if res != html4Expected {
			b.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html4Expected, res)
		}
	}
}

func BenchmarkMarshallElemSnippetNotFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := MarshallElem(html5)
		if res != html5Expected {
			b.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html5Expected, res)
		}
	}
}*/

const html1 = (`<html><head></head><body><div data-lift="ChangeName"><p name="name">Diego</p><p data-lift="ChangeLastName">Medina</p></div></body></html>`)
const html1Expected = (`<html><head></head><body><div><p name="name">Gabriel</p><p>Bauman</p></div></body></html>`)

const html2 = (`<html><head></head><body><div data-lift="ChangeName"><p name="name">Diego</p><p class="pretty-last-name">Medina</p></div></body></html>`)
const html2Expected = (`<html><head></head><body><div><p name="name">Gabriel</p><p class="pretty-last-name">Medina</p></div></body></html>`)

const html3 = (`<html><head></head><body><div data-lift="ChangeName"><p name="name">Diego1</p><div data-lift="ChangeName"><p name="name">Diego</p></div></div></body></html>`)
const html3Expected = (`<html><head></head><body><div><p name="name">Gabriel1</p><div><p name="name">Gabriel</p></div></div></body></html>`)

const html4 = (`<html><head></head><body><div data-lift="ChangeName"><p name="name">Diego</p><p class="pretty-last-name">Medina</p></div><div data-lift="ChangeName"><p name="name">Diego1</p><p class="pretty-last-name">Medina</p></div></body></html>`)
const html4Expected = (`<html><head></head><body><div><p name="name">Gabriel</p><p class="pretty-last-name">Medina</p></div><div><p name="name">Gabriel1</p><p class="pretty-last-name">Medina</p></div></body></html>`)

const html5 = (`<html><head></head><body><div data-lift="DoesNotExist"><p name="name">Diego</p></div></body></html>`)
const html5Expected = (`Did not find function: 'DoesNotExist'`)

const html6 = (`<html><head></head><body><div><p name="name">Diego</p></div></body></html>`)
const html6Expected = (`<html><head></head><body><div><p name="name">Diego</p></div></body></html>`)

const html7 = (`<div><p><span data-lift="ChangeLastName">Medina</span></p><p>Here is some random string nobody changed.</p></div>`)
const html7Expected = (`<div><p><span>Bauman</span></p><p>Here is some random string nobody changed.</p></div>`)

const html8 = (`<meta http-equiv="X-UA-Compatible" content="IE=Edge"></meta>`)
const html8Expected = (`<meta http-equiv="X-UA-Compatible" content="IE=Edge"></meta>`)

const html9 = (`<html><head></head><body><div data-lift="ReplaceInnerSpan"><p>Diego</p><p class="last-name">Bauman</p></div></body></html>`)
const html9Expected = (`<html><head></head><body><div><p>Diego</p><p class="last-name">Medina</p></div></body></html>`)

func ChangeName(node *html.Node) *html.Node {
	tree := h5.NewTree(node)
	t := transform.New(&tree)
	replacement := h5.Text("Hayley")
	t.Apply(transform.ReplaceChildren(replacement), "name=name")
	return t.Doc()
}

func ChangeLastName(node *html.Node) *html.Node {
	tree := h5.NewTree(node)
	t := transform.New(&tree)
	replacement := h5.Text("Bauman")
	t.Apply(transform.ReplaceChildren(replacement), "p")
	return t.Doc()
}
