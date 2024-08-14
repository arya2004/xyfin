package mail

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"path/filepath"
	"strings"
)

const (
	smtpAuthAddress   = "smtp.office365.com"
	smtpServerAddress = "smtp.office365.com:587"
)

type EmailSender interface {
	SendEmail(subject string, templatePath string, to []string, cc []string, bcc []string, attachFiles []string, data map[string]string) error
}

type GmailSender struct {
	name      string
	fromEmail string
	password  string
}

func NewGmailSender(name string, fromEmail string, password string) EmailSender {
	return &GmailSender{name, fromEmail, password}
}

// Function to read and apply a template with data
func loadTemplate(templatePath string, data map[string]string) (string, error) {
	tmplContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("email").Parse(string(tmplContent))
	if err != nil {
		return "", err
	}

	var content bytes.Buffer
	err = tmpl.Execute(&content, data)
	if err != nil {
		return "", err
	}

	return content.String(), nil
}

func (sender *GmailSender) SendEmail(subject string, templatePath string, to []string, cc []string, bcc []string, attachFiles []string, data map[string]string) error {
	// Load the email content from the template with data
	content, err := loadTemplate(templatePath, data)
	if err != nil {
		return err
	}

	// Build the email body
	var msg bytes.Buffer
	writer := multipart.NewWriter(&msg)

	// Set up the main headers
	msg.WriteString(fmt.Sprintf("From: %s\r\n", sender.fromEmail))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ", ")))
	if len(cc) > 0 {
		msg.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(cc, ", ")))
	}
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", writer.Boundary()))

	// Add the HTML content as part of the multipart message
	htmlPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": {"text/html; charset=UTF-8"},
	})
	if err != nil {
		return err
	}
	htmlPart.Write([]byte(content))

	// Attach files
	for _, file := range attachFiles {
		if err := attachFile(writer, file); err != nil {
			return err
		}
	}

	writer.Close()

	// Send the email
	auth := LoginAuth(sender.fromEmail, sender.password)
	return smtp.SendMail(smtpServerAddress, auth, sender.fromEmail, append(to, cc...), msg.Bytes())
}

// Function to attach a file to the email
func attachFile(writer *multipart.Writer, filePath string) error {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(filePath)
	part, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":        {mime.TypeByExtension(filepath.Ext(filePath))},
		"Content-Disposition": {fmt.Sprintf("attachment; filename=\"%s\"", fileName)},
		"Content-Transfer-Encoding": {"base64"},
	})
	if err != nil {
		return err
	}

	part.Write([]byte(base64.StdEncoding.EncodeToString(fileData)))

	return nil
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown fromServer")
		}
	}
	return nil, nil
}
