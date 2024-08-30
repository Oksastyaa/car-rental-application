package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
func SendEmail(to, subject, body string) error {
	// Get Mailtrap configuration from environment variables
	apiHost := os.Getenv("MAILTRAP_HOST") // Ini sekarang hanya host dasar
	apiKey := os.Getenv("MAILTRAP_API_TOKEN")
	fromEmail := os.Getenv("MAILTRAP_ADDRESS")
	fromName := os.Getenv("MAILTRAP_NAME")
	apiUrl := "https://sandbox.api.mailtrap.io/api/send/3108103"

	// Validate environment variables
	if apiHost == "" || apiKey == "" || fromEmail == "" {
		return fmt.Errorf("missing Mailtrap configuration: MAILTRAP_HOST, MAILTRAP_API_TOKEN, or MAILTRAP_ADDRESS is not set")
	}

	// Ensure the URL includes the scheme (https)

	// Construct the email request payload
	emailRequest := MailtrapEmailRequest{
		From: MailtrapAddress{
			Email: fromEmail,
			Name:  fromName,
		},
		To: []MailtrapAddress{
			{
				Email: to,
				Name:  to,
			},
		},
		Subject:  subject,
		HTMLBody: body,
	}

	// Convert request to JSON
	jsondata, err := json.Marshal(emailRequest)
	if err != nil {
		return fmt.Errorf("error when marshalling email request: %v", err)
	}

	// Create a new HTTP POST request to Mailtrap API
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsondata))
	if err != nil {
		return fmt.Errorf("error when creating request: %v", err)
	}

	// Add headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send request and get response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error when sending Email: %v", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to send email, status: %s", resp.Status)
	}

	fmt.Println("Email sent successfully")
	return nil
}

// Helper function to check if a URL starts with http or https
func startsWithHTTP(url string) bool {
	return len(url) >= 4 && (url[:4] == "http" || url[:5] == "https")
}
