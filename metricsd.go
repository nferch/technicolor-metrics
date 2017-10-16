package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/headzoo/surf/browser"
	"github.com/influxdata/influxdb/client/v2"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/headzoo/surf.v1"
)

type metricsdConfig struct {
	PollDelay int
	Modem     modemConfig
	InfluxDB  influxDBConfig
}
type modemConfig struct {
	Address  string
	Port     int
	Username string
	Password string
}

type influxDBConfig struct {
	Protocol    string
	Address     string
	Port        int
	Username    string
	Password    string
	Database    string
	Measurement string
	Noop        bool
}

const defaultModemIP string = "192.168.100.1"
const defaultModemPort int = 80
const defaultModemUsername string = "admin"
const defaultInfluxDBPort = 8086
const defaultInfluxMeasurement = "cablemodem"
const defaultPollDelay int = 600
const networkStatsURL string = "vendor_network.asp"

func main() {
	app := cli.NewApp()
	app.Name = "technicolor-metrics"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: filepath.Join("/etc", app.Name, "metrics.conf"),
			Usage: "Use this configuration file instead of ",
		},
		cli.StringFlag{
			Name:  "measurement, m",
			Usage: "Use measurement instead of default or from config file",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Enable debugging output",
		},
		cli.BoolFlag{
			Name:  "noop, n",
			Usage: "Dry run mode, don't send metrics",
		},
	}
	app.Action = run
	app.Run(os.Args)
}

func run(c *cli.Context) error {
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	// TODO: commandline usage of noop should override config
	config := metricsdConfig{
		PollDelay: defaultPollDelay,
		Modem:     modemConfig{Address: defaultModemIP, Port: defaultModemPort, Username: defaultModemUsername},
		InfluxDB:  influxDBConfig{Port: defaultInfluxDBPort, Measurement: defaultInfluxMeasurement, Noop: c.Bool("noop")},
	}

	readConfig(&config, c)
	if c.GlobalString("measurement") != "" {
		config.InfluxDB.Measurement = c.GlobalString("measurement")
		log.Warningf("writing to measurement %s", config.InfluxDB.Measurement)
	}
	ifc, iferr := Connect(&config.InfluxDB)
	if iferr != nil {
		log.Fatal(iferr)
	}

	bow := surf.NewBrowser()
	for {
		SubmitMetrics(&config, bow, ifc)
		time.Sleep(time.Duration(config.PollDelay) * time.Second)
	}
}

func readConfig(config *metricsdConfig, ctx *cli.Context) {
	tomlData, ferr := ioutil.ReadFile(ctx.GlobalString("config"))
	if ferr != nil {
		log.Fatalf("Can't read config file %s: %s", ctx.GlobalString("config"), ferr)
	}

	if _, err := toml.Decode(string(tomlData), config); err != nil {
		log.Fatalf("Can't parse config: %s", err)
	}
}

// SubmitMetrics orchestrates the connection/collection/submission process.
func SubmitMetrics(config *metricsdConfig, bow *browser.Browser, ifc client.Client) {
	modemURL := fmt.Sprintf("http://%s:%d/%s", config.Modem.Address, config.Modem.Port, networkStatsURL)
	log.Debugf("Connecting to %s", modemURL)
	if err := bow.Open(modemURL); err != nil {
		log.Fatalf("Can't connect to modem: %s", err)
	}

	log.Debugf("fetched page %s", bow.Title())
	/*
		for _, e := range bow.Forms() {
			fmt.Printf("%v\n", e)
		}
	*/
	if len(bow.Find("#login").Nodes) > 0 {
		log.Debug("logging in")
		if lerr := login(bow, config); lerr != nil {
			log.Fatalf("Error logging in: %s", lerr)
		}
	}
	// fmt.Println(bow.Body())
	s := WANStatsPage{Body: bow.Dom()}

	dhl := DownstreamResultList{}
	usl := UpstreamResultList{}

	if err := dhl.ParseFromStatsPage(&s); err != nil {
		log.Fatalf("Can't parse Downstream stats: %s", err)
	}
	if err := usl.ParseFromStatsPage(&s); err != nil {
		log.Fatalf("Can't parse Upstream stats: %s", err)
	}

	dhl.EmitToInfluxDB(ifc, &config.InfluxDB)
	usl.EmitToInfluxDB(ifc, &config.InfluxDB)

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
	log.Debug("submitting form")
	if lerr := lform.Submit(); lerr != nil {
		log.Fatalf("Error submitting login form: %v", lerr)
	}
	if bow.StatusCode() > 200 {
		log.Fatalf("Error submitting login form: %v", bow.StatusCode())
	}
	return nil
}
