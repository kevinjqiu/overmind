package overmind

import (
	"os"

	"github.com/go-kit/kit/log"
)

// Logger is used as the logger for the system
var Logger log.Logger

func init() {
	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger, "ts", log.DefaultTimestampUTC)
	Logger = log.With(Logger, "caller", log.DefaultCaller)
}
