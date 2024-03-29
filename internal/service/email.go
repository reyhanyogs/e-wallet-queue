package service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"

	"github.com/hibiken/asynq"
	"github.com/reyhanyogs/e-wallet-queue/domain"
	"github.com/reyhanyogs/e-wallet-queue/dto"
	"github.com/reyhanyogs/e-wallet-queue/internal/component"
	"github.com/reyhanyogs/e-wallet-queue/internal/config"
)

type emailService struct {
	config *config.Config
}

func NewEmail(config *config.Config) domain.EmailService {
	return &emailService{
		config: config,
	}
}

func (s *emailService) Send(to string, subject string, body string) error {
	from := mail.Address{Name: "", Address: s.config.Mail.User}
	toMail := mail.Address{Name: "", Address: to}

	// Setup Headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = toMail.String()
	headers["Subject"] = subject

	// Setup Message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := s.config.Mail.Host + ":" + s.config.Mail.Port

	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", s.config.Mail.User, s.config.Mail.Password, host)

	// TLS Config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsConfig)
	if err != nil {
		component.Log.Errorf("Send(Dial): err = %s", err.Error())
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		component.Log.Errorf("Send(NewClient): err = %s", err.Error())
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		component.Log.Errorf("Send(Auth): err = %s", err.Error())
		return err
	}

	// From & To
	if err = c.Mail(from.Address); err != nil {
		component.Log.Errorf("Send(Mail): err = %s", err.Error())
		return err
	}

	if err = c.Rcpt(toMail.Address); err != nil {
		component.Log.Errorf("Send(Rcpt): err = %s", err.Error())
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		component.Log.Errorf("Send(Data): err = %s", err.Error())
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		component.Log.Errorf("Send(Write): err = %s", err.Error())
		return err
	}

	return w.Close()
}

func (s *emailService) SendEmailQueue() (string, func(ctx context.Context, task *asynq.Task) error) {
	return "send:email", func(ctx context.Context, task *asynq.Task) error {
		var data dto.EmailSendReq
		_ = json.Unmarshal(task.Payload(), &data)

		component.Log.Infof("SendEmailQueue(Send): To = %s", data.To)
		return s.Send(data.To, data.Subject, data.Body)
	}
}
