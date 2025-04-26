package server

import (
	"os"
)

type Config struct {
	port    string
	fileDir string
}

func NewConfig() *Config {
	return &Config{
		port:    getPort(),
		fileDir: getFileDir(),
	}
}

func (c *Config) Port() string {
	return c.port
}

func (c *Config) FileDir() string {
	return c.fileDir
}

func getPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "4221"
}

func getFileDir() string {
	for i, arg := range os.Args {
		if arg == "--directory" && i+1 < len(os.Args) {
			return os.Args[i+1]
		}
	}
	return "/tmp"
}
