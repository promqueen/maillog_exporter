package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	t.Run("connect_smtp", func(t *testing.T) {
		evt := parseLine("Feb  1 15:31:29 mail postfix/submission/smtpd[44640]: connect from p200300DCFBCD00000B4C106896D280000.dip0.t-ipconnect.de[2003:dc:fbcd:0000:b4c1:689:6d28:0000]")
		assert.Equal(t, connectSMTP, evt)
	})

	t.Run("outgoing_mail", func(t *testing.T) {
		evt := parseLine("Feb  1 15:31:31 mail postfix/smtp[44644]: 741F51009D: to=<asd@asdasd>, relay=asdasd[217.0.0.0]:25, delay=1.9, delays=0.09/0/0.38/1.4, dsn=2.0.0, status=sent (250 2.0.0 Ok: queued as AAAAAAAA)")
		assert.Equal(t, outgoingMail, evt)
	})

	t.Run("incoming_mail", func(t *testing.T) {
		evt := parseLine("Feb  1 07:54:25 mail postfix/lmtp[27157]: 05DAA10087: to=<test@example.com>, relay=some.example.com[private/dovecot-lmtp], delay=1.4, delays=0.9/0/0.01/0.48, dsn=2.0.0, status=sent (250 2.0.0 <some@example.com> AAAAAA Saved)")
		assert.Equal(t, incomingMail, evt)
	})

	t.Run("greylisted", func(t *testing.T) {
		evt := parseLine("Feb  1 10:23:46 mail postgrey[1147]: action=greylist, reason=new, client_name=arqqeepi.my-addr.com, client_address=217.0.0.0, sender=koyot266190@arqqeepi.my-addr.com, recipient=test@example.com")
		assert.Equal(t, greylisted, evt)
	})

	t.Run("greylist_pass", func(t *testing.T) {
		evt := parseLine("Feb  1 07:54:25 mail postgrey[1147]: action=pass, reason=triplet found, delay=499, client_name=test.example.com, client_address=test.example.com, sender=a@test.example.com, recipient=receiver@example.com")
		assert.Equal(t, greylistPass, evt)
	})

	t.Run("user unknown", func(t *testing.T) {
		evt := parseLine("Feb  1 15:34:43 mail postfix/smtpd[44718]: NOQUEUE: reject: RCPT from nm27.bullet.mail.ne1.yahoo.com[98.138.90.90]: 550 5.1.1 <test@example.com>: Recipient address rejected: User unknown in virtual mailbox table; from=<frank.richard112@yahoo.com> to=<test@example.com> proto=ESMTP helo=<nm27.bullet.mail.ne1.yahoo.com>")
		assert.Equal(t, unknownUser, evt)
	})

	t.Run("imap login", func(t *testing.T) {
		evt := parseLine("Feb 01 16:34:40 imap-login: Info: Login: user=<thomas@example.com>, method=PLAIN, rip=217.0.0.0, lip=217.0.0.1, mpid=46118, TLS, session=<P/AAAAA>")
		assert.Equal(t, connectIMAP, evt)
	})

	t.Run("imap logout", func(t *testing.T) {
		evt := parseLine("Feb 01 16:34:40 imap(thomas@example.com): Info: Logged out in=34 out=471")
		assert.Equal(t, disconnectIMAP, evt)
	})

	t.Run("dovecot starting", func(t *testing.T) {
		evt := parseLine("Feb 16 07:33:20 master: Info: Dovecot v2.2.27 (c0f36b0) starting up for imap, lmtp")
		assert.Equal(t, resetIMAP, evt)
	})
}
