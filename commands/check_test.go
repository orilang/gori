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

func TestCommandsCheck(t *testing.T) {
	assert := assert.New(t)

	t.Run("type_checker_success", func(t *testing.T) {
		configDir := "../testdata"
		configFile := filepath.Join(configDir, "typecheck/success/main.ori")

		cmd := Check()
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan error, 1)
		go func() {
			done <- cmd.Run(ctx, []string{"lex", "--file", configFile})
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

	t.Run("type_checker_error_no_such_file_or_directory", func(t *testing.T) {
		configDir := "../testdata"
		configFile := filepath.Join(configDir, "main.ori")

		cmd := Check()
		assert.Error(cmd.Run(context.Background(), []string{"check", "--file", configFile}))
	})

	t.Run("type_checker_error_no_file_or_directory", func(t *testing.T) {
		cmd := Check()
		assert.ErrorIs(walk.ErrNoFileOrDirectoryPassed, cmd.Run(context.Background(), []string{"check"}))
	})

	t.Run("type_checker_error_no_such_file_or_directory", func(t *testing.T) {
		configDir := "../testdata/typecheck/errors"
		configFile := filepath.Join(configDir, "main.ori")

		cmd := Check()
		assert.Error(cmd.Run(context.Background(), []string{"check", "--file", configFile}))
	})

	t.Run("type_checker_test_data", func(t *testing.T) {
		workingDir, err := os.Getwd()
		assert.Nil(err)

		testdata := "../testdata/typecheck"
		err = filepath.Walk(filepath.Join(workingDir, testdata),
			func(file string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					cmd := Check()
					if strings.Contains(file, "success") {
						assert.NoError(cmd.Run(context.Background(), []string{"check", "--file", file}))
					} else {
						assert.Error(cmd.Run(context.Background(), []string{"check", "--file", file}))
					}
				}
				return nil
			})
		assert.Nil(err)
	})
}
