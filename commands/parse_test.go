package commands

import (
	"context"
	"path/filepath"
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
}
