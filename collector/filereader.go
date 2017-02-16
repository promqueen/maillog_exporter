package collector

import (
	"log"
	"strings"

	"github.com/hpcloud/tail"
)

//go:generate stringer -type=event

type event int

const (
	eventEmpty event = iota
	connectSMTP
	disconnectSMTP
	connectIMAP
	disconnectIMAP
	resetIMAP
	incomingMail
	outgoingMail
	unknownUser
	greylisted
	greylistPass
	spamBlocked
	bounced
)

func ConsumeLogs(logs []string) {
	for _, logFile := range logs {
		go func(logFile string) {
			t, err := tail.TailFile(logFile, tail.Config{
				ReOpen: true,
				Follow: true,
			})
			if err != nil {
				log.Fatalf("could not consume log %v, because: %v", logFile, err)
			}

			for line := range t.Lines {
				switch parseLine(line.Text) {
				case eventEmpty:
					// nothing
				case connectSMTP:
					openConnections.WithLabelValues("smtp").Inc()
				case disconnectSMTP:
					openConnections.WithLabelValues("smtp").Dec()
				case connectIMAP:
					openConnections.WithLabelValues("imap").Inc()
				case disconnectIMAP:
					openConnections.WithLabelValues("imap").Dec()
				case resetIMAP:
					openConnections.WithLabelValues("imap").Set(0)
				case incomingMail:
					mailsTransferred.WithLabelValues("in").Add(1)
				case outgoingMail:
					mailsTransferred.WithLabelValues("out").Add(1)
				case greylisted:
					greylistedMail.WithLabelValues("blocked").Add(1)
				case greylistPass:
					greylistedMail.WithLabelValues("passed").Add(1)
				case spamBlocked:
					blockedCounter.Add(1)
				case bounced:
					bouncedCounter.Add(1)
				default:
					log.Printf("unhandled event: %s", parseLine(line.Text))
				}
			}
		}(logFile)
	}
}

func parseLine(line string) event {
	parts := strings.Split(line, " ")
	if len(parts) < 6 {
		return eventEmpty
	}

	logger := parts[5]
	msg := strings.Join(parts[6:], " ")
	if strings.HasPrefix(parts[3], "imap") || strings.HasPrefix(parts[3], "master") {
		logger = parts[3]
		msg = strings.Join(parts[4:], " ")
	}

	switch {
	case strings.Contains(logger, "postfix/submission"):
		if strings.HasPrefix(msg, "connect from") {
			return connectSMTP
		}
		if strings.HasPrefix(msg, "disconnect from") {
			return disconnectSMTP
		}

	case strings.HasPrefix(logger, "postfix/lmtp"):
		if strings.Contains(msg, "status=sent") {
			return incomingMail
		}

	case strings.HasPrefix(logger, "postfix/smtpd"):
		if strings.HasPrefix(msg, "disconnect from") {
			return disconnectSMTP
		}
		if strings.HasPrefix(msg, "connect from") {
			return connectSMTP
		}
		if strings.Contains(msg, "User unknown in virtual mailbox table") {
			return unknownUser
		}
		if strings.Contains(msg, "blocked using") {
			return spamBlocked
		}

	case strings.HasPrefix(logger, "postfix/smtp"):
		if strings.HasPrefix(msg, "connect from") {
			return connectSMTP
		}
		if strings.Contains(msg, "status=sent") {
			return outgoingMail
		}
		if strings.Contains(msg, "status=bounced") {
			return bounced
		}

	case strings.HasPrefix(logger, "imap-login"):
		if strings.HasPrefix(msg, "Info: Login:") {
			return connectIMAP
		}

	case strings.HasPrefix(logger, "imap"):
		if strings.HasPrefix(msg, "Info: Logged out") {
			return disconnectIMAP
		}

	case strings.HasPrefix(logger, "postgrey"):
		if strings.Contains(msg, "action=greylist") {
			return greylisted
		}
		if strings.Contains(msg, "action=pass") {
			return greylistPass
		}

	case strings.HasPrefix(logger, "master"):
		if strings.Contains(msg, "starting up") {
			return resetIMAP
		}
	}
	return eventEmpty
}
