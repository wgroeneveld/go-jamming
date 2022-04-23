package notifier

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"encoding/base64"
	"github.com/rs/zerolog/log"
	"net/mail"
	"net/smtp"
)

type MailNotifier struct {
	Conf *common.Config
}

// sendMail is a utility function that sends a mail without authentication to localhost. Tested using postfix.
// cheers https://github.com/gadelkareem/go-helpers/blob/master/helpers.go
func sendMail(from, subject, body, toName, toAddress string) error {
	c, err := smtp.Dial("127.0.0.1:25")
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Mail(from); err != nil {
		return err
	}

	to := (&mail.Address{toName, toAddress}).String()
	if err = c.Rcpt(to); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	msg := "To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func (mn *MailNotifier) NotifyReceived(wm mf.Mention, indieweb *mf.IndiewebData) {
	err := sendMail(
		"admin@brainbaking.com",
		"Webmention in moderation from "+wm.SourceDomain(),
		BuildNotification(wm, indieweb, mn.Conf),
		"Go-Jamming User",
		"wouter@brainbaking.com")

	if err != nil {
		log.Err(err).Msg("Unable to send notification mail, check localhost postfix settings?")
	}
}
