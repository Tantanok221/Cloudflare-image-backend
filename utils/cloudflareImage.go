package utils

import "github.com/tantanok221/cloudflare-image-backend/internal/helper"

func GetEndpoint() string {
	id := helper.GetEnv("Cloudflare_AccountID")
	endpoint := "https://api.cloudflare.com/client/v4/accounts/" + id + "/images/v1"
	return endpoint
}

func GetAuthToken() string {
	token := "Bearer " + helper.GetEnv("Cloudflare_Token")
	return token
}

func IsImage(contentType string) bool {
	return contentType == "image/jpeg" || contentType == "image/png" || contentType == "image/gif"
}
