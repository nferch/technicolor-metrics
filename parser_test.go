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
	for i := range fixtureDsr.DownstreamChannels {
		if drl.DownstreamChannels[i] != fixtureDsr.DownstreamChannels[i] {
			t.Errorf("Downstream() results differ for %d: %v vs %v", i, drl.DownstreamChannels[i], fixtureDsr.DownstreamChannels[i])
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
	for i := 0; i < Max(len(url.UpstreamChannels), len(fixtureUsr.UpstreamChannels)); i++ {
		if len(url.UpstreamChannels)-1 < i {
			t.Errorf("Upstream() missing results for idx %d, should be %#v\n", i, fixtureUsr.UpstreamChannels[i])
		} else if len(fixtureUsr.UpstreamChannels)-1 < i {
			t.Errorf("Upstream() missing results for idx %d, should be %#v\n", i, url.UpstreamChannels[i])
		} else if url.UpstreamChannels[i] != fixtureUsr.UpstreamChannels[i] {
			t.Errorf("Upstream() results differ for %d: %#v vs %#v", i, url.UpstreamChannels[i], fixtureUsr.UpstreamChannels[i])
		}
	}
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
