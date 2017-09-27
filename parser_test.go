package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-cmp/cmp"
)

// Ensure interfaces are fully implemented
var _ ResultList = (*DownstreamResultList)(nil)
var _ ResultList = (*UpstreamResultList)(nil)

func compareResultList(t *testing.T, r ResultList, f ResultList, s *Stats) {
	userr := r.ParseFromStats(s)
	if userr != nil {
		t.Errorf("Error %s ParseFromStats(): %v", reflect.TypeOf(r), userr)
	}
	if !cmp.Equal(r, f) {
		t.Errorf("%s results differ for %v vs %v", reflect.TypeOf(r), r, f)
	}
}

func TestParsers(t *testing.T) {
	fixture, err := os.Open(fixtureFilename)
	if err != nil {
		t.Fatalf("Can't open fixture file: %v", err)
	}

	sel, derr := goquery.NewDocumentFromReader(fixture)
	if derr != nil {
		t.Fatalf("Can't parse fixture file: %v", derr)
	}

	s := &Stats{body: sel}

	dhl := DownstreamResultList{}
	usl := UpstreamResultList{}

	compareResultList(t, &dhl, &fixtureDsr, s)
	compareResultList(t, &usl, &fixtureUsr, s)
}
