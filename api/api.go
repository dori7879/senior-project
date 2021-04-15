package api

import (
	"context"
	"math/rand"
)

// ReportError notifies an external service of errors. No-op by default.
var ReportError = func(ctx context.Context, err error, args ...interface{}) {}

// ReportPanic notifies an external service of panics. No-op by default.
var ReportPanic = func(err interface{}) {}

// Characters for random string generation
var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var digits = []rune("0123456789")

// RandStringSeq is a function to generate random sequence of chars of a set size
func RandStringSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// RandDigitSeq is a function to generate random sequence of digits of a set size
func RandDigitSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}
