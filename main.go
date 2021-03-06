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
	SampleTrackingEvent       = "<Tracking event=\"close\"><![CDATA[https://mytracker.com]]></Tracking>"
	linearEndTag              = "</Linear>"
	nonLinearEndTag           = "</NonLinearAds>"
	trackerUrl                = "https://mytracker.com"

	TrackingEvent = "<Tracking event=\"%s\"><![CDATA[%s]]></Tracking>"
)

type VastModifier struct {
	vast string
}

func (v *VastModifier) SetVast(vast string) {
	v.vast = vast
}

func (v *VastModifier) InjectTrackerEvent(event, trackerUrl string) {
	trackingEvent := fmt.Sprintf(TrackingEvent, event, trackerUrl)
	list := strings.SplitAfter(v.vast, "/Creative>")

	for i, cr := range list {
		ci := strings.Index(cr, trackingEventsTagCloseTag)

		if ci == -1 {

			li := strings.Index(cr, linearEndTag)
			if li != -1 {
				list[i] = strings.Replace(cr, linearEndTag, trackingEventsTagOpenTag+trackingEvent+trackingEventsTagCloseTag+linearEndTag, 1)
			}

			nli := strings.Index(cr, nonLinearEndTag)
			if nli != -1 {
				list[i] = strings.Replace(cr, nonLinearEndTag, trackingEventsTagOpenTag+trackingEvent+trackingEventsTagCloseTag+nonLinearEndTag, 1)
			}
		}
		list[i] = strings.Replace(cr, trackingEventsTagCloseTag, trackingEvent+trackingEventsTagCloseTag, 1)
	}

	v.vast = strings.Join(list, "")
}

func (v *VastModifier) ToString() string {
	return v.vast
}

func stringBased(vast string) string {
	vm := VastModifier{}
	vm.SetVast(vast)
	vm.InjectTrackerEvent("close", "https://mytracker.com")
	return vm.ToString()
}

// func stringBased(vast string) string {

// 	list := strings.SplitAfter(vast, "/Creative>")

// 	for i, cr := range list {
// 		ci := strings.Index(cr, trackingEventsTagCloseTag)

// 		if ci == -1 {

// 			li := strings.Index(cr, linearEndTag)
// 			if li != -1 {
// 				list[i] = strings.Replace(cr, linearEndTag, trackingEventsTagOpenTag+SampleTrackingEvent+trackingEventsTagCloseTag+linearEndTag, 1)
// 			}

// 			nli := strings.Index(cr, nonLinearEndTag)
// 			if nli != -1 {
// 				list[i] = strings.Replace(cr, nonLinearEndTag, trackingEventsTagOpenTag+SampleTrackingEvent+trackingEventsTagCloseTag+nonLinearEndTag, 1)
// 			}
// 		}
// 		list[i] = strings.Replace(cr, trackingEventsTagCloseTag, SampleTrackingEvent+trackingEventsTagCloseTag, 1)
// 	}

// 	return strings.Join(list, "")
// }

