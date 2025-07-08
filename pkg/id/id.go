package id

import gonanoid "github.com/matoous/go-nanoid/v2"

const (
	base64URLSafe = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
	nanoIDSize    = 10
)

func GenerateID() (string, error) {
	id, err := gonanoid.Generate(base64URLSafe, nanoIDSize) // entropia 60 bits
	if err != nil {
		return "", err
	}
	return id, nil
}
