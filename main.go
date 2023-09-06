package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation/v2"
	github "github.com/google/go-github/v53/github"
	"github.com/k0kubun/pp"
)

func run() error {
	appID, err := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse GITHUB_APP_ID: %w", err)
	}

	privateKey := []byte(os.Getenv("GITHUB_APP_PRIVATE_KEY"))
	ctx := context.Background()

	tr := http.DefaultTransport

	appsTransport, err := ghinstallation.NewAppsTransport(tr, appID, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create apps transport to discover installation ID: %w", err)
	}

	client := github.NewClient(&http.Client{Transport: appsTransport})
	installation, _, err := client.Apps.FindRepositoryInstallation(ctx, "grafana", "deployment_tools")
	if err != nil {
		return fmt.Errorf("failed to find installation ID: %w", err)
	}

	pp.Print("installation %#+v\n", installation)

	installationID := installation.GetID()

	// and then create a new transport for the installation so we can make API
	// calls with it
	clientTransport, err := ghinstallation.New(tr, appID, installationID, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create apps transport to get token: %w", err)
	}

	clientToken, err := clientTransport.Token(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	fmt.Printf("token: %s\n", clientToken)

	return nil

}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
