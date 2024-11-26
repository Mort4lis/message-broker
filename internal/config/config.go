package config

import (
	"time"
)

type HTTPServer struct {
	Listen string `yaml:"listen" env-default:":8080"`
}

type Queue struct {
	Name           string `yaml:"name" env-required:"true"`
	MaxMessages    int64  `yaml:"max_messages" env-default:"100"`
	MaxSubscribers int64  `yaml:"max_subscribers" env-default:"3"`
}

type Config struct {
	HTTPServer      HTTPServer    `yaml:"http_server"`
	Queues          []Queue       `yaml:"queues"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"10s"`
}
