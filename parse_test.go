package arriba

import (
	"testing"
)

func TestMarshallElemDifferentSnippets(t *testing.T) {
	res := MarshallElem(html1)
	if res != html1Expected {
		t.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html1Expected, res)
	}
}

func TestMarshallElemSingleSnippet(t *testing.T) {
	res := MarshallElem(html2)
	if res != html2Expected {
		t.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html2Expected, res)
	}
}

func TestMarshallElemMultipleSnippetSameLevel(t *testing.T) {
	res := MarshallElem(html4)
	if res != html4Expected {
		t.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html4Expected, res)
	}
}

func TestMarshallElemNestedSnippet(t *testing.T) {
	res := MarshallElem(html3)
	if res != html3Expected {
		t.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html3Expected, res)
	}
}

func TestMarshallElemSnippetNotFound(t *testing.T) {
	res := MarshallElem(html5)
	if res != html5Expected {
		t.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html5Expected, res)
	}
}

const html1 = (`<html><head></head><body><div data-lift="ChangeName"><p name="name">Diego</p><p data-lift="ChangeLastName">Medina</p></div></body></html>`)
const html1Expected = (`<html><head></head><body><div><p name="name">Gabriel</p><p>Bauman</p></div></body></html>`)

const html2 = (`<html><head></head><body><div data-lift="ChangeName"><p name="name">Diego</p><p class="pretty-last-name">Medina</p></div></body></html>`)
const html2Expected = (`<html><head></head><body><div><p name="name">Gabriel</p><p class="pretty-last-name">Medina</p></div></body></html>`)

const html3 = (`<html><head></head><body><div data-lift="ChangeName"><p name="name">Diego</p><div data-lift="ChangeName"><p name="name">Diego</p></div></div></body></html>`)
const html3Expected = (`<html><head></head><body><div><p name="name">Gabriel</p><div><p name="name">Gabriel</p></div></div></body></html>`)

const html4 = (`<html><head></head><body><div data-lift="ChangeName"><p name="name">Diego</p><p class="pretty-last-name">Medina</p></div><div data-lift="ChangeName"><p name="name">Diego1</p><p class="pretty-last-name">Medina</p></div></body></html>`)
const html4Expected = (`<html><head></head><body><div><p name="name">Gabriel</p><p class="pretty-last-name">Medina</p></div><div><p name="name">Gabriel1</p><p class="pretty-last-name">Medina</p></div></body></html>`)

const html5 = (`<html><head></head><body><div data-lift="DoesNotExist"><p name="name">Diego</p></div></body></html>`)
const html5Expected = (`Did not find function DoesNotExist`)
