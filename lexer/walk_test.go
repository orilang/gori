package lexer

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer_walk(t *testing.T) {
	assert := assert.New(t)

	t.Run("err_no_files_found", func(t *testing.T) {
		_, err := walk(Config{})
		assert.ErrorIs(err, ErrNoFilesFound)
	})

	t.Run("err_no_such_file", func(t *testing.T) {
		_, err := walk(Config{File: "xxxx.ori"})
		assert.ErrorIs(err, syscall.Errno(2))
	})

	t.Run("err_no_such_directory", func(t *testing.T) {
		_, err := walk(Config{Directory: "xxxx"})
		assert.ErrorIs(err, syscall.Errno(2))
	})

	t.Run("err_empty_directory", func(t *testing.T) {
		_, err := walk(Config{Directory: "testdata/empty"})
		assert.ErrorIs(err, ErrNoFilesFound)
	})

	t.Run("err_extension_no_such_file", func(t *testing.T) {
		_, err := walk(Config{File: "testdata/fake/fake"})
		assert.ErrorIs(err, ErrNoFilesFound)
	})

	t.Run("success_file", func(t *testing.T) {
		_, err := walk(Config{File: "testdata/fake/fake.ori"})
		assert.Nil(err)
	})

	t.Run("success_directory", func(t *testing.T) {
		_, err := walk(Config{Directory: "testdata/fake"})
		assert.Nil(err)
	})
}
