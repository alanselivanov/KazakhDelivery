package mail

import (
	"fmt"
	"net/smtp"
	"user-service/internal/config"
)

type MailService struct {
	config *config.SMTPConfig
	auth   smtp.Auth
}

func NewMailService(cfg *config.Config) *MailService {
	auth := smtp.PlainAuth("", cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Host)

	return &MailService{
		config: &cfg.SMTP,
		auth:   auth,
	}
}

func (s *MailService) SendRegistrationConfirmation(to, username string) error {
	subject := "Registration Confirmation in KazakhDelivery"
	body := fmt.Sprintf(`
	<html>
		<body>
			<h2>Welcome to KazakhDelivery, %s!</h2>
			<p>Thank you for registering with our service. Your account has been successfully created.</p>
			<p>You can now log in using your email and password.</p>
			<p>Best regards,<br>The KazakhDelivery Team</p>
		</body>
	</html>
	`, username)

	return s.sendMail(to, subject, body)
}

func (s *MailService) sendMail(to, subject, htmlBody string) error {
	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.Username)

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)

	return smtp.SendMail(
		addr,
		s.auth,
		s.config.Username,
		[]string{to},
		[]byte(message),
	)
}
