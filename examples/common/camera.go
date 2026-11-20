package common

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/andriantp/camera"
	"github.com/andriantp/camera/driver/isapi"
	"github.com/andriantp/camera/driver/onvif"
	"github.com/joho/godotenv"
)

func Camera() (*camera.Client, camera.Credential) {
	_, b, _, _ := runtime.Caller(0)
	path := filepath.Dir(b)
	filename := path + "/.env"

	if err := godotenv.Load(filename); err != nil {
		log.Fatalf("Load .env: %v", err)
	}

	cfg := camera.Config{
		Timeout:   Timeout,
		TLSConfig: tlsConfig(),
	}

	cred := camera.Credential{
		Host:     os.Getenv("Host"),
		Username: os.Getenv("User"),
		Password: os.Getenv("Password"),
	}

	return camera.New(cfg), cred
}

func ONVIF() (*onvif.Client, onvif.Camera) {
	client, cred := Camera()

	c := onvif.Camera{
		Credential: cred,
		Profile:    envString("Profile", DefaultProfile),
	}

	return onvif.New(client), c
}

func ISAPI() (*isapi.Client, isapi.Camera) {
	client, cred := Camera()

	c := isapi.Camera{
		Credential: cred,
		Channel:    envInt("Channel", DefaultChannel),
	}

	return isapi.New(client), c
}

// tlsConfig builds the TLS settings for https:// hosts from .env:
//
//	CACert   = path to the camera certificate (PEM) to trust
//	Insecure = true to skip certificate verification
//
// It returns nil (verify against system roots) when neither is set.
func tlsConfig() *tls.Config {
	caCert := os.Getenv("CACert")
	insecure := envBool("Insecure", false)

	if caCert == "" && !insecure {
		return nil
	}

	cfg := &tls.Config{}

	if caCert != "" {
		pem, err := os.ReadFile(caCert)
		if err != nil {
			log.Fatalf("read CACert: %v", err)
		}

		roots := x509.NewCertPool()
		if !roots.AppendCertsFromPEM(pem) {
			log.Fatalf("CACert %s: no PEM certificate found", caCert)
		}
		cfg.RootCAs = roots
	}

	if insecure {
		log.Println("warning: Insecure=true, TLS certificate verification is disabled")
		cfg.InsecureSkipVerify = true //nolint:gosec // opt-in via .env for self-signed cameras
	}

	return cfg
}

func envString(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("%s must be a number, got %q", key, v)
	}
	return n
}

func envBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Fatalf("%s must be true or false, got %q", key, v)
	}
	return b
}
