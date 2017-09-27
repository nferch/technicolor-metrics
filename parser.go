package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

/* Container for Stats page */
type Stats struct {
	body *goquery.Document
}

/* Downstream stats */
type DownstreamResultList struct {
	DownstreamChannels []DownstreamResult
}

/* A single channel's stats */
type DownstreamResult struct {
	Index      uint8
	LockStatus string
	Frequency  int16
	SNR        float32
	Power      float32
	Modulation string
}

func (s *Stats) Downstream() (*DownstreamResultList, error) {
	// #content > div:nth-child(6)
	// s.body.Find("//*[@id=\"module\"]/table/thead/tr/th[1]").Each(func(i int, s *goquery.Selection) {
	//*[@id="content"]/div[5]/table/thead/tr/th[1]
	tab, err := findStatsTable(s.body, "Downstream")
	if err != nil {
		return nil, err
	}
	return parseDownstream(tab)
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
			if len(drl.DownstreamChannels) < j+1 {
				drl.DownstreamChannels = append(drl.DownstreamChannels, DownstreamResult{})
			}
			tdint, _ := strconv.Atoi(strings.Split(strings.TrimSpace(td.Text()), " ")[0])
			tdfloat, _ := strconv.ParseFloat(strings.Split(strings.TrimSpace(td.Text()), " ")[0], 32)
			switch rowTitle {
			case "Index":
				drl.DownstreamChannels[j].Index = uint8(tdint)
			case "Lock Status":
				drl.DownstreamChannels[j].LockStatus = td.Text()
			case "Frequency":
				drl.DownstreamChannels[j].Frequency = int16(tdint)
			case "SNR":
				drl.DownstreamChannels[j].SNR = float32(tdfloat)
			case "Power":
				drl.DownstreamChannels[j].Power = float32(tdfloat)
			case "Modulation":
				drl.DownstreamChannels[j].Modulation = td.Text()
			}
		})
	})
	return &drl, nil
}
