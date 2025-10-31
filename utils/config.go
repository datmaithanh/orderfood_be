package utils

import "time"

const (
	DBDriver = "postgres"
	DBSource = "postgresql://root:secret@localhost:5432/orderfood?sslmode=disable"
	ServerAddress = ":8080"
	TokenSymmetricKey = "92947293748924923749237498262347"
	TokenDuration time.Duration = 15 * time.Minute
)