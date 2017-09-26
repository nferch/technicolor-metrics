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
