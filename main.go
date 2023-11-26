package main

import (
	"flag"
	"fmt"
	"github.com/eugecm/gometar/metar/parser"
	"github.com/eugecm/gometar/sky"
	"github.com/eugecm/gometar/weather"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type METARReport struct {
	AirportCode          string
	ObservationTime      time.Time
	WindDirection        int
	WindSpeed            int
	Visibility           int
	AtmosphericPhenomena []string
	Clouds               []string
	CloudBase            []int
	Temperature          int
	DewPoint             int
	Pressure             int
	SignificantChanges   bool
}

func CreateMETARReport(
	airportCode string,
	observationTime time.Time,
	windDirection int,
	windSpeed int,
	visibility int,
	atmosphericPhenomena []string,
	clouds []string,
	cloudBase []int,
	temperature int,
	dewPoint int,
	pressure int,
	significantChanges bool,
) *METARReport {
	return &METARReport{
		AirportCode:          airportCode,
		ObservationTime:      observationTime,
		WindDirection:        windDirection,
		WindSpeed:            windSpeed,
		Visibility:           visibility,
		AtmosphericPhenomena: atmosphericPhenomena,
		Clouds:               clouds,
		CloudBase:            cloudBase,
		Temperature:          temperature,
		DewPoint:             dewPoint,
		Pressure:             pressure,
		SignificantChanges:   significantChanges,
	}
}

func mapSlice[T any, M any](a []T, f func(T) M) []M {
	n := make([]M, len(a))
	for i, e := range a {
		n[i] = f(e)
	}
	return n
}

func atoiOrZero(strVal string) int {
	result, err := strconv.Atoi(strVal)
	if err != nil {
		return 0
	}

	return result
}

func decodeMETAR(rawReport string) (*METARReport, error) {
	parserObj := parser.New()
	rawReport = strings.Replace(rawReport, "METAR ", "", -1)
	parsedReport, err := parserObj.Parse(rawReport)

	distanceInt, _ := strconv.Atoi(parsedReport.Visibility.Distance)
	qnhPressureInt, _ := strconv.Atoi(parsedReport.Qnh.Pressure)

	return CreateMETARReport(
		parsedReport.Station,
		parsedReport.DateTime,
		parsedReport.Wind.Source,
		parsedReport.Wind.Speed.Speed,
		distanceInt,
		mapSlice(parsedReport.Weather.Phenomena, func(p weather.Phenomenon) string {
			return string(p)
		}),
		mapSlice(parsedReport.Clouds, func(c sky.CloudInformation) string {
			return string(c.Amount)
		}),
		mapSlice(parsedReport.Clouds, func(c sky.CloudInformation) int {
			return atoiOrZero(c.Height) * 100
		}),
		parsedReport.Temperature.Temperature,
		parsedReport.Temperature.DewPoint,
		qnhPressureInt,
		parsedReport.Remarks == "",
	), err
}

func (m *METARReport) compare(o *METARReport) bool {
	return m.AirportCode == o.AirportCode &&
		m.ObservationTime == o.ObservationTime &&
		m.WindDirection == o.WindDirection &&
		m.WindSpeed == o.WindSpeed &&
		m.Visibility == o.Visibility &&
		reflect.DeepEqual(m.AtmosphericPhenomena, o.AtmosphericPhenomena) &&
		reflect.DeepEqual(m.Clouds, o.Clouds) &&
		reflect.DeepEqual(m.CloudBase, o.CloudBase) &&
		m.Temperature == o.Temperature &&
		m.DewPoint == o.DewPoint &&
		m.Pressure == o.Pressure &&
		m.SignificantChanges == o.SignificantChanges
}

func main() {
	rawMetar := flag.String("metar", "", "RAW METAR report")
	flag.Parse()

	decodedMetar, err := decodeMETAR(*rawMetar)

	if err != nil {
		panic(err)
	}

	fmt.Println(*decodedMetar)
}
