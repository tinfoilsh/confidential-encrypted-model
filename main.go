package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	listenAddr     = flag.String("l", ":8080", "listen address")
	model          = flag.String("m", "quay.io/nates/qwen3-0.6b:encrypted", "model to download")
	destinationDir = flag.String("d", "decrypted", "destination directory")
)

func main() {
	flag.Parse()

	status := "waiting for key"

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(status + "\n"))
	})

	http.HandleFunc("/key", func(w http.ResponseWriter, r *http.Request) {
		if status != "waiting for key" {
			http.Error(w, "key already received", http.StatusConflict)
			return
		}

		key, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tmpFile, err := os.CreateTemp("", "key.pem")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpFile.Write(key)
		tmpFile.Close()

		status = "decrypting"

		args := []string{
			"skopeo", "copy", "--insecure-policy",
			"--decryption-key", tmpFile.Name(),
			"docker://" + *model,
			"dir:" + *destinationDir,
		}
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run skopeo command in background
		go func() {
			defer os.Remove(tmpFile.Name())
			if err := cmd.Run(); err != nil {
				status = "failed: " + err.Error()
			} else {
				status = "ready"
			}
		}()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("key loaded"))
	})

	log.Println("Starting server on", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
