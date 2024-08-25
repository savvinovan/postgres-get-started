package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"log/slog"
	"os"
	"sort"

	"github.com/jackc/pgx/v5"
)

func main() {
	connConfig, err := pgx.ParseConfig("postgres://postgres:postgres@localhost:5432")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}
	z, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create zap logger: %v\n", err)
		os.Exit(1)
	}
	// Example of using zap logger with slog.
	l := slog.New(slogzap.Option{Logger: z}.NewZapHandler())
	connConfig.Tracer = &tracelog.TraceLog{
		Logger:   &Logger{sl: l},
		LogLevel: 6,
	}
	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)

	dbPool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5432")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	var greeting2 string
	err = dbPool.QueryRow(context.Background(), "select 'Hello, world from pool!'").Scan(&greeting2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting2)

}

type Logger struct {
	sl *slog.Logger
}

// Log I grab this adapter from this issue:
// https://github.com/jackc/pgx/issues/1582
func (l *Logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var attrs []slog.Attr
	for _, k := range keys {
		attrs = append(attrs, slog.Any(k, data[k]))
	}
	l.sl.LogAttrs(ctx, translateLevel(level), msg, attrs...)
}

func translateLevel(level tracelog.LogLevel) slog.Level {
	switch level {
	case tracelog.LogLevelTrace:
		return slog.LevelDebug
	case tracelog.LogLevelDebug:
		return slog.LevelDebug
	case tracelog.LogLevelInfo:
		return slog.LevelInfo
	case tracelog.LogLevelWarn:
		return slog.LevelWarn
	case tracelog.LogLevelError:
		return slog.LevelError
	case tracelog.LogLevelNone:
		return slog.LevelError
	default:
		return slog.LevelError
	}
}
