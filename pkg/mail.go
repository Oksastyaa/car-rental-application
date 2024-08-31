package pkg

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"os"
)

// MailtrapEmailRequest represents the email request
type MailtrapEmailRequest struct {
	From     MailtrapAddress   `json:"from"`
	To       []MailtrapAddress `json:"to"`
	Subject  string            `json:"subject"`
	HTMLBody string            `json:"html"`
}

// MailtrapAddress represents the email address
type MailtrapAddress struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// SendEmail sends an email using Mailtrap API
//func SendEmail(to, subject, body string) error {
//	// Get Mailtrap configuration from environment variables
//	apiHost := os.Getenv("MAILTRAP_HOST") // Ini sekarang hanya host dasar
//	apiKey := os.Getenv("MAILTRAP_API_TOKEN")
//	fromEmail := os.Getenv("MAILTRAP_ADDRESS")
//	fromName := os.Getenv("MAILTRAP_NAME")
//	apiUrl := "https://send.api.mailtrap.io/api/send"
//
//	// Validate environment variables
//	if apiHost == "" || apiKey == "" || fromEmail == "" {
//		return fmt.Errorf("missing Mailtrap configuration: MAILTRAP_HOST, MAILTRAP_API_TOKEN, or MAILTRAP_ADDRESS is not set")
//	}
//
//	// Ensure the URL includes the scheme (https)
//
//	// Construct the email request payload
//	emailRequest := MailtrapEmailRequest{
//		From: MailtrapAddress{
//			Email: fromEmail,
//			Name:  fromName,
//		},
//		To: []MailtrapAddress{
//			{
//				Email: to,
//				Name:  to,
//			},
//		},
//		Subject:  subject,
//		HTMLBody: body,
//	}
//
//	// Convert request to JSON
//	jsondata, err := json.Marshal(emailRequest)
//	if err != nil {
//		return fmt.Errorf("error when marshalling email request: %v", err)
//	}
//
//	// Create a new HTTP POST request to Mailtrap API
//	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsondata))
//	if err != nil {
//		return fmt.Errorf("error when creating request: %v", err)
//	}
//
//	// Add headers
//	req.Header.Set("Authorization", "Bearer "+apiKey)
//	req.Header.Set("Content-Type", "application/json")
//
//	// Send request and get response
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		return fmt.Errorf("error when sending Email: %v", err)
//	}
//	defer resp.Body.Close()
//
//	// Check HTTP status code
//	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
//		return fmt.Errorf("failed to send email, status: %s", resp.Status)
//	}
//
//	fmt.Println("Email sent successfully")
//	return nil
//}

// SendEmail sends an email using SendGrid API
func SendEmail(to, subject, body, htmlContent string) error {
	from := mail.NewEmail("Car Rental", os.Getenv("SENDGRID_FROM_EMAIL"))
	toEmail := mail.NewEmail("User", to)
	message := mail.NewSingleEmail(from, subject, toEmail, body, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		return err
	}

	response, err := client.Send(message)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("failed to send email, status: %s", response.Body)
	}
	return nil
}
