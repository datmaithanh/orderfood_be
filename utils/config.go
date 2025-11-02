package utils

import "time"

const (
	DBDriver                            = "postgres"
	DBSource                            = "postgresql://root:secret@localhost:5432/orderfood?sslmode=disable"
	ServerAddress                       = ":8080"
	TokenSymmetricKey                   = "92947293748924923749237498262347"
	TokenDuration         time.Duration = 15 * time.Minute
	UrlToWebsiteOrderFood               = "http://localhost:3000"
	CLOUDINARY_URL                      = "cloudinary://535384412283863:9kzrBUZX3R-Kh8-Z82Okh4UcDgQ@dxnuzxb59"
)
