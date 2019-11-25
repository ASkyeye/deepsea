package rmailer

import (
	"crypto/tls"
	"deepsea/global"
	"errors"
	gomail "github.com/gophish/gomail"
	thtml "html/template"
	ttext "html/template"
	"io"
	"log"
	"path/filepath"
	"strings"
)

func GenMail(
	server string,
	from string,
	subject string,
	bodyTextTemplate string,
	bodyHtmlTemplate string,
	attachments []string,
	embeds []string,
	headers map[string]string,
	tdata *TemplateData) (*gomail.Message, error)  {

	var err error

	log.Printf("[Debug] Identifier: %s | Email: %s | First Name: %s | Last Name: %s\n",
		tdata.Mark.Identifier, tdata.Mark.Email,
		tdata.Mark.Firstname, tdata.Mark.Lastname)

	m := gomail.NewMessage()
	m.SetHeader("Subject", subject)
	m.SetHeader("From", from)
	m.SetHeader("To", tdata.Mark.Email)

	// Set Headers
	for key, value := range headers {
		m.SetHeader(key, value)
	}

	// Create a Message-Id:
	msgId := strings.Join([]string{global.RandString(16), server}, "@")
	m.SetHeader("Message-ID", "<"+msgId+">")

	// Templates HTML/Text
	th, err := thtml.ParseFiles(bodyHtmlTemplate)
	if err != nil {
		return new(gomail.Message), err
	}
	tt, err := ttext.ParseFiles(bodyTextTemplate)
	if err != nil {
		return new(gomail.Message), err
	}

	// Embedded images
	l := len(embeds)
	if l != 0 {
		tdata.EmbedImage = make([]string, l)
		for ix, file := range embeds {
			log.Println("[Info] Embedding: ", file)

			if global.FileExists(file) {
				m.Embed(file)
				tdata.EmbedImage[ix] = filepath.Base(file)
			} else {
				return new(gomail.Message),
				errors.New("Embedded File Not Found: " + file)
			}
		}
	}

	// URLs
	// Compile and write Templates
	m.AddAlternativeWriter("text/plain", func(w io.Writer) error {
		return tt.Execute(w, tdata)
	})
	m.AddAlternativeWriter("text/html", func(w io.Writer) error {
		return th.Execute(w, tdata)
	})

	// Attachments
	for _, file := range attachments {
		log.Println("[Info] Attaching asset : ", file)
		if global.FileExists(file) {
			m.Attach(file)
		} else {
			return new(gomail.Message),
			errors.New("Attachment File Not Found: " + file)
		}
	}
	return m, nil
}

func DialSend(
	m *gomail.Message,
	server string, port int, username string, password string, usetls string) {

	d := gomail.NewDialer(server, port, username, password)
	if strings.ToLower(usetls) == "yes" {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("[Error] Could not dial and send: %v", err)
	}
}
