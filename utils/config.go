package utils

import "time"

const (
	DBDriver                            = "postgres"
	// DBSource                            = "postgresql://root:secret@localhost:5432/orderfood?sslmode=disable"
	DBSource                            = "postgresql://neondb_owner:npg_DRjtm1F7ThoI@ep-round-violet-a1nlcl0m.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"
	ServerAddress                       = ":8888"
	TokenSymmetricKey                   = "92947293748924923749237498262347"
	TokenDuration         time.Duration = 15 * time.Minute
	RefreshTokenDuration  time.Duration = 7 * 24 * time.Hour
	UrlToWebsiteOrderFood               = "http://localhost:3000"
	CLOUDINARY_URL                      = "cloudinary://535384412283863:9kzrBUZX3R-Kh8-Z82Okh4UcDgQ@dxnuzxb59"
)
