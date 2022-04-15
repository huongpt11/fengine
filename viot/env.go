package viot

import (
	_logger "github.com/duclmse/fengine/pkg/logger"
	"github.com/goccy/go-json"
	"github.com/subosito/gotenv"
	"io"
	"net/http"
	"os"
)

const version string = "0.10.0"

// VersionInfo contains version endpoint response.
type VersionInfo struct {
	// Service contains service name.
	Service string `json:"service"`
	// Version contains service current version value.
	Version string `json:"version"`
}

// Version exposes an HTTP handler for retrieving service version.
func Version(service string) http.HandlerFunc {
	return func(rw http.ResponseWriter, _ *http.Request) {
		res := VersionInfo{service, version}

		data, _ := json.Marshal(res)

		_, _ = rw.Write(data)
	}
}

// Env reads specified environment variable. If no value has been found, fallback is returned.
func Env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

//LoadEnvFile loads environment variables defined in an .env formatted file.
func LoadEnvFile(envFilePath string) error {
	return gotenv.Load(envFilePath)
}

type UUIDProvider interface {
	// ID generates the unique identifier.
	ID() (string, error)
}

// Response contains HTTP response specific methods.
type Response interface {
	// Code returns HTTP response code.
	Code() int

	// Headers returns map of HTTP headers with their values.
	Headers() map[string]string

	// Empty indicates if HTTP response has content.
	Empty() bool
}

func Close(logger _logger.Logger, name string) func(io.Closer) {
	return func(closer io.Closer) {
		err := closer.Close()
		if err != nil {
			logger.Error("cannot close %s: %s", name, err.Error())
		}
	}
}
