# Maillog Exporter for Prometheus

Maillog Exporter consumes one or more log files from Postfix, Dovecot and Postgrey and exposes the metrics for Prometheus.

## Running

	Usage of ./maillog_exporter:
	  -listen string
	    	address to listen on (default ":9290")
	  -logpath string
	    	locations of log file that will be grepped (default "/var/log/maillog /var/log/dovecot.log")

## Building

Requires Go ≥ 1.7

	make build
