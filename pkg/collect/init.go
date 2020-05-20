package collect

import (
	"os"

	logging "gopkg.in/op/go-logging.v1"
)

func init() {
	// initialing the logger is necessary for yq to not print
	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05} %{shortfunc} [%{level:.4s}]%{color:reset} %{message}`,
	)
	backend := logging.AddModuleLevel(logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format))
	backend.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend)
}
