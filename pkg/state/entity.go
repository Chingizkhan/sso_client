package state

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

// State is used to pass forward. It's an random variable, for tracking user's oauth flow
type State string

func New(st string) (State, error) {
	if st == "" {
		return "", errors.New("incorrect state")
	}
	return State(st), nil
}

func Generate() (State, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("unable randomize state: %s", err)
	}
	state := base64.StdEncoding.EncodeToString(b)
	return State(state), nil
}
