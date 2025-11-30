package utils

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	DBDriver                            = "postgres"
	ServerAddress                       = ":8888"
	GrpcServerAddress                   = ":9090"
	TokenDuration         time.Duration = 15 * time.Minute
	RefreshTokenDuration  time.Duration = 7 * 24 * time.Hour
	UrlToWebsiteOrderFood               = "http://localhost:3000"
	CLOUDINARY_URL                      = "cloudinary://535384412283863:9kzrBUZX3R-Kh8-Z82Okh4UcDgQ@dxnuzxb59"
	END_POINT                          = "https://api.datmt.id.vn"
)

var (
	DBSource          = getDBSource()
	TokenSymmetricKey = getTokenSymmetricKey()
	Redis_Addr        = getRedisAddr()
	Redis_Password    = getRedisPassword()
	Redis_ServerName  = getRedisServerName()
)

var envLoaded = false

func init() {
	if !envLoaded {
		LoadConfig()
		envLoaded = true
	}
}

func getDBSource() string {
	if dbSource := os.Getenv("DB_SOURCE"); dbSource != "" {
		return dbSource
	}
	return ""
}

func getTokenSymmetricKey() string {
	if key := os.Getenv("TokenSymmetricKey"); key != "" {
		return key
	}
	return ""
}

func getRedisAddr() string {
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		return addr
	}
	return ""
}

func getRedisPassword() string {
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		return password
	}
	return ""
}

func getRedisServerName() string {
	if serverName := os.Getenv("REDIS_SERVER_NAME"); serverName != "" {
		return serverName
	}
	return ""
}

func LoadConfig() {
	godotenv.Load(".env.prod")
	DBSource = getDBSource()
	TokenSymmetricKey = getTokenSymmetricKey()
	Redis_Addr = getRedisAddr()
	Redis_Password = getRedisPassword()
	Redis_ServerName = getRedisServerName()
}
