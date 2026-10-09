package commands

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/orilang/gori/walk"
	"github.com/stretchr/testify/assert"
)

func TestCommandsParse(t *testing.T) {
	assert := assert.New(t)

	t.Run("success", func(t *testing.T) {
		configDir := "../testdata"
		configFile := filepath.Join(configDir, "success/main.ori")

		cmd := Parse()
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan error, 1)
		go func() {
			done <- cmd.Run(ctx, []string{"lex", "--file", configFile, "--output"})
		}()

		time.Sleep(time.Second)
		cancel()

		select {
		case err := <-done:
			assert.NoError(err)
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for Run() to stop")
		}
	})

	t.Run("error_no_such_file_or_directory", func(t *testing.T) {
		configDir := "../testdata"
		configFile := filepath.Join(configDir, "main.ori")

		cmd := Parse()
		assert.Error(cmd.Run(context.Background(), []string{"lex", "--file", configFile}))
	})

	t.Run("error_no_file_or_directory", func(t *testing.T) {
		cmd := Parse()
		assert.ErrorIs(walk.ErrNoFileOrDirectoryPassed, cmd.Run(context.Background(), []string{"lex"}))
	})

	t.Run("parser_error_no_such_file_or_directory", func(t *testing.T) {
		configDir := "../testdata/parser/errors"
		configFile := filepath.Join(configDir, "main.ori")

		cmd := Parse()
		assert.Error(cmd.Run(context.Background(), []string{"parse", "--file", configFile}))
	})

	t.Run("parser_test_data", func(t *testing.T) {
		workingDir, err := os.Getwd()
		assert.Nil(err)

		testdata := "../testdata/parser"
		err = filepath.Walk(filepath.Join(workingDir, testdata),
			func(file string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					cmd := Parse()
					if strings.Contains(file, "success") {
						assert.NoError(cmd.Run(context.Background(), []string{"parse", "--file", file}))
					} else {
						assert.Error(cmd.Run(context.Background(), []string{"parse", "--file", file}))
					}
				}
				return nil
			})
		assert.Nil(err)
	})
}
