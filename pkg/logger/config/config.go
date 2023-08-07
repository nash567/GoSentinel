package config

import (
	"io"

	"github.com/nash567/GoSentinel/pkg/logger/model"
)

type Config struct {
	// Level is the lowest level of log message that should be emitted. Any log
	// messages logged at the specified level or any level more severe will be
	// emitted. The default level is DEBUG.
	Level model.Level

	// Output is the destination for log messages. By default, it is os.Stdout.
	Output io.Writer
}
