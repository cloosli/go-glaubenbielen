package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/cloosli/go-glaubenbielen/util"
	gpx "github.com/cloosli/go-glaubenbielen/xml"
)

var (
	flagFilename string
	flagOutput   string
	flagDebug    bool
	flagLocal    bool
	flagSteps    int
)

func run() error {

	if !strings.HasSuffix(flagFilename, "gpx") {
		return errors.New("filename: not a valid gpx file")
	}

	q, err := gpx.ParseGpx(flagFilename)
	if err != nil {
		return err
	}

	if flagOutput == "" {
		if q.Track.Name != "" && len(q.Track.Name) > 3 {
			flagOutput = util.NormalizeText(q.Track.Name) + ".csv"
		} else {
			flagOutput = strings.TrimSuffix(flagFilename, "gpx") + "csv"
		}
	}
	util.CreateFile(flagOutput)

	f, err := os.OpenFile(flagOutput, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	//set := make(map[string] struct{})
	set := make(map[string]bool)

	for _, t := range q.Track.TrackSegments {
		totalTrackPoints := len(t.TrackPoints)
		f.WriteString("date,lat,lon,ele,temp,Road,Village,City,Town,City2,Neighbourhood,State,Postcode,Country,DisplayName,\n")
		for i, trackPoint := range t.TrackPoints {
			if i%flagSteps != 0 && i+1 < totalTrackPoints {
				continue
			}

			log.Printf("%d/%d\t%+v", i+1, totalTrackPoints, trackPoint)

			u, err := url.Parse("http://nominatim.openstreetmap.org/reverse?format=json&zoom=18&addressdetails=1")
			if err != nil {
				log.Fatal(err)
			}
			q := u.Query()
			q.Set("lat", util.FloatToString(trackPoint.Lat))
			q.Set("lon", util.FloatToString(trackPoint.Lon))
			u.RawQuery = q.Encode()

			var res Result
			err = getJson(u.String(), &res)
			if err != nil {
				log.Fatal(err)
			}

			city := res.Address.GetBestCity()
			set[city] = true

			printCSV(f, trackPoint, res)
		}
	}

	if len(set) > 0 {
		fmt.Println("\nAll cities:")
		for key := range set {
			fmt.Print(key)
			fmt.Print(", ")
		}
		fmt.Print("\n\n")
	}

	fmt.Printf("CSV file created > %s\n", f.Name())

	return nil
}

func printCSV(f *os.File, tp gpx.TrackPoint, res Result) {
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
	b.WriteString(fmt.Sprintf("%q,", addr.GetBestCity()))
	b.WriteString(fmt.Sprintf("%q,", addr.Neighbourhood))
	b.WriteString(fmt.Sprintf("%q,", addr.State))
	b.WriteString(fmt.Sprintf("%q,", addr.Postcode))
	b.WriteString(fmt.Sprintf("%q,", addr.Country))
	b.WriteString(fmt.Sprintf("%q,", res.DisplayName))
	b.WriteString(fmt.Sprintf("\n"))

	f.WriteString(b.String())
}

func (addr *Address) GetBestCity() string {
	if addr.Village != "" {
		return addr.Village
	} else if addr.City != "" {
		return addr.City
	} else if addr.Town != "" {
		return addr.Town
	} else if addr.Neighbourhood != "" {
		return addr.Neighbourhood
	} else {
		return addr.State
	}
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

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func main() {
	flag.BoolVar(&flagDebug, "v", false, "verbose logging")
	flag.BoolVar(&flagLocal, "local", false, "Use local city list file")
	flag.StringVar(&flagFilename, "i", "", "Input file: -i <path-to-file.gpx>")
	flag.StringVar(&flagOutput, "o", "", "Output file: -o <path-to-file.csv>")
	flag.IntVar(&flagSteps, "s", 100, "Check every x steps the location")
	showUsage := flag.Bool("h", false, "Show usage")
	flag.Parse()

	if *showUsage {
		flag.Usage()
		os.Exit(0)
	}

	if flagFilename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(); err != nil {
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
