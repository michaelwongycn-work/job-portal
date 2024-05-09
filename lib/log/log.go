package log

import (
	"context"
	"log"
)

func PrintLogErr(ctx context.Context, errorMsg string, err error) {
	log.Printf("%+v\n%s: %v", ctx, errorMsg, err)
}

func PrintLogAPIErr(ctx context.Context, errorMsg string, statusCode int) {
	log.Printf("%+v\n%s: %d", ctx, errorMsg, statusCode)
}
