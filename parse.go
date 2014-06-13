package arriba

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
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

func MarshallElem(in string) string {
	funcMap["ChangeTime"] = ChangeTime
	funcMap["ChangeName"] = ChangeName
	completeHTML := ""
	parentTag := ""
	snipetHTML := ""
	//node := Elem{}
	//err := xml.Unmarshal([]byte(in), &node)

	decoder := xml.NewDecoder(bytes.NewBufferString(in))

	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.StartElement:

			if completeHTML != "" {
				completeHTML = completeHTML + "<" + startElement.Name.Local
			} else {
				completeHTML = "<" + startElement.Name.Local
			}
			//parentTag = "<2" + startElement.Name.Local

			open := 1
			close := 0
		Loop:
			for _, value := range startElement.Attr {
				//fmt.Printf("Att: %v ==> value: %v\n", value.Name.Local, value.Value)
				parentTag = parentTag + " " + value.Name.Local + "=\"" + value.Value + "\""
				if value.Name.Local == "data-lift" {

					for {
						tok, err := decoder.Token()
						if err != nil {
							return err.Error()
						}
						switch innerTok := tok.(type) {
						case xml.StartElement:
							if snipetHTML == "" {
								snipetHTML = parentTag + ">" //we found first inner node, so close the parent
							}

							snipetHTML = snipetHTML + "<" + innerTok.Name.Local
							for _, attr := range innerTok.Attr {
								snipetHTML = snipetHTML + " " + attr.Name.Local + "=\"" + attr.Value + "\""
							}
							snipetHTML = snipetHTML + ">"
							open++
							//fmt.Printf("1111: %v\n", innerTok.Name.Local)
						case xml.CharData:
							snipetHTML = snipetHTML + string(innerTok)

						case xml.EndElement:
							//fmt.Printf("2222 %v\n", innerTok.Name.Local)
							snipetHTML = snipetHTML + "</" + innerTok.Name.Local + ">"
							close++
							if open == close {
								fmt.Printf("End ==>> %v, %v\n", open, close)
								//return snipetHTML
								completeHTML = completeHTML + snipetHTML
								break Loop
							}
						}
						fmt.Printf(" ==>> snipetHTML  %v\n\n", snipetHTML)
					}

				}
			}

			fmt.Printf("Start: %v\n", startElement.Name.Local)
			//fmt.Printf("Debug: %v\n", snipetHTML)
			if snipetHTML == "" {
				completeHTML = completeHTML + ">"
			}

			fmt.Printf("completeHTML: %v\n", completeHTML)
		case xml.CharData:
			fmt.Printf("\n\nCharData: %+v\n", string(startElement))
		case xml.EndElement:
			fmt.Printf("\n\nEnd: %+v\n", startElement.Name.Local)
			completeHTML = completeHTML + "</" + startElement.Name.Local + ">"
		case xml.Comment:
			fmt.Printf("\n\nComment: %+v\n", startElement)
		case xml.Directive:
			fmt.Printf("\n\nDirective: %+v\n", string(startElement))
		case xml.Token:
			fmt.Printf("\n\n4: %+v\n", startElement)

		default:
			fmt.Printf("\n=========$$$$$$$$$$$$$$$$\n")
		}

	}
	fmt.Printf("completeHTML %v \n", completeHTML)
	return completeHTML
}

func ChangeTime(html string) string {
	return strings.Replace(html, "Time goes here", time.Now().Format("2006-01-02T15:04:05.999999999Z07:00"), 1)
}

func ChangeName(html string) string {
	return strings.Replace(html, "Diego", "Gabriel", 1)
}
