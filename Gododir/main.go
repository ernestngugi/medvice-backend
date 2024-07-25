package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/godo.v2"
)

func main() {
	godo.Godo(tasks)
}

func tasks(p *godo.Project) {

	p.Task("test", nil, func(c *godo.Context) {
		envFilePath := c.Args.AsString("e")
		envMap, envStr := buildEnv(envFilePath)

		c.Bash(fmt.Sprintf("docker-compose exec -T postgres psql -c \"DROP DATABASE IF EXISTS %v;\" -U %v -d template1;", envMap["DATABASE_NAME"], envMap["DATABASE_USER"]))
		c.Bash(fmt.Sprintf("docker-compose exec -T postgres psql -c \"CREATE DATABASE %v\" -U %v -d template1;", envMap["DATABASE_NAME"], envMap["DATABASE_USER"]))
		c.Bash(fmt.Sprintf("goose -dir internal/db/migrations postgres %v up", envMap["DATABASE_URL"]))
		c.Bash(fmt.Sprintf("%v go test -race ./...", envStr))
	})

	p.Task("test-lite", nil, func(c *godo.Context) {
		envFilePath := c.Args.AsString("e")
		_, envStr := buildEnv(envFilePath)

		c.Bash(fmt.Sprintf("%v go test ./...", envStr))
	})

	p.Task("coverage", nil, func(c *godo.Context) {
		envFilePath := c.Args.AsString("e")
		_, envStr := buildEnv(envFilePath)

		c.Bash(fmt.Sprintf("%v go test -coverprofile=coverage.out ./...", envStr))
		c.Bash(fmt.Sprintf("%v go tool cover -html=coverage.out", envStr))
	})
}

func buildEnv(filePath string) (map[string]string, string) {

	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	envMap, err := godotenv.Parse(f)
	if err != nil {
		panic(err)
	}

	for key, val := range envMap {
		if key == "DATABASE_URL" {
			dbURL, err := url.Parse(val)
			if err != nil {
				panic(err)
			}

			envMap["DATABASE_USER"] = dbURL.User.Username()
			envMap["DATABASE_NAME"] = strings.TrimPrefix(dbURL.Path, "/")
		}
	}

	envs := make([]string, 0)
	for key, val := range envMap {
		envs = append(envs, fmt.Sprintf("%v=\"%v\"", key, val))
	}

	return envMap, strings.Join(envs, " ")
}
