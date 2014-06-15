package arriba

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

type HTMLTransform func(string) string

type Elem struct {
	XMLName  xml.Name
	Comment  xml.Comment
	Attr     xml.Attr
	InnerXML string `xml:",innerxml"`
	Kids     []Elem `xml:",any"`
}

var funcMap = make(map[string]HTMLTransform)
var ch = make(chan string)

func MarshallElem(in string) string {
	fmt.Println("\n\n\n\n\n")

	go readChan(ch)
	funcMap["ChangeName"] = ChangeName
	completeHTML := ""

	decoder := xml.NewDecoder(bytes.NewBufferString(in))

	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch element := token.(type) {
		case xml.StartElement:

			if completeHTML != "" {
				completeHTML = completeHTML + "<" + element.Name.Local
			} else {
				completeHTML = "<" + element.Name.Local
			}
			functionName := ""
			for _, value := range element.Attr {
				_, res := processSnippet(value, decoder)
				completeHTML = completeHTML + res

				//}
			}
			if !strings.HasSuffix(completeHTML, ">") {
				completeHTML = completeHTML + ">"
			}

			if functionName != "" {
				fmt.Printf("functionName: %v\n", functionName)
			}
		case xml.CharData:
			fmt.Printf("CharData: %+v\n", string(element))
		case xml.EndElement:
			completeHTML = completeHTML + "</" + element.Name.Local + ">"
		case xml.Comment:
			fmt.Printf("Comment: %+v\n", element)
		case xml.Directive:
			fmt.Printf("Directive: %+v\n", string(element))
		case xml.Token:
			fmt.Printf("4: %+v\n", element)

		default:
			fmt.Errorf("\nIf yo uare here, you are missing a type: %v\n", element)
		}

	}
	return completeHTML
}

func processSnippet(value xml.Attr, decoder *xml.Decoder) (error, string) {

	//functionName := value.Value

	snippetHTML := ""
	open := 1
	closingTags := 0

	//parentTag = "<" + parentTag + " " + value.Name.Local + "=\"" + value.Value + "\""
	//parentTag = parentTag + " - " + value.Name.Local + "=\"" + value.Value + "\""
	parentTag := ""
	if value.Name.Local != "data-lift" {
		parentTag = " " + value.Name.Local + "=\"" + value.Value + "\""
	}
	if value.Name.Local == "data-lift" {
		//fmt.Printf("111\n")
		for {
			fmt.Println("parentTag " + parentTag)
			tok, err := decoder.Token()
			if err != nil {
				//We are done processing tokens, let's end.
				close(ch)
				return err, snippetHTML
			}
			switch innerTok := tok.(type) {
			case xml.StartElement:
				if snippetHTML == "" && parentTag != "" && !strings.HasSuffix(parentTag, ">") {
					//if snippetHTML == "" {

					snippetHTML = parentTag + "/>" //we found first inner node, so close the parent
				} else if snippetHTML == "" && parentTag == "" && !strings.HasSuffix(parentTag, ">") {

					snippetHTML = parentTag + "ss>"
				}

				snippetHTML = snippetHTML + "<" + innerTok.Name.Local
				fmt.Printf("2========== snippetHTML %v\n", snippetHTML)
				//fmt.Println("1")
				for _, attr := range innerTok.Attr {
					if attr.Name.Local != "data-lift" {
						snippetHTML = snippetHTML + " " + attr.Name.Local + "=\"" + attr.Value + "\""
					}

					ch <- snippetHTML
					_, super := processSnippet(attr, decoder)
					//_, super := processSnippet(attr, decoder, innerTok.Name.Local)
					//fmt.Printf("snippetHTML > super %v>%v\n", snippetHTML, super)
					if strings.HasSuffix(snippetHTML, ">") {

						snippetHTML = snippetHTML + super
					} else {
						fmt.Println("1 " + super)
						snippetHTML = snippetHTML + ">" + super
					}
					ch <- snippetHTML
				}
				//fmt.Printf("1 snippetHTML is %v\n", snippetHTML)
				//snippetHTML = snippetHTML + ">============"
				//fmt.Println("2")
				open++
			case xml.CharData:
				snippetHTML = snippetHTML + string(innerTok)
			case xml.EndElement:
				snippetHTML = snippetHTML + "</" + innerTok.Name.Local + ">"
				//fmt.Printf(" ==>> snippetHTML  %v\n", snippetHTML)
				closingTags++
				//fmt.Printf("Open: %v, closing tag: %v\n", open, closingTags)
				if open == closingTags { //do we have our matching closing tag? //This fails with autoclose tags I think
					//fmt.Printf("2 snippetHTML is %v\n", snippetHTML)
					//fmt.Printf("3 %v\n", snippetHTML)
					//fmt.Printf("33 %v\n", ChangeName(snippetHTML))
					ch <- snippetHTML
					return nil, ChangeName(snippetHTML)
				}
			}
			//fmt.Printf(" ==>> snippetHTML  %v\n", snippetHTML)
		}

	}
	return nil, ""
}

func ChangeName(html string) string {
	return strings.Replace(html, "Diego", "Gabriel", 1)
}

func readChan(ch chan string) {
	var buffer string
	for {
		select {
		case data, ok := <-ch:
			if ok == true {
				fmt.Printf("got: %v\n", data)
				buffer = data
				//buffer =  buffer + data
			} else {
				fmt.Printf("sending: %v\n", buffer)
				return
			}
		}
	}

}
