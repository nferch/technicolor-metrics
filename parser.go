package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

// Stats is a container for data scraped from the vendor_network.asp page.
type Stats struct {
	body *goquery.Document
}

// ResultList is a common interface for parsed data
type ResultList interface {
	ParseFromStats(*Stats) error
	parseFromSelection(*goquery.Selection) error
}

// DownstreamResultList is the stats for downstream.
type DownstreamResultList struct {
	Channels []DownstreamResult
}

// DownstreamResult is a single channel's stats
type DownstreamResult struct {
	Index      uint8
	LockStatus string
	Frequency  int16
	SNR        float32
	Power      float32
	Modulation string
}

// UpstreamResultList is the stats for downstream.
type UpstreamResultList struct {
	Channels []UpstreamResult
}

// UpstreamResult is a single channel's stats
type UpstreamResult struct {
	Index      uint8
	LockStatus string
	Frequency  int16
	SymbolRate int16
	Power      float32
	Modulation string
}

func fromStats(r ResultList, s *Stats, direction string) error {
	tab, err := findStatsTable(s.body, direction)
	if err != nil {
		return err
	}
	return r.parseFromSelection(tab)
}

// ParseFromStats extracts the stats from the HTML in the Stats page
func (drl *DownstreamResultList) ParseFromStats(s *Stats) error {
	return fromStats(drl, s, "Downstream")
}

// ParseFromStats extracts the stats from the HTML in the Stats page
func (url *UpstreamResultList) ParseFromStats(s *Stats) error {
	return fromStats(url, s, "Upstream")
}

func findStatsTable(sel *goquery.Document, direction string) (*goquery.Selection, error) {
	var tab *goquery.Selection
	sel.Find(".module > table").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if s.Find("thead > tr > th").First().Text() == direction {
			tab = s
			return false
		}
		return true
	})
	// fmt.Printf("sel is %v\n", tab.Text())
	if tab != nil {
		return tab, nil
	}
	return nil, errors.New("cannot find table in HTML")
}

func (drl *DownstreamResultList) parseFromSelection(sel *goquery.Selection) error {
	rowTitle := ""
	sel.Find("tbody > tr").Each(func(i int, tr *goquery.Selection) {
		// fmt.Printf("tr %#v\n", tr)
		tr.Find("th").Each(func(j int, th *goquery.Selection) {
			// fmt.Printf("th %#v\n", th.Text())
			rowTitle = th.Text()
		})
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			// fmt.Printf("td %d %#v\n", j, td.Text())
			if len(drl.Channels) < j+1 {
				drl.Channels = append(drl.Channels, DownstreamResult{})
			}
			tdint, _ := strconv.Atoi(strings.Split(strings.TrimSpace(td.Text()), " ")[0])
			tdfloat, _ := strconv.ParseFloat(strings.Split(strings.TrimSpace(td.Text()), " ")[0], 32)
			switch rowTitle {
			case "Index":
				drl.Channels[j].Index = uint8(tdint)
			case "Lock Status":
				drl.Channels[j].LockStatus = td.Text()
			case "Frequency":
				drl.Channels[j].Frequency = int16(tdint)
			case "SNR":
				drl.Channels[j].SNR = float32(tdfloat)
			case "Power":
				drl.Channels[j].Power = float32(tdfloat)
			case "Modulation":
				drl.Channels[j].Modulation = td.Text()
			}
		})
	})
	fmt.Printf("Parsed as %#v\n", drl)
	return nil
}

func (url *UpstreamResultList) parseFromSelection(sel *goquery.Selection) error {
	rowTitle := ""
	sel.Find("tbody > tr").Each(func(i int, tr *goquery.Selection) {
		// fmt.Printf("tr %#v\n", tr)
		tr.Find("th").Each(func(j int, th *goquery.Selection) {
			// fmt.Printf("th %#v\n", th.Text())
			rowTitle = th.Text()
		})
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			// fmt.Printf("td %d %#v\n", j, td.Text())
			if len(url.Channels) < j+1 {
				url.Channels = append(url.Channels, UpstreamResult{})
			}
			tdint, _ := strconv.Atoi(strings.Split(strings.TrimSpace(td.Text()), " ")[0])
			tdfloat, _ := strconv.ParseFloat(strings.Split(strings.TrimSpace(td.Text()), " ")[0], 32)
			switch rowTitle {
			case "Index":
				url.Channels[j].Index = uint8(tdint)
			case "Lock Status":
				url.Channels[j].LockStatus = td.Text()
			case "Frequency":
				url.Channels[j].Frequency = int16(tdint)
			case "Symbol Rate":
				url.Channels[j].SymbolRate = int16(tdint)
			case "Power":
				url.Channels[j].Power = float32(tdfloat)
			case "Modulation":
				url.Channels[j].Modulation = td.Text()
			}
		})
	})
	fmt.Printf("Parsed as %#v\n", url)
	return nil
}
