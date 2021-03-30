package email

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

type Credentials struct {
	Email    string
	Password string
}

type Smtp struct {
	Host string
	Port string
}

type StudentsCsv struct {
	CarneIndex int
	EmailIndex int
}

type EmailOptions struct {
	Credentials Credentials
	StudentsCsv StudentsCsv
	Smtp        Smtp
	Subject     string
	Body        string
}

func SendEmail(opt EmailOptions, to []string) {
	// Sender data.
	from := opt.Credentials.Email
	password := opt.Credentials.Password

	// smtp server configuration.
	smtpHost := opt.Smtp.Host
	smtpPort := opt.Smtp.Port

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n%s", opt.Subject, mimeHeaders, opt.Body)))

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}

// Code from https://gist.github.com/guinso/8405a991d8a095b01427b9ea83934d67
func SendEmailWithAttachment(opt EmailOptions, tos []string, attachmentFilePath string, filename string) {
	var delimeter = "**=myohmy689407924327"

	tlsConfig := tls.Config{
		ServerName:         opt.Smtp.Host,
		InsecureSkipVerify: true,
	}

	// log.Println("Establish TLS connection")
	conn, connErr := tls.Dial("tcp", fmt.Sprintf("%s:%s", opt.Smtp.Host, opt.Smtp.Port), &tlsConfig)
	if connErr != nil {
		log.Panic(connErr)
	}
	defer conn.Close()

	// log.Println("create new email client")
	client, clientErr := smtp.NewClient(conn, opt.Smtp.Host)
	if clientErr != nil {
		log.Panic(clientErr)
	}
	defer client.Close()

	// log.Println("setup authenticate credential")
	auth := smtp.PlainAuth("", opt.Credentials.Email, opt.Credentials.Password, opt.Smtp.Host)

	if err := client.Auth(auth); err != nil {
		log.Panic(err)
	}

	// log.Println("Start write mail content")
	// log.Println("Set 'FROM'")
	if err := client.Mail(opt.Credentials.Email); err != nil {
		log.Panic(err)
	}
	// log.Println("Set 'TO(s)'")
	for _, to := range tos {
		if err := client.Rcpt(to); err != nil {
			log.Panic(err)
		}
	}

	writer, writerErr := client.Data()
	if writerErr != nil {
		log.Panic(writerErr)
	}

	// Basic email headers
	sampleMsg := fmt.Sprintf("From: %s\r\n", opt.Credentials.Email)
	sampleMsg += fmt.Sprintf("To: %s\r\n", strings.Join(tos, ";"))

	sampleMsg += fmt.Sprintf("Subject: %s\r\n", opt.Subject)

	// log.Println("Mark content to accept multiple contents")
	sampleMsg += "MIME-Version: 1.0\r\n"
	sampleMsg += fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", delimeter)

	// Place HTML message
	// log.Println("Put HTML message")
	sampleMsg += fmt.Sprintf("\r\n--%s\r\n", delimeter)
	sampleMsg += "Content-Type: text/html; charset=\"utf-8\"\r\n"
	sampleMsg += "Content-Transfer-Encoding: 7bit\r\n"
	sampleMsg += fmt.Sprintf("\r\n%s", opt.Body)

	// Place file
	// log.Println("Put file attachment")
	sampleMsg += fmt.Sprintf("\r\n--%s\r\n", delimeter)
	sampleMsg += "Content-Type: text/plain; charset=\"utf-8\"\r\n"
	sampleMsg += "Content-Transfer-Encoding: base64\r\n"
	sampleMsg += "Content-Disposition: attachment;filename=\"" + filename + "\"\r\n"

	// Read file
	rawFile, fileErr := ioutil.ReadFile(attachmentFilePath)
	if fileErr != nil {
		log.Panic(fileErr)
	}
	sampleMsg += "\r\n" + base64.StdEncoding.EncodeToString(rawFile)

	// Write into email client stream writter
	// log.Println("Write content into client writter I/O")
	if _, err := writer.Write([]byte(sampleMsg)); err != nil {
		log.Panic(err)
	}

	if closeErr := writer.Close(); closeErr != nil {
		log.Panic(closeErr)
	}

	client.Quit()

	log.Println("Email sent to: " + tos[0])
}

func GetEmailByCarne(carne string, records [][]string, scsv StudentsCsv) (string, error) {
	for _, i := range records {
		if strings.Compare(i[scsv.CarneIndex], carne) == 0 {
			return i[scsv.EmailIndex], nil
		}
	}

	return "", errors.New("El siguiente carnet no existe en el csv de alumnos: " + carne)
}

func ReadStudentsCsv() [][]string {
	b, err := os.ReadFile("Alumnos.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(strings.NewReader(string(b)))
	r.Comma = ','

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return records
}

func EmailOnly(opt EmailOptions) {
	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	students := ReadStudentsCsv()
	for _, file := range entries {
		filename := file.Name()
		ext := filepath.Ext(filename)
		if strings.Compare(ext, ".pdf") == 0 {
			carnet := strings.Replace(filename, ext, "", 1)

			to, err := GetEmailByCarne(carnet, students, opt.StudentsCsv)

			if err != nil {
				fmt.Println(err)
				continue
			}
			SendEmailWithAttachment(opt, []string{to}, filename, filename)
		}
	}
}
