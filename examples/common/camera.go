package common

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

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
		Timeout: Timeout,
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
		Profile:    Profile,
	}

	return onvif.New(client), c
}

func ISAPI() (*isapi.Client, isapi.Camera) {
	client, cred := Camera()

	c := isapi.Camera{
		Credential: cred,
		Channel:    Channel,
	}

	return isapi.New(client), c
}
