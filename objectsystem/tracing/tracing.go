package tracing

import (
	"os"
)

func IsDebugMode() bool {
	return os.Getenv("DEBUG") == "true"
}
