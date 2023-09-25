package internal

import "github.com/google/uuid"

func RandomSting() string {
	return uuid.Must(uuid.NewRandom()).String()
}
