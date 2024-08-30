package config

import (
	"github.com/xendit/xendit-go/v6"
	"os"
)

var XenditClient *xendit.APIClient

func InitXendit() {
	apiKey := os.Getenv("XENDIT_API_KEY")
	XenditClient = xendit.NewClient(apiKey)
}
