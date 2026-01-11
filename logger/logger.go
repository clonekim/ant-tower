package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"twn-monitor/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Setup(cfg *config.Config) {

	fileLogger := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.HomeDir, cfg.LogFile),
		MaxSize:    10, //10MB
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	fileWriter := zerolog.ConsoleWriter{
		Out:        fileLogger,
		NoColor:    true,
		TimeFormat: "2006-01-02 15:04:05",
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("| %s", i)
		},
	}

	var output io.Writer
	if cfg.UseConsole {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
		output = zerolog.MultiLevelWriter(fileWriter, consoleWriter)
	} else {
		output = fileWriter
	}

	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()

	log.Info().Str("path", cfg.HomeDir).Msg("Logger initialized")
}
