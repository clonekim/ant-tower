package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Port         string
	UseConsole   bool
	DBFile       string
	LogFile      string
	LogRetention int
	HomeDir      string
}

func Load() *Config {
	cfg := &Config{}

	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}

	if home == "" {
		var err error
		home, err = os.UserHomeDir()
		if err != nil {
			panic(fmt.Sprintf("FATAL: Cannot find home directory. Please check HOME/USERPROFILE env vars: %v", err))
		}
	}

	cfg.HomeDir = home

	flag.StringVar(&cfg.Port, "port", "5001", "Server Port")
	flag.BoolVar(&cfg.UseConsole, "console", false, "Show console window (false to hide)")

	flag.StringVar(&cfg.DBFile, "db", "twn.db", "SQLite Database path")
	flag.StringVar(&cfg.LogFile, "logfile", "trace.log", "Log filename (saved in HomeDir/trace.log)")
	flag.IntVar(&cfg.LogRetention, "retention", 30, "Log retention days")

	flag.Parse()

	cfg.DBFile = filepath.Join(home, cfg.DBFile)

	if envPort := os.Getenv("TWN_PORT"); envPort != "" {
		cfg.Port = envPort
	}

	return cfg
}
