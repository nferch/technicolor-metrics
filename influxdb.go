package main

import (
	"fmt"
	"log"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

func Connect(ifconf *influxDBConfig) (client.Client, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     fmt.Sprintf("%s://%s:%d", ifconf.Protocol, ifconf.Address, ifconf.Port),
		Username: ifconf.Username,
		Password: ifconf.Password,
	})
	if err != nil {
		log.Fatalf("Cannot instantiate InfluxDB Client: %s", err)
		return nil, err
	}
	return c, err
}

func (drl *DownstreamResultList) EmitToInfluxDB(clt client.Client, ifconf *influxDBConfig) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  ifconf.Database,
		Precision: "s",
	})
	if err != nil {
		log.Printf("Can't create batch points: %s", err)
	}

	for _, d := range drl.Channels {
		tags := map[string]string{"channel": fmt.Sprintf("%d", (d.Index))}
		fields := map[string]interface{}{
			"Frequency": d.Frequency,
			"SNR":       d.SNR,
			"Power":     d.Power,
		}
		pt, err := client.NewPoint("downstream", tags, fields, time.Now())
		if err != nil {
			log.Printf("Can't create point: %s", err)
		}
		bp.AddPoint(pt)
	}
	if err := clt.Write(bp); err != nil {
		log.Printf("Error writing points: %s", err)
	}
}

func (drl *UpstreamResultList) EmitToInfluxDB(clt client.Client, ifconf *influxDBConfig) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  ifconf.Database,
		Precision: "s",
	})
	if err != nil {
		log.Printf("Can't create batch points: %s", err)
	}

	for _, d := range drl.Channels {
		tags := map[string]string{"channel": fmt.Sprintf("%d", (d.Index))}
		fields := map[string]interface{}{
			"Frequency":  d.Frequency,
			"SymbolRate": d.SymbolRate,
			"Power":      d.Power,
		}
		pt, err := client.NewPoint("upstream", tags, fields, time.Now())
		if err != nil {
			log.Printf("Can't create point: %s", err)
		}
		bp.AddPoint(pt)
	}
	if err := clt.Write(bp); err != nil {
		log.Printf("Error writing points: %s", err)
	}
}
