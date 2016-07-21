package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// http://nominatim.openstreetmap.org/reverse?format=json&lat=47.07038573920726776123046875&lon=8.19601108320057392120361328125&zoom=18&addressdetails=1
// const xml = `<?xml version="1.0" encoding="UTF-8"?>
//<gpx creator="Garmin Connect" version="1.1"
//  xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/11.xsd"
//  xmlns="http://www.topografix.com/GPX/1/1"
//  xmlns:ns3="http://www.garmin.com/xmlschemas/TrackPointExtension/v1"
//  xmlns:ns2="http://www.garmin.com/xmlschemas/GpxExtensions/v3" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
//  <metadata>
//    <link href="connect.garmin.com">
//      <text>Garmin Connect</text>
//    </link>
//    <time>2016-07-10T04:50:18.000Z</time>
//  </metadata>
//  <trk>
//    <name>Glaubenbielen</name>
//    <type>cycling</type>
//    <trkseg>
//      <trkpt lat="47.26096071302890777587890625" lon="7.82374846749007701873779296875">
//        <ele>427.399993896484375</ele>
//        <time>2016-07-10T04:50:18.000Z</time>
//        <extensions>
//          <ns3:TrackPointExtension>
//            <ns3:atemp>19.0</ns3:atemp>
//            <ns3:hr>90</ns3:hr>
//            <ns3:cad>0</ns3:cad>
//          </ns3:TrackPointExtension>
//        </extensions>
//      </trkpt>`

type Query struct {
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

func run(file string) error {

	if !strings.HasSuffix(file, "gpx") {
		return errors.New("filename: not a valid gpx file")
	}

	xmlFile, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer xmlFile.Close()

	var q Query
	decoder := xml.NewDecoder(xmlFile)
	err = decoder.Decode(&q)
	if err != nil {
		return err
	}

	var outputfile string
	if q.Track.Name != "" && len(q.Track.Name) > 3 {
		outputfile = NormalizeText(q.Track.Name) + ".csv"
	} else {
		outputfile = strings.TrimSuffix(file, "gpx") + "csv"
	}
	createFile(outputfile)

	f, err := os.OpenFile(outputfile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, t := range q.Track.TrackSegments {
		totalTrackPoints := len(t.TrackPoints)
		f.WriteString("date,lat,lon,ele,temp,Road,Village,City,Town,City2,Neighbourhood,State,Postcode,Country,DisplayName,\n")
		for i, trackPoint := range t.TrackPoints {
			if i%100 != 0 {
				continue
			}

			log.Printf("%d/%d\t%+v", i, totalTrackPoints, trackPoint)

			u, err := url.Parse("http://nominatim.openstreetmap.org/reverse?format=json&zoom=18&addressdetails=1")
			if err != nil {
				log.Fatal(err)
			}
			q := u.Query()
			q.Set("lat", floatToString(trackPoint.Lat))
			q.Set("lon", floatToString(trackPoint.Lon))
			u.RawQuery = q.Encode()

			var res Result
			err = getJson(u.String(), &res)
			if err != nil {
				log.Fatal(err)
			}

			printCSV(f, trackPoint, res)
		}
	}

	return nil
}

func createFile(p string) {
	f, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()
}

func printCSV(f *os.File, tp TrackPoint, res Result) {
	addr := res.Address
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%v,", tp.Date))
	b.WriteString(fmt.Sprintf("%v,", tp.Lat))
	b.WriteString(fmt.Sprintf("%v,", tp.Lon))
	b.WriteString(fmt.Sprintf("%v,", tp.Ele))
	b.WriteString(fmt.Sprintf("%v,", tp.Temp))
	b.WriteString(fmt.Sprintf("%q,", addr.Road))
	b.WriteString(fmt.Sprintf("%q,", addr.Village))
	b.WriteString(fmt.Sprintf("%q,", addr.City))
	b.WriteString(fmt.Sprintf("%q,", addr.Town))
	if addr.Village != "" {
		b.WriteString(fmt.Sprintf("%q,", addr.Village))
	} else if addr.City != "" {
		b.WriteString(fmt.Sprintf("%q,", addr.City))
	} else if addr.Town != "" {
		b.WriteString(fmt.Sprintf("%q,", addr.Town))
	} else if addr.Neighbourhood != "" {
		b.WriteString(fmt.Sprintf("%q,", addr.Neighbourhood))
	} else {
		b.WriteString(fmt.Sprintf("%q,", addr.State))
	}
	b.WriteString(fmt.Sprintf("%q,", addr.Neighbourhood))
	b.WriteString(fmt.Sprintf("%q,", addr.State))
	b.WriteString(fmt.Sprintf("%q,", addr.Postcode))
	b.WriteString(fmt.Sprintf("%q,", addr.Country))
	b.WriteString(fmt.Sprintf("%q,", res.DisplayName))
	b.WriteString(fmt.Sprintf("\n"))

	f.WriteString(b.String())
}

type Result struct {
	DisplayName string `json:"display_name"`
	Address     Address
}
type Address struct {
	Village       string
	Road          string
	Neighbourhood string
	Town          string
	City          string
	State         string
	Postcode      string
	Country       string
}

func floatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 8, 64)
}
func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func NormalizeText(s string) string {
	isMn := func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsDigit(r))
	}
	tf := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	name, _, _ := transform.String(tf, s)
	return name
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-file.gpx>\n", filepath.Base(os.Args[0]))
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

//:::  This routine calculates the distance between two points (given the     :::
//:::  latitude/longitude of those points). It is being used to calculate     :::
//:::  the distance between two locations using GeoDataSource (TM) prodducts  :::
//:::                                                                         :::
//:::  Definitions:                                                           :::
//:::    South latitudes are negative, east longitudes are positive           :::
//:::                                                                         :::
//:::  Passed to function:                                                    :::
//:::    lat1, lon1 = Latitude and Longitude of point 1 (in decimal degrees)  :::
//:::    lat2, lon2 = Latitude and Longitude of point 2 (in decimal degrees)  :::
//:::    unit = the unit you desire for results                               :::
//:::           where: 'M' is statute miles (default)                         :::
//:::                  'K' is kilometers                                      :::
//:::                  'N' is nautical miles                                  :::
//:::                                                                         :::
//:::  Worldwide cities and other features databases with latitude longitude  :::
//:::  are available at http://www.geodatasource.com
//func distance(lat1, lon1, lat2, lon2 float64, unit string) {
//	var radlat1 = math.Pi * lat1 / 180
//	var radlat2 = math.Pi * lat2 / 180
//	var theta = lon1 - lon2
//	var radtheta = math.Pi * theta / 180
//	var dist = math.Sin(radlat1) * math.Sin(radlat2) + math.Cos(radlat1) * math.Cos(radlat2) * math.Cos(radtheta);
//	dist = math.Acos(dist)
//	dist = dist * 180 / math.Pi
//	dist = dist * 60 * 1.1515
//	if (unit == "K") {
//		dist = dist * 1.609344
//	}
//	if (unit == "N") {
//		dist = dist * 0.8684
//	}
//	return dist
//}
