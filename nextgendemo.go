package main

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

var (
	LengthErr  = 4
	LettersErr = []byte("0123456789")
)

var ErrNoRowsFound error = errors.New("No rows found.")

type MyError struct {
	internal error
	external error
	location string
	code     string
}

func New(internal error, external string) error {
	b := make([]byte, LengthErr)
	for i := range b {
		b[i] = LettersErr[rand.Intn(len(LettersErr))]
	}
	_, file, line, _ := runtime.Caller(1)
	return MyError{
		internal: internal,
		external: errors.New(external),
		location: fmt.Sprintf("%s:%d", file, line),
		code:     string(b),
	}
}

func (m MyError) Error() string {
	if m.internal == nil {
		m.internal = errors.New("unspecified internal error")
	}
	if m.external == nil {
		m.external = errors.New("unknown error")
	}
	return fmt.Sprintf("error: %s [e%s]\n\t%s\n\t%s\n",
		m.external, m.code, m.internal, m.location)
}

func (m MyError) Unwrap() error {
	return m.internal
}

func FunctionCausingError() error {
	err := fmt.Errorf("No userid found in users table: %w", ErrNoRowsFound)
	return New(err, "invalid userid")
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	err := FunctionCausingError()
	if err != nil {
		if errors.Is(err, ErrNoRowsFound) {
			fmt.Print(err)
		}
	}
}
