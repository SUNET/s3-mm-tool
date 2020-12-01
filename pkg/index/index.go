package index

import (
	"net/http"
	"os"
	"strings"
	"time"
)

var IndexServer string

func RegisterURL(url string) error {
	var client = &http.Client{
		Timeout: time.Second * 10,
	}

	body := strings.NewReader(url)
	_, err := client.Post(IndexServer, "application/x-www-form-urlencoded", body)
	return err
}

func init() {
	IndexServer = os.Getenv("S3_MM_INDEXSERVER")
}
