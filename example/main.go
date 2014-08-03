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

const html3 = (`<html><body><div data-lift="ChangeName">PEPE1<p name="name2">Diego</p></div>PIPI1<p data-lift="ChangeLastName">zzzzzzzzzzzz</p></body></html>`)
/* *
<html>
	<body>
		<div data-lift="ChangeName">
			PEPE
			<p name="name2">xxxxxxx</p>
		</div>
		PIPI
		<p data-lift="ChangeLastName">zzzzzzzzzzzz</p>
	</body>
</html>
 */
func Children(tokenArray []xml.Token, initialIndex int) ([]xml.Token, int, *Stack) {
	openTags := 1
	closedTags := 0
	stack := new(Stack)
	childStackIndex := new(Stack)
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
				break
			}

		case xml.StartElement:
			if (openTags - closedTags) == 1 {
				childStackIndex.Push(i)
			}
			openTags++

		case xml.CharData:
			if (openTags - closedTags) == 1 {
				fmt.Println("CHAR CHILDREN " + string(innerTok))
				stack.Push(xml.CopyToken(innerTok))
			}
		}
	}
	var ret = make([]xml.Token, stack.Size())
	var i = stack.Size()
	for i > 0 {
		ret[i - 1] = xml.CopyToken(stack.Pop())
		i--
	}
	return ret, closingTag, childStackIndex
}

func toTokenArray(decoder *xml.Decoder) []xml.Token {
	var i = 0
	var array = make([]xml.Token, 100) // todo FIX THIS
	for {
		tok, err := decoder.Token()
		if err != nil {
			//We are done processing tokens, let's end.
			break
		}

		switch innerTok := tok.(type) {
		case xml.StartElement:
			array[i] = innerTok
		case xml.CharData:
			array[i] = xml.CopyToken(innerTok)
		case xml.EndElement:
			array[i] = innerTok
		}

		fmt.Print("--> ")
		fmt.Println(tok)
		i++
	}
	return array
}

type Stackednode struct {
	NodePosition int
	ClosePosition int
	level int
	Visited bool // this field could be removed and use the closeposition to check if it is -1
}

type StackedResult struct {
	str string
	level int
}

func StackedResultToString(stack *Stack, myLevel int) string {
	var result = ""
	for stack.Size() > 0 && stack.Head().(StackedResult).level > myLevel {
		//		fmt.Print("Comparing level: ")
		//		fmt.Print(stack.Head().(StackedResult).level)
		//		fmt.Print("with level: ")
		//		fmt.Print(myLevel)
		result = stack.Pop().(StackedResult).str+result
	}
	return result
}

func processTag(token xml.Token, visited bool, resultStack *Stack, level int) string {
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
				lvl := resultStack.Head().(StackedResult).level
				if (visited && resultStack.Size() > 0 && lvl > level) {
					childNodesStr = StackedResultToString(resultStack, level)
				}
				// get the function from the map
				var funcResult = ChangeName(childNodesStr) // TODO get the function from the map
				resultStack.Push(StackedResult{funcResult, lvl})
			}
		}
		stringTag = stringTag + ">"
	case xml.CharData:
		stringTag = stringTag + string(innerTok)
	case xml.EndElement:
		stringTag = stringTag+ "</" + innerTok.Name.Local + ">"
	}
	return stringTag
}

func Princ() {

	decoder := xml.NewDecoder(bytes.NewBufferString(html3))
	var tokenArray = toTokenArray(decoder)
	resultStack := new(Stack)
	itStack := new(Stack)

	itStack.Push(Stackednode{0, -1, 0, false})
	var currentNode Stackednode

	for itStack.Size() > 0 {
		currentNode = itStack.Pop().(Stackednode)
		switch innerTok := tokenArray[currentNode.NodePosition].(type) {
		case xml.EndElement:
			fmt.Println("")
		case xml.CharData:
			var processed = string(innerTok)
			resultStack.Push(StackedResult{processed, currentNode.level })
		default:
			fmt.Println(innerTok)
			if !currentNode.Visited {
				var children, last, childrenIndexes = Children(tokenArray, currentNode.NodePosition + 1)
				currentNode.ClosePosition = last
				currentNode.Visited = true
				// Agregar aca el children
				for _,child := range children {
					fmt.Print("Adding child char " + string(child.(xml.CharData)) + " for ")
					fmt.Println(currentNode.level)
					resultStack.Push(StackedResult{string(child.(xml.CharData)), currentNode.level + 1})
				}

				if childrenIndexes.Size() > 0 {
					itStack.Push(currentNode)
					for childrenIndexes.Size() > 0 {
						index := childrenIndexes.Pop().(int)
						fmt.Print("Adding ")
						fmt.Println(tokenArray[index])
						fmt.Print("Level ")
						fmt.Println(currentNode.level)
						itStack.Push(Stackednode{index, -1, currentNode.level + 1, false})
					}
				} else {
					// traverse the attributes of the node and apply the function specified by data-lift attr
					// save the string version of the node in result slice with the closing tag
					var processed = processTag(tokenArray[currentNode.NodePosition], currentNode.Visited, resultStack, currentNode.level)
					processedChildren := ""
	//				if resultStack.Size() > 0 && resultStack.Head().(StackedResult).level > currentNode.level {
						processedChildren = StackedResultToString(resultStack, currentNode.level)
	//				}
					resultStack.Push(StackedResult{processed + processedChildren + processTag(tokenArray[currentNode.ClosePosition], false, resultStack, currentNode.level), currentNode.level})
					fmt.Println("Result stackv: " + resultStack.Head().(StackedResult).str)
				}
			} else {
				// traverse the attributes of the node and apply the function specified by data-lift attr
				// save the string version of the node in result slice with the closing tag
				var processed = processTag(tokenArray[currentNode.NodePosition], currentNode.Visited, resultStack, currentNode.level)
				processedChildren := ""
	//			if resultStack.Size() > 0 && resultStack.Head().(StackedResult).level > currentNode.level {
					processedChildren = StackedResultToString(resultStack, currentNode.level)
	//			}
				resultStack.Push(StackedResult{processed + processedChildren + processTag(tokenArray[currentNode.ClosePosition], false, resultStack, currentNode.level), currentNode.level})
				fmt.Print("Result stack/: ")
				fmt.Println(resultStack.Head().(StackedResult))
			}
		}
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
}

func ChangeName(html string) string {
	return strings.Replace(html, "Diego", "Gabriel", 1)
}
