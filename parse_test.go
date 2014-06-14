package arriba

import (
	"testing"
)

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

/*func TestMarshallElemNestedSnippet(t *testing.T) {
	res := MarshallElem(html3)
	if res != html3Expected {
		t.Errorf("Got a different html, expeted: \n%v\n but got:\n%v\n", html3Expected, res)
	}
}
*/
const html1 = (`
<!DOCTYPE html>
<html>
  <head >
    <meta content="text/html; charset=UTF-8" http-equiv="content-type" />
    <title>Home</title>
  </head>
  <body>
    <div>
      <h2>Welcome to your project!</h2>
      <p><span data-lift="ChangeTime">Welcome to your Lift app at <span id="time">Time goes here</span></span></p>
    </div>
  </body>
</html>
`)

const html2 = (`<html><head></head><body><div data-lift="ChangeName"><p name="name">Diego</p><p class="pretty-last-name">Medina</p></div></body></html>`)
const html2Expected = (`<html><head></head><body><div data-lift="ChangeName"><p name="name">Gabriel</p><p class="pretty-last-name">Medina</p></div></body></html>`)

const html3 = (`<html><head></head><body><div data-lift="ChangeName"><p name="name">Diego</p><div data-lift="ChangeName"><p name="name">Diego</p></div></div></body></html>`)
const html3Expected = (`<html><head></head><body><div data-lift="ChangeName"><p name="name">Gabriel</p><div data-lift="ChangeName"><p name="name">Gabriel</p></div></div></body></html>`)

const html4 = (`<html><head></head><body><div data-lift="ChangeName"><p name="name">Diego</p><p class="pretty-last-name">Medina</p></div><div data-lift="ChangeName"><p name="name">Diego1</p><p class="pretty-last-name">Medina</p></div></body></html>`)
const html4Expected = (`<html><head></head><body><div data-lift="ChangeName"><p name="name">Gabriel</p><p class="pretty-last-name">Medina</p></div><div data-lift="ChangeName"><p name="name">Gabriel1</p><p class="pretty-last-name">Medina</p></div></body></html>`)
