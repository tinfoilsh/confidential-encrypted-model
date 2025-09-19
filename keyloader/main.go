package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/tinfoilsh/verifier/client"
)

var (
	enclave = flag.String("enclave", "encrypted-model.inf6.tinfoil.sh", "")
	repo    = flag.String("repo", "tinfoilsh/confidential-encrypted-model", "")
	key     = flag.String("key", "private.pem", "")
)

func main() {
	flag.Parse()

	log.Println("Verifying enclave")
	tinfoilClient := client.NewSecureClient(*enclave, *repo)
	httpClient, err := tinfoilClient.HTTPClient()
	if err != nil {
		log.Fatalf("failed to create HTTP client: %v", err)
	}

	keyFile, err := os.ReadFile(*key)
	if err != nil {
		log.Fatalf("failed to read key: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/key", *enclave), bytes.NewBuffer(keyFile))
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("failed to load key: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %v", err)
	}

	log.Println(string(body))
}
