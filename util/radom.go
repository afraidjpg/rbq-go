package util

import "github.com/google/uuid"

func RandomName() string {
	return uuid.Must(uuid.NewRandom()).String()
}
