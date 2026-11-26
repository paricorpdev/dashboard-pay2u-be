package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/paricorpdev/dashboard-pay2u-be/internal/config"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: go run ./cmd/migrate <up|down|version|force> [arguments]")
	}

	viperConfig := config.NewViper()
	sslmode := viperConfig.GetString("database.postgresql.sslmode")
	if sslmode == "" {
		sslmode = "disable"
	}
	databaseURL := postgresURL(
		viperConfig.GetString("database.postgresql.username"),
		viperConfig.GetString("database.postgresql.password"),
		viperConfig.GetString("database.postgresql.host"),
		viperConfig.GetInt("database.postgresql.port"),
		viperConfig.GetString("database.postgresql.name"),
		sslmode,
	)

	migratePath, err := exec.LookPath("migrate")
	if err != nil {
		return fmt.Errorf("golang-migrate CLI not found in PATH; install it before running migrations")
	}

	commandArgs := normalizeCommandArgs(os.Args[1:])
	if len(commandArgs) > 1 && commandArgs[0] == "force" && commandArgs[1] == "0" {
		return fmt.Errorf("use `force -1` to reset before the first migration; version 0 does not exist")
	}

	args := []string{"-path", "db/migrations", "-database", databaseURL}
	args = append(args, commandArgs...)
	command := exec.Command(migratePath, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	if err := command.Run(); err != nil {
		return fmt.Errorf("migration command failed: %w", err)
	}
	return nil
}

func normalizeCommandArgs(input []string) []string {
	args := append([]string(nil), input...)
	if len(args) == 0 {
		return args
	}

	// Ignore a stray backslash pasted at the end of the command.
	args[len(args)-1] = strings.TrimSuffix(args[len(args)-1], `\`)

	// The migrate CLI parses subcommand arguments as flags, so stop flag parsing
	// before a negative version such as `force -1`.
	if len(args) > 1 && args[0] == "force" && strings.HasPrefix(args[1], "-") {
		args = append(args[:1], append([]string{"--"}, args[1:]...)...)
	}
	return args
}

func postgresURL(username, password, host string, port int, database, sslmode string) string {
	query := url.Values{}
	query.Set("sslmode", sslmode)
	// MS and Dashboard can point at the same PostgreSQL database, so keep
	// their golang-migrate version histories separate.
	query.Set("x-migrations-table", "dashboard_schema_migrations")
	return (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(username, password),
		Host:     net.JoinHostPort(host, fmt.Sprint(port)),
		Path:     "/" + database,
		RawQuery: query.Encode(),
	}).String()
}
