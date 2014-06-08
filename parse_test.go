package template

import (
	"testing"
)

func TestGetFunctions(*testing.T) {
	GetFunctions(html1)
}

const html1 = (`
<!DOCTYPE html>
<html>
  <head >
    <meta content="text/html; charset=UTF-8" http-equiv="content-type" />
    <title>Home</title>
  </head>
  <body>
    <div data-lift="surround?with=default;at=content">
      <h2>Welcome to your project!</h2>
      <p><span data-lift="helloWorld.howdy">Welcome to your Lift app at <span id="time">Time goes here</span></span></p>
    </div>
  </body>
</html>
`)

const html2 = (`<html><head></head><body><div data-lift="goFunctionName"><p name="name">Diego</p><p class="pretty-last-name">Medina</p></div></body></html>`)
const html3 = (`<span data-lift="helloWorld.howdy">Welcome to your Lift app at <span id="time">Time goes here</span></span>`)
