package collector

import "github.com/prometheus/client_golang/prometheus"

const namespace = "mail"

var (
	openConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mail_open_connections",
		Help: "Number of open connections to the server.",
	}, []string{"protocol"})

	mailsTransferred = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "mail_mails_transferred",
		Help: "Transferred mails in any direction",
	}, []string{"direction"})

	greylistedMail = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "mail_greylisted",
		Help: "Count of greylisted or passed mails.",
	}, []string{"status"})

	blockedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "mail_incoming_blocked",
		Help: "Blocked incoming mails.",
	})

	bouncedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "mail_outgoing_bounced",
		Help: "Outgoing mails that bounced.",
	})
)

func RegisterMetrics() {
	prometheus.MustRegister(openConnections)
	prometheus.MustRegister(mailsTransferred)
	prometheus.MustRegister(greylistedMail)
	prometheus.MustRegister(blockedCounter)
	prometheus.MustRegister(bouncedCounter)
}
