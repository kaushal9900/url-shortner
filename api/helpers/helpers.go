package helpers

import (
	"strings"

	"github.com/kaushal9900/url-shortner/configs"
)

func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

func RemoveDomainError(url string) bool {
	domain := configs.EnvConfigs.Domain

	// Check if the URL is equal to the domain
	if url == domain {
		return false
	}

	// Remove "http://" or "https://"
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)

	// Remove "www."
	newURL = strings.Replace(newURL, "www.", "", 1)

	// Extract the domain from the remaining URL
	newURL = strings.Split(newURL, "/")[0]

	// Check if the extracted domain is equal to the domain
	return newURL != domain
}
