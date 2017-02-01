package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/thomersch/maillog_exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	listen := flag.String("listen", ":9290", "address to listen on")
	logpath := flag.String("logpath", "/var/log/maillog /var/log/dovecot.log", "locations of log file that will be grepped")
	flag.Parse()

	collector.RegisterMetrics()
	go collector.ConsumeLogs(strings.Split(*logpath, " "))

	ph := prometheus.Handler()
	http.Handle("/metrics", ph)

	log.Printf("Starting server: %s", *listen)
	err := http.ListenAndServe(*listen, nil)
	if err != nil {
		log.Fatal(err)
	}
}
