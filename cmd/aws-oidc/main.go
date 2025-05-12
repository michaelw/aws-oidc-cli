package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

// CLI config using Kong
var CLI struct {
	Process struct {
		Provider  string `help:"OIDC provider name (as in config)" required:""`
		Role      string `help:"AWS Role ARN to assume" required:""`
		Account   string `help:"AWS Account ID" required:""`
		UseSecret bool   `help:"Use secret for code verifier (optional)"`
	} `cmd:"process" help:"Process OIDC flow and vend AWS credentials"`
	Config string `help:"Path to config file" default:"oidc-providers.json"`
}

// ProviderConfig holds API gateway URL for a provider
// (expand as needed for more provider config)
type ProviderConfig struct {
	Name   string `json:"name"`
	ApiURL string `json:"api_url"`
}

type Providers struct {
	Providers []ProviderConfig `json:"providers"`
}

type AwsCredsResponse struct {
	Version         int
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
	Expiration      time.Time
}

func main() {
	ctx := kong.Parse(&CLI)

	if ctx.Command() != "process" {
		ctx.PrintUsage(false)
		os.Exit(1)
	}

	// Load providers config
	file, err := os.Open(CLI.Config)
	if err != nil {
		log.Fatalf("failed to open config: %v", err)
	}
	defer file.Close()
	var providers Providers
	if err := json.NewDecoder(file).Decode(&providers); err != nil {
		log.Fatalf("failed to decode config: %v", err)
	}

	var provider *ProviderConfig
	for _, p := range providers.Providers {
		if p.Name == CLI.Process.Provider {
			provider = &p
			break
		}
	}
	if provider == nil {
		log.Fatalf("provider '%v' not found in config", CLI.Process.Provider)
	}

	// Start local server for redirect
	port := randomPort()
	redirectURI := fmt.Sprintf("http://127.0.0.1:%d/creds", port)
	server := &http.Server{Addr: ":" + strconv.Itoa(port)}

	codeCh := make(chan string)
	state := randomState()

	http.HandleFunc("/creds", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "invalid state", http.StatusBadRequest)
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "missing code", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Authentication complete. You may close this window.")
		codeCh <- code
	})

	// Start server in background
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Handle shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	// Begin OIDC flow (browser open, etc.)
	challenge, verifier := generatePKCE()
	// Construct OIDC auth URL (this would be provider-specific)
	_ = redirectURI
	authURL := fmt.Sprintf("%s/auth?challenge=%s&state=%s", strings.TrimSuffix(provider.ApiURL, "/"), challenge, state)
	fmt.Fprintf(os.Stderr, "Open the following URL in your browser to authenticate:\n  %s\n", authURL)
	// Open the URL in the default browser
	err = browser.OpenURL(authURL)
	if err != nil {
		log.Fatalf("failed to open URL: %v", err)
	}

	// Wait for code or interrupt
	var code string
	select {
	case code = <-codeCh:
		// got code
	case <-stop:
		log.Println("Interrupted")
		os.Exit(1)
	}

	log.Println("Login successful!")

	// Shutdown server
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = server.Shutdown(ctxTimeout)

	// Exchange code for credentials
	creds, err := exchangeCodeForCreds(provider.ApiURL, code, verifier, CLI.Process.Account, CLI.Process.Role)
	if err != nil {
		log.Fatalf("failed to get credentials: %v", err)
	}

	// Print credentials in AWS credential_process format
	output, _ := json.MarshalIndent(creds, "", "  ")
	fmt.Println(string(output))
}

// randomPort returns a random port between 49152–65535
func randomPort() int {
	return 8080 //49152 + int(time.Now().UnixNano()%int64(65535-49152))
}

// randomState returns a random string for OIDC state
func randomState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// generatePKCE returns a code_challenge and code_verifier (stub)
func generatePKCE() (challenge, verifier string) {
	verifier = oauth2.GenerateVerifier()
	challenge = oauth2.S256ChallengeFromVerifier(verifier)
	return
}

// exchangeCodeForCreds calls the /creds endpoint and returns credentials
func exchangeCodeForCreds(apiURL, code, verifier, account, role string) (*AwsCredsResponse, error) {
	// Compose request body
	body := map[string]string{
		"code":     code,
		"verifier": verifier,
		"account":  account,
		"role":     role,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	credsURL := fmt.Sprintf("%s/creds", apiURL)
	resp, err := http.Post(credsURL, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to POST to /creds: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("/creds error: %s", string(b))
	}

	var creds AwsCredsResponse
	if err := json.NewDecoder(resp.Body).Decode(&creds); err != nil {
		return nil, fmt.Errorf("failed to decode credentials: %w", err)
	}
	creds.Expiration = creds.Expiration.Local()
	return &creds, nil
}
