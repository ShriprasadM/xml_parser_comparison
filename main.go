package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/beevik/etree"
)

const (
	trackingEventsTag         = "TrackingEvents"
	trackingEventsTagOpenTag  = "<" + trackingEventsTag + ">"
	trackingEventsTagCloseTag = "</" + trackingEventsTag + ">"
	SampleTrackingEvent       = "<Tracking event=\"shriprasad\"><![CDATA[https://mytracker.com]]></Tracking>"
	linearEndTag              = "</Linear>"
	trackerUrl                = "https://mytracker.com"
)

func stringBased(vast string) (string, error) {
	ci := strings.Index(vast, trackingEventsTagCloseTag)

	// no impression tag - pass it as it is
	if ci == -1 {
		// return vast, nil
		// append tracking event
		// oi := strings.Index(vast, trackingEventsTagOpenTag)

		oi := strings.Index(vast, linearEndTag)

		// if ci-oi == len(trackingEventsTagOpenTag) {
		if oi != -1 {
			return strings.Replace(vast, linearEndTag, trackingEventsTagOpenTag+SampleTrackingEvent+trackingEventsTagCloseTag+linearEndTag, 1), nil
		}
		return vast, nil // single replacement
	}

	return strings.Replace(vast, trackingEventsTagCloseTag, SampleTrackingEvent+trackingEventsTagCloseTag, 1), nil
}

func etreeBased(vast string) (string, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(vast); err != nil {
		panic(err)
	}

	// ele := doc.SelectElement("impressionTag")

	ele := doc.FindElement("VAST/Ad/InLine/Creatives/Creative/Linear/" + trackingEventsTag)
	newEle := doc.CreateElement("Tracking")
	newEle.CreateAttr("event", "shriprasad")
	newEle.SetCData(trackerUrl)
	ele.AddChild(newEle)
	return doc.WriteToString()
}

func xmlEncodingBased(vast string) (string, error) {
	type Tracking struct {
		Name xml.Name `xml:"Tracking"`
		Attr xml.Attr
	}

	decoder := xml.NewDecoder(strings.NewReader(vast))
	var buf bytes.Buffer
	encoder := xml.NewEncoder(&buf)
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error getting token: %v\n", err)
			break
		}

		switch v := token.(type) {
		case xml.EndElement:
			if v.Name.Local == trackingEventsTag {
				// var tracking Tracking
				// if err = decoder.DecodeElement(&tracking, &v); err != nil {
				// 	log.Fatal(err)
				// }
				// // encoder.Encode()

				se := xml.StartElement{
					Name: xml.Name{
						Local: "Tracking",
					},
					Attr: []xml.Attr{
						{
							Name:  xml.Name{Local: "event"},
							Value: "shriprasad",
						},
					},
				}

				ee := xml.EndElement{
					Name: xml.Name{Local: "Tracking"},
				}

				if err := encoder.EncodeToken(se); err != nil {
					log.Fatal(err)
				}

				bt := fmt.Sprintf("[CDATA[%s]]", trackerUrl)

				if err := encoder.EncodeToken(xml.Directive(bt)); err != nil {
					log.Fatal(err)
				}

				if err := encoder.EncodeToken(ee); err != nil {
					log.Fatal(err)
				}
				if err := encoder.EncodeToken(xml.CopyToken(v)); err != nil {
					log.Fatal(err)
				}

				continue
			}
		}

		if err := encoder.EncodeToken(xml.CopyToken(token)); err != nil {
			log.Fatal(err)
		}

	}

	// must call flush, otherwise some elements will be missing
	if err := encoder.Flush(); err != nil {
		log.Fatal(err)
	}

	// fmt.Println(buf.String())

	return buf.String(), nil
}
