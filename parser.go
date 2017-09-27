package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

// Stats is a container for data scraped from the vendor_network.asp page.
type Stats struct {
	body *goquery.Document
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

// Downstream parses downstream stats from HTML
func (s *Stats) Downstream() (*DownstreamResultList, error) {
	// #content > div:nth-child(6)
	//*[@id="content"]/div[5]/table/thead/tr/th[1]
	tab, err := findStatsTable(s.body, "Downstream")
	if err != nil {
		return nil, err
	}
	return parseDownstream(tab)
}

// Upstream parses upstream stats from HTML
func (s *Stats) Upstream() (*UpstreamResultList, error) {
	// #content > div:nth-child(6)
	//*[@id="content"]/div[5]/table/thead/tr/th[1]
	tab, err := findStatsTable(s.body, "Upstream")
	if err != nil {
		return nil, err
	}
	return parseUpstream(tab)
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

func parseDownstream(sel *goquery.Selection) (*DownstreamResultList, error) {
	drl := DownstreamResultList{}
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
	return &drl, nil
}

func parseUpstream(sel *goquery.Selection) (*UpstreamResultList, error) {
	url := UpstreamResultList{}
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
	return &url, nil
}
