package server

import (
	"os"
	"strings"
)

func EnforceHTTPS(url string) string {

	if !strings.Contains(url, "://") {
		return "https://" + url
	}
	return strings.Replace(url, "http://", "https://", 1)
}

func RemoveProhibitedUrls(url string) bool {
	if url == os.Getenv("DOMAIN") || url == os.Getenv("DOMAIN")+"/" {
		return false
	}

	// Remove protocol (https:// or http://)
	newUrl := strings.TrimPrefix(url, "https://")
	newUrl = strings.TrimPrefix(newUrl, "http://")

	// Remove www. prefix
	newUrl = strings.TrimPrefix(newUrl, "www.")

	// Get only the domain part (before first /)
	newUrl = strings.Split(newUrl, "/")[0]

	if newUrl == os.Getenv("DOMAIN") || newUrl == os.Getenv("DOMAIN")+"/" {
		return false
	}

	return true
}
