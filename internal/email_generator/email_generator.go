package email_generator

import (
	"github.com/jordan-wright/email"
	"github.com/matcornic/hermes/v2"
	"net/smtp"
	"net/textproto"
	"os"
)

var h = hermes.Hermes{
	// Optional Theme
	// Theme: new(Default)
	Product: hermes.Product{
		// Appears in header & footer of e-mails
		Name: "Go-Rss",
		Link: "",
		// Optional product logo
		Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
	},
}

func SendEmail(emailAddr string, emailType hermes.Email) error {
	//TODO: Extract this to only call one
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := h.GenerateHTML(emailType)
	if err != nil {
		return err
	}

	e := &email.Email{
		To:      []string{emailAddr},
		From:    "Go-Rss<frankidatank@gmail.com>",
		Subject: "Welcome to Go-Rss",
		HTML:    []byte(emailBody),
		Headers: textproto.MIMEHeader{},
	}
	err = e.Send("smtp.gmail.com:587", smtp.PlainAuth("", username, password, "smtp.gmail.com"))
	if err != nil {
		return err
	}

	return nil
}

func GenerateWelcomeEmail(name string) hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: name,
			Intros: []string{
				"Welcome to Go-Rss! We're very excited to have you on board.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with Hermes, please click here:",
					Button: hermes.Button{
						Color: "#22BC66", // Optional action button color
						Text:  "Confirm your account",
						Link:  "https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
					},
				},
			},
		},
	}
}
