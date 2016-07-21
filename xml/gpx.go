package xml

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

type Gpx struct {
	XMLName xml.Name `xml:"gpx"`
	Track   Track    `xml:"trk"`
}

type Track struct {
	Name          string         `xml:"name"`
	Type          string         `xml:"type"`
	TrackSegments []TrackSegment `xml:"trkseg"`
}

type TrackSegment struct {
	TrackPoints []TrackPoint `xml:"trkpt"`
}

type TrackPoint struct {
	Lat  float64   `xml:"lat,attr"`
	Lon  float64   `xml:"lon,attr"`
	Date time.Time `xml:"time"`
	Ele  float64   `xml:"ele"` // Höhe in m
	Temp string    `xml:"extensions>TrackPointExtension>atemp"`
}

func ParseGpx(filename string) (Gpx, error) {
	xmlFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return Gpx{}, err
	}
	defer xmlFile.Close()

	var q Gpx
	decoder := xml.NewDecoder(xmlFile)
	err = decoder.Decode(&q)
	return q, err
}
