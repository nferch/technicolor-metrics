package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/headzoo/surf/browser"
	"gopkg.in/headzoo/surf.v1"
)

type metricsdConfig struct {
	Modem modemConfig
}
type modemConfig struct {
	IP       string
	Port     int
	Username string
	Password string
}

const defaultModemIP string = "192.168.100.1"
const defaultModemPort int = 80
const defaultModemUsername string = "admin"
const networkStatsURL string = "vendor_network.asp"

func main() {
	config := metricsdConfig{Modem: modemConfig{IP: defaultModemIP, Port: defaultModemPort, Username: defaultModemUsername}}

	tomlData, ferr := ioutil.ReadFile("metricsd.conf")
	if ferr != nil {
		log.Fatalf("Can't read config file: %s", ferr)
	}

	if _, err := toml.Decode(string(tomlData), &config); err != nil {
		log.Fatalf("Can't parse config: %s", err)
	}

	bow := surf.NewBrowser()
	modemURL := fmt.Sprintf("http://%s:%d/%s", config.Modem.IP, config.Modem.Port, networkStatsURL)
	log.Printf("Connecting to %s", modemURL)
	if err := bow.Open(modemURL); err != nil {
		log.Fatalf("Can't connect to modem: %s", err)
	}

	fmt.Println(bow.Title())
	/*
		for _, e := range bow.Forms() {
			fmt.Printf("%v\n", e)
		}
	*/
	if len(bow.Find("#login").Nodes) > 0 {
		log.Println("logging in")
		if lerr := login(bow, config); lerr != nil {
			log.Fatalf("Error logging in: %s", lerr)
		}
	}
	// fmt.Println(bow.Body())
	s := Stats{body: bow.Dom()}

	dhl := DownstreamResultList{}
	usl := UpstreamResultList{}

	if err := dhl.ParseFromStats(&s); err != nil {
		log.Fatalf("Can't parse Downstream stats: %s", err)
	}
	if err := usl.ParseFromStats(&s); err != nil {
		log.Fatalf("Can't parse Upstream stats: %s", err)
	}

	// #content > div:nth-child(6) > table > tbody > tr:nth-child(4) > th
}

func login(bow *browser.Browser, config *metricsdConfig) error {
	lform, err := bow.Form("#pageForm")

	if err != nil {
		log.Fatalf("Cannot find login form: %v", err)
	}

	if err = lform.Input("loginUsername", config.Modem.Username); err != nil {
		log.Fatalf("Can't find username field in HTML: %s", err)
	}
	if err = lform.Input("loginPassword", config.Modem.Password); err != nil {
		log.Fatalf("Can't find password field in HTML: %s", err)
	}
	log.Println("submitting form")
	if lerr := lform.Submit(); lerr != nil {
		log.Fatalf("Error submitting login form: %v", lerr)
	}
	if bow.StatusCode() > 200 {
		log.Fatalf("Error submitting login form: %v", bow.StatusCode())
	}
	return nil
}