func etreeBased(vast string) (string, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(vast); err != nil {
		panic(err)
	}

	for _, ele := range doc.FindElements("VAST/Ad/InLine/Creatives/Creative") {
		trackEle := ele.FindElement("Linear/" + trackingEventsTag)
		if trackEle != nil {
			newEle := doc.CreateElement("Tracking")
			newEle.CreateAttr("event", "close")
			newEle.SetCData(trackerUrl)
			trackEle.AddChild(newEle)
		}

		trackEle = ele.FindElement("NonLinearAds/" + trackingEventsTag)
		if trackEle != nil {
			newEle := doc.CreateElement("Tracking")
			newEle.CreateAttr("event", "close")
			newEle.SetCData(trackerUrl)
			trackEle.AddChild(newEle)
		}

	}

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
							Value: "close",
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

func main() {
	// vast := stringBased(vast)
	//vast, _ := xmlEncodingBased(vast)
	vast, _ := etreeBased(vast)
	fmt.Println(vast)
}

var vast string = `
<VAST version="4.2" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns="http://www.iab.com/VAST">
  <Ad id="20001" >
    <InLine>
      <AdSystem version="1">iabtechlab</AdSystem>
      <Error><![CDATA[https://example.com/error]]></Error>
      <Impression id="Impression-ID"><![CDATA[https://example.com/track/impression]]></Impression>
      <AdServingId>a532d16d-4d7f-4440-bd29-2ec05553fc80</AdServingId>
      <AdTitle>Inline Simple Ad</AdTitle>
      <AdVerifications></AdVerifications>
      <Advertiser>IAB Sample Company</Advertiser>
      <Category authority="https://www.iabtechlab.com/categoryauthority">AD CONTENT description category</Category>
      <Creatives>
        <Creative id="5480" sequence="1" adId="2447226">
          <Linear>
            <TrackingEvents>
              <Tracking event="start" ><![CDATA[https://example.com/tracking/start]]></Tracking>
              <Tracking event="progress" offset="00:00:10"><![CDATA[http://example.com/tracking/progress-10]]></Tracking>
              <Tracking event="firstQuartile"><![CDATA[https://example.com/tracking/firstQuartile]]></Tracking>
              <Tracking event="midpoint"><![CDATA[https://example.com/tracking/midpoint]]></Tracking>
              <Tracking event="thirdQuartile"><![CDATA[https://example.com/tracking/thirdQuartile]]></Tracking>
              <Tracking event="complete"><![CDATA[https://example.com/tracking/complete]]></Tracking>
            </TrackingEvents>
            <Duration>00:00:16</Duration>
            <MediaFiles>
              <MediaFile id="5241" delivery="progressive" type="video/mp4" bitrate="2000" width="1280" height="720" minBitrate="1500" maxBitrate="2500" scalable="1" maintainAspectRatio="1" codec="H.264">
                <![CDATA[https://iab-publicfiles.s3.amazonaws.com/vast/VAST-4.0-Short-Intro.mp4]]>
              </MediaFile>
              <MediaFile id="5244" delivery="progressive" type="video/mp4" bitrate="1000" width="854" height="480" minBitrate="700" maxBitrate="1500" scalable="1" maintainAspectRatio="1" codec="H.264">
                <![CDATA[https://iab-publicfiles.s3.amazonaws.com/vast/VAST-4.0-Short-Intro-mid-resolution.mp4]]>
              </MediaFile>
              <MediaFile id="5246" delivery="progressive" type="video/mp4" bitrate="600" width="640" height="360" minBitrate="500" maxBitrate="700" scalable="1" maintainAspectRatio="1" codec="H.264">
                <![CDATA[https://iab-publicfiles.s3.amazonaws.com/vast/VAST-4.0-Short-Intro-low-resolution.mp4]]>
              </MediaFile>
            </MediaFiles>
            <VideoClicks>
              <ClickThrough id="blog">
                <![CDATA[https://iabtechlab.com]]>
              </ClickThrough>
            </VideoClicks>
          </Linear>
          <UniversalAdId idRegistry="Ad-ID">8465</UniversalAdId>
        </Creative>
		<Creative id="5480" sequence="1" adId="2447226">
          <Linear>
            <TrackingEvents>
              <Tracking event="start" ><![CDATA[https://example.com/tracking/start]]></Tracking>
              <Tracking event="progress" offset="00:00:10"><![CDATA[http://example.com/tracking/progress-10]]></Tracking>
              <Tracking event="firstQuartile"><![CDATA[https://example.com/tracking/firstQuartile]]></Tracking>
              <Tracking event="midpoint"><![CDATA[https://example.com/tracking/midpoint]]></Tracking>
              <Tracking event="thirdQuartile"><![CDATA[https://example.com/tracking/thirdQuartile]]></Tracking>
              <Tracking event="complete"><![CDATA[https://example.com/tracking/complete]]></Tracking>
            </TrackingEvents>
            <Duration>00:00:16</Duration>
            <MediaFiles>
              <MediaFile id="5241" delivery="progressive" type="video/mp4" bitrate="2000" width="1280" height="720" minBitrate="1500" maxBitrate="2500" scalable="1" maintainAspectRatio="1" codec="H.264">
                <![CDATA[https://iab-publicfiles.s3.amazonaws.com/vast/VAST-4.0-Short-Intro.mp4]]>
              </MediaFile>
              <MediaFile id="5244" delivery="progressive" type="video/mp4" bitrate="1000" width="854" height="480" minBitrate="700" maxBitrate="1500" scalable="1" maintainAspectRatio="1" codec="H.264">
                <![CDATA[https://iab-publicfiles.s3.amazonaws.com/vast/VAST-4.0-Short-Intro-mid-resolution.mp4]]>
              </MediaFile>
              <MediaFile id="5246" delivery="progressive" type="video/mp4" bitrate="600" width="640" height="360" minBitrate="500" maxBitrate="700" scalable="1" maintainAspectRatio="1" codec="H.264">
                <![CDATA[https://iab-publicfiles.s3.amazonaws.com/vast/VAST-4.0-Short-Intro-low-resolution.mp4]]>
              </MediaFile>
            </MediaFiles>
            <VideoClicks>
              <ClickThrough id="blog">
                <![CDATA[https://iabtechlab.com]]>
              </ClickThrough>
            </VideoClicks>
          </Linear>
          <UniversalAdId idRegistry="Ad-ID">8465</UniversalAdId>
        </Creative>
        <Creative id="5480" sequence="1" adId="2447226">
          <NonLinearAds>
            <NonLinear width="350" height="350">
              <StaticResource creativeType="image/png">
                <![CDATA[https://mms.businesswire.com/media/20150623005446/en/473787/21/iab_tech_lab.jpg]]>
              </StaticResource>
              <NonLinearClickThrough>
               <![CDATA[https://iabtechlab.com]]>
              </NonLinearClickThrough>
              <NonLinearClickTracking>
               <![CDATA[https://example.com/tracking/clickTracking]]>
              </NonLinearClickTracking>
            </NonLinear>
            <TrackingEvents>
            </TrackingEvents>
           </NonLinearAds>
           <UniversalAdId idRegistry="Ad-ID">8465</UniversalAdId>
        </Creative>
      </Creatives>
    </InLine>
  </Ad>
</VAST>
`
