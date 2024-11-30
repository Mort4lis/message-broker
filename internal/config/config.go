package config

import (
	"time"
)

type Logging struct {
	Level  string `env-default:"info" yaml:"level"`
	Format string `env-default:"text" yaml:"format"`
}

type HTTPServer struct {
	Listen            string        `env-default:":8080" yaml:"listen"`
	ReadTimeout       time.Duration `env-default:"10s"   yaml:"read_timeout"`
	ReadHeaderTimeout time.Duration `env-default:"5s"    yaml:"read_header_timeout"`
	WriteTimeout      time.Duration `env-default:"10s"   yaml:"write_timeout"`
}

type Queue struct {
	Name           string `env-required:"true" yaml:"name"`
	MaxMessages    int64  `env-default:"100"   yaml:"max_messages"`
	MaxSubscribers int64  `env-default:"3"     yaml:"max_subscribers"`
}

type Config struct {
	Logging         Logging       `yaml:"logging"`
	Queues          []Queue       `yaml:"queues"`
	HTTPServer      HTTPServer    `yaml:"http_server"`
	ShutdownTimeout time.Duration `env-default:"10s"  yaml:"shutdown_timeout"`
}
