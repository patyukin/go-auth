package main

import (
	"context"
	"log"

	"github.com/bogatyr285/auth-go/cmd/commands"
)

func main() {
	ctx := context.Background()

	cmd := commands.NewServeCmd()

	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Fatalf("smth went wrong: %s", err)
	}
}
