package mail

import (
	"fmt"
	"go-skeleton/bootstrap"
	"go-skeleton/lib/utils"
	"log"

	mail "github.com/xhit/go-simple-mail/v2"
)

// Html filename
const (
	UserVerifyEmail    = "user_verify_registration"
	UserForgotPassword = "user_forgot_password"
	UserUpdateEmail    = "user_update_email"
)

var MailSubj = map[string]string{
	UserUpdateEmail:    "[Detect Data] Email Verification",
	UserVerifyEmail:    "[Detect Data] Email Verification",
	UserForgotPassword: "[Detect Data] Forgot Password",
}

type EmailData struct {
	Name        string
	Email       string
	Link        string
	Value       int
	Description string
}

type Contract struct {
	app *bootstrap.App
}

func New(app *bootstrap.App) *Contract {
	return &Contract{app}
}

func (c *Contract) SendMail(usedFor, subject, to string, emailData interface{}) error {
	fn := fmt.Sprintf("%s/%s.html", c.app.Config.GetString("resource_path"), usedFor)

	server := mail.NewSMTPClient()

	// SMTP Server
	server.Host = c.app.Config.GetString("mail.host")
	server.Port = c.app.Config.GetInt("mail.port")
	server.Username = c.app.Config.GetString("mail.username")
	server.Password = c.app.Config.GetString("mail.password")
	server.Encryption = mail.EncryptionSTARTTLS

	// SMTP client
	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	// fill the html body
	tpl, err := utils.ParseTpl(fn, emailData)
	if err != nil {
		return err
	}

	// New email simple html with inline and CC
	from := fmt.Sprintf("%s <%s>", c.app.Config.GetString("mail.mail_name"), c.app.Config.GetString("mail.mail_from"))
	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(to).
		SetSubject(subject)

	email.SetBody(mail.TextHTML, tpl)
	if email.Error != nil {
		return err
	}

	// Call Send and pass the client
	err = email.Send(smtpClient)
	if err != nil {
		return err
	} else {
		log.Println("Email Sent")
	}

	return nil
}
