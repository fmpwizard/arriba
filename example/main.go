package main

//This is a sample application that shows how to use arriba
//cd into the example folder and run go run main.go

import (
	"github.com/fmpwizard/arriba"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"encoding/xml"
	"bytes"
	"fmt"
)

var funcMap = make(map[string]HTMLTransform)

func main() {
	http.HandleFunc("/index", home)
	Princ()
	//funcMap holds a map of function names as they appear on the html and maps to the real function to call
	funcMap["ChangeTime"] = ChangeTime
	http.ListenAndServe(":7070", nil)
}

func home(rw http.ResponseWriter, req *http.Request) {
	t, err := ioutil.ReadFile("index.html")
	if err != nil {
		panic(err)
	}

	for functionName, html := range arriba.GetFunctions(string(t)) {
		//This is a silly way to replace the old html with new one, because
		//it will fail if we have the same raw html more than once.
		t = []byte(strings.Replace(string(t), html, funcMap[functionName](html), 1))
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


type Item struct {
	value interface{}
	next *Item // Next stack item
}

type Stack struct {
	top *Item // Top item of the stack
	size int  // item count of the stack
}

// Put the item on top of the stack
func (stack *Stack) Push(value interface{}) {
	stack.top = &Item { value, stack.top }
	stack.size++
}

// Put the items on top of the stack
func (stack *Stack) PushArray(value []interface{}) {
	for item := range value {
		stack.top = &Item { item, stack.top }
		stack.size++
	}
}

// If the stack is not empty, remove the top element and return the value
// If the stack is empty, return nil
func (stack *Stack) Pop() (value interface{}) {
	if stack.size > 0 {
		value, stack.top = stack.top.value, stack.top.next
		stack.size--
		return
	}
	return nil
}

// If the stack is not empty, remove the top element and return the value
// If the stack is empty, return nil
func (stack *Stack) Head() (value interface{}) {
	if stack.size > 0 {
		return stack.top.value
	}
	return nil
}

// item count of the stack
func (stack *Stack) Size() int {
	return stack.size
}




const html1 = (`<html><body><div data-lift="ChangeName"><p name="name">Diego</p><p data-lift="ChangeLastName">Medina</p></div></body></html>`)

const html2 = (`<html><body><div data-lift="ChangeName"><p name="name1">xxxxxxx</p></div><p data-lift="ChangeLastName">zzzzzzzzzzzz</p></body></html>`)

func Children(tokenArray []xml.Token, initialIndex int) ([]xml.Token, int) {
	openTags := 1
	closedTags := 0
	stack := new(Stack)
	var closingTag int
	for i:= initialIndex; i < len(tokenArray); i++ {
		if openTags == closedTags {
			break
		}
		tok := tokenArray[i]
		switch innerTok := tok.(type) {
		case xml.EndElement:
			closedTags++
			if openTags == closedTags {
				closingTag = i
				fmt.Print("hizo break en")
				fmt.Println(closingTag)
				break
			}
			stack.Push(innerTok)

		case xml.StartElement:
			openTags++
			stack.Push(innerTok)

		case xml.CharData:
			stack.Push(innerTok)
		}
	}
	var ret = make([]xml.Token, stack.Size())
	var i = stack.Size()
	for i > 0 {
		ret[i - 1] = stack.Pop()
		fmt.Print("metido al stack")
		fmt.Println(ret[i - 1])
		i--
	}
	return ret, closingTag
}

func toTokenArray(decoder *xml.Decoder) []xml.Token {
	stack := new(Stack)
	for {
		tok, err := decoder.Token()
		if err != nil {
			//We are done processing tokens, let's end.
			break
		}
		stack.Push(tok)
	}
	var array = make([]xml.Token, stack.Size())
	var i = stack.Size()
	for i > 0 {
		array[i - 1] = stack.Pop()
		i--
	}
	return array
}

type Stackednode struct {
	NodePosition int
	ClosePosition int
	Visited bool // this field could be removed and use the closeposition to check if it is -1
}

func processTag(token xml.Token, visited bool, resultStack *Stack) string {
	stringTag := ""
	switch innerTok := token.(type) {
	case xml.StartElement:
		stringTag = "<"+innerTok.Name.Local
		for _, attr := range innerTok.Attr {
			if attr.Name.Local != "data-lift" {
				stringTag = stringTag+" "+attr.Name.Local+"=\""+attr.Value+"\""
			} else {
				// call the function
				childNodesStr := ""
				if (visited && resultStack.Size() > 0) {
					childNodesStr = resultStack.Pop().(string)
				}
				fmt.Println("child debug " + childNodesStr)
				// get the function from the map
				var funcResult = ChangeName(childNodesStr) // TODO get the function from the map
				resultStack.Push(funcResult)
			}
		}
		stringTag = stringTag + ">"
	case xml.CharData:
		fmt.Println("char data " + string(innerTok))
		stringTag = stringTag + string(innerTok)
	case xml.EndElement:
		stringTag = stringTag+ "</" + innerTok.Name.Local + ">"
	}
	return stringTag
}

func Princ() {

	decoder := xml.NewDecoder(bytes.NewBufferString(html2))

	decoder2 := xml.NewDecoder(bytes.NewBufferString(html2))
	fmt.Println("Array")
	var tokenArrayp = toTokenArray(decoder2)
	for _, item := range tokenArrayp {
		switch innerTok := item.(type) {
		case xml.StartElement:
			fmt.Println(innerTok)
		case xml.CharData:
			fmt.Println(string(innerTok))
		case xml.EndElement:
			fmt.Println(innerTok)
		}
	}

	var tokenArray = toTokenArray(decoder)

	resultStack := new(Stack)
	itStack := new(Stack)

	itStack.Push(Stackednode{0, -1, false})
	var currentNode Stackednode

	for itStack.Size() > 0 {
		currentNode = itStack.Pop().(Stackednode)
//		fmt.Print("Visiting ")
//		fmt.Print(itStack.Size())
		switch innerTok := tokenArray[currentNode.NodePosition].(type) {
		case xml.EndElement:
			fmt.Println("")
		case xml.CharData:
			var processed = processTag(tokenArray[currentNode.NodePosition], currentNode.Visited, resultStack)
//			processedChildren := ""
//			if resultStack.Size() > 0 {
//				processedChildren = resultStack.Pop().(string)
//			}
			resultStack.Push(processed)
		default:
			fmt.Println(innerTok)
			if !currentNode.Visited {
				var children, last = Children(tokenArray, currentNode.NodePosition + 1)
				fmt.Println("children")
				fmt.Println(children)
				fmt.Println(last)
				fmt.Println("--")
				currentNode.ClosePosition = last
				currentNode.Visited = true
				if len(children) > 0 {
					itStack.Push(currentNode)
					itStack.Push(Stackednode{currentNode.NodePosition + 1, -1, false})
				} else {
					// traverse the attributes of the node and apply the function specified by data-lift attr
					// save the string version of the node in result slice with the closing tag
					var processed = processTag(tokenArray[currentNode.NodePosition], currentNode.Visited, resultStack)
					resultStack.Push(processed + processTag(tokenArray[currentNode.ClosePosition], false, resultStack))
					fmt.Println("Result stackv: " + resultStack.Head().(string))
				}
			} else {
				// traverse the attributes of the node and apply the function specified by data-lift attr
				// save the string version of the node in result slice with the closing tag
				var processed = processTag(tokenArray[currentNode.NodePosition], currentNode.Visited, resultStack)
				processedChildren := ""
				if resultStack.Size() > 0 {
					processedChildren = resultStack.Pop().(string)
				}
				resultStack.Push(processed + processedChildren + processTag(tokenArray[currentNode.ClosePosition], false, resultStack))
				fmt.Println("Result stack: " + resultStack.Head().(string))
			}
		}
	}

	//fmt.Println(tokenArray)

	//var children, _ = Children(tokenArray, 1)

	for resultStack.Size() > 0 {
		fmt.Println(resultStack.Pop().(string))
	}

	/*
	init result stack empty
	init nodes stack with root and the state not-visited

	while stack.size > 0
	 Get the xml node from the stack
	  if the node has children:
		 push the node to the stack marking it as visited and after that push the children nodes as not-visited
	  if the node has been visited before get the children from the result stack and apply the data-lift function if exists, pushing the transformation into the result stack
	  if the node has not children check if the node has the attributes (ie "data-lift") and execute the binded function pushing the transformation into the result stack

	return nodesStack.pop
	*/

	//fmt.Println(children)
	//
	//	token2, _ := decoder.Token()
	//
	//	fmt.Println(Children(&token2))



}

func ChangeName(html string) string {
	return strings.Replace(html, "Diego", "Gabriel", 1)
}
