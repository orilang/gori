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

func TestCommandsLower(t *testing.T) {
	assert := assert.New(t)

	t.Run("lower_success", func(t *testing.T) {
		configDir := "../testdata"
		configFile := filepath.Join(configDir, "lower/success/main.ori")

		cmd := Lower()
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan error, 1)
		go func() {
			done <- cmd.Run(ctx, []string{"lower", "--file", configFile})
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

	t.Run("lower_error_no_such_file_or_directory", func(t *testing.T) {
		configDir := "../testdata"
		configFile := filepath.Join(configDir, "main.ori")

		cmd := Lower()
		assert.Error(cmd.Run(context.Background(), []string{"lower", "--file", configFile}))
	})

	t.Run("lower_error_no_file_or_directory", func(t *testing.T) {
		cmd := Lower()
		assert.ErrorIs(walk.ErrNoFileOrDirectoryPassed, cmd.Run(context.Background(), []string{"lower"}))
	})

	t.Run("lower_error_no_such_file_or_directory", func(t *testing.T) {
		configDir := "../testdata/lower/errors"
		configFile := filepath.Join(configDir, "main.ori")

		cmd := Lower()
		assert.Error(cmd.Run(context.Background(), []string{"lower", "--file", configFile}))
	})

	t.Run("lower_test_data", func(t *testing.T) {
		workingDir, err := os.Getwd()
		assert.Nil(err)

		testdata := "../testdata/lower"
		err = filepath.Walk(filepath.Join(workingDir, testdata),
			func(file string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					cmd := Lower()
					if strings.Contains(file, "success") {
						assert.NoError(cmd.Run(context.Background(), []string{"lower", "--file", file}))
					} else {
						assert.Error(cmd.Run(context.Background(), []string{"lower", "--file", file}))
					}
				}
				return nil
			})
		assert.Nil(err)
	})
}
