package main

import (
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestDownstream(t *testing.T) {
	fixture, err := os.Open(fixtureFilename)
	if err != nil {
		t.Fatalf("Can't open fixture file: %v", err)
	}

	sel, derr := goquery.NewDocumentFromReader(fixture)
	if derr != nil {
		t.Fatalf("Can't parse fixture file: %v", derr)
	}

	s := &Stats{body: sel}
	drl, dserr := s.Downstream()
	if dserr != nil {
		t.Errorf("Error Downstream(): %v", dserr)
	}
	for i := range fixtureDsr.Channels {
		if drl.Channels[i] != fixtureDsr.DownstreamChannels[i] {
			t.Errorf("Downstream() results differ for %d: %v vs %v", i, drl.Channels[i], fixtureDsr.DownstreamChannels[i])
		}
	}
}

func TestUpstream(t *testing.T) {
	fixture, err := os.Open(fixtureFilename)
	if err != nil {
		t.Fatalf("Can't open fixture file: %v", err)
	}

	sel, derr := goquery.NewDocumentFromReader(fixture)
	if derr != nil {
		t.Fatalf("Can't parse fixture file: %v", derr)
	}

	s := &Stats{body: sel}
	url, userr := s.Upstream()
	if userr != nil {
		t.Errorf("Error Upstream(): %v", userr)
	}
	for i := 0; i < Max(len(url.Channels), len(fixtureUsr.Channels)); i++ {
		if len(url.Channels)-1 < i {
			t.Errorf("Upstream() missing results for idx %d, should be %#v\n", i, fixtureUsr.Channels[i])
		} else if len(fixtureUsr.Channels)-1 < i {
			t.Errorf("Upstream() missing results for idx %d, should be %#v\n", i, url.Channels[i])
		} else if url.Channels[i] != fixtureUsr.Channels[i] {
			t.Errorf("Upstream() results differ for %d: %#v vs %#v", i, url.Channels[i], fixtureUsr.Channels[i])
		}
	}
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
