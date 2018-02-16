package log

import (
	"os"

	"github.com/op/go-logging"
)

var Log = logging.MustGetLogger("FrankLog")

var format = logging.MustStringFormatter(
	`%{color}%{level:.4s} %{time:15:04:05.000} %{shortfunc} ▶%{color:reset} %{message}`,
)

func InitLogger() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)

	formatter := logging.NewBackendFormatter(backend, format)
	//formatter.SetLevel(logging.DEBUG, "")

	// Set the backends to be used.
	logging.SetBackend(formatter)
}
