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
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	formatter := logging.NewBackendFormatter(backend1, format)
	//formatter.SetLevel(logging.DEBUG, "")

	// Set the backends to be used.
	logging.SetBackend(formatter)
}