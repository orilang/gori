package walk

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer_Walk(t *testing.T) {
	assert := assert.New(t)

	t.Run("err_no_files_found", func(t *testing.T) {
		_, err := Walk(Config{})
		assert.ErrorIs(err, ErrNoFilesFound)
	})

	t.Run("err_no_such_file", func(t *testing.T) {
		_, err := Walk(Config{File: "xxxx.ori"})
		assert.ErrorIs(err, syscall.Errno(2))
	})

	t.Run("err_no_such_directory", func(t *testing.T) {
		_, err := Walk(Config{Directory: "xxxx"})
		assert.ErrorIs(err, syscall.Errno(2))
	})

	t.Run("err_empty_directory", func(t *testing.T) {
		workingDir, err := os.Getwd()
		assert.Nil(err)

		_, err = Walk(Config{Directory: filepath.Join(workingDir, "..", "testdata/empty")})
		assert.ErrorIs(err, ErrNoFilesFound)
	})

	t.Run("err_extension_no_such_file", func(t *testing.T) {
		workingDir, err := os.Getwd()
		assert.Nil(err)

		_, err = Walk(Config{File: filepath.Join(workingDir, "..", "testdata/fake/fake")})
		assert.ErrorIs(err, ErrNoFilesFound)
	})

	t.Run("success_file", func(t *testing.T) {
		workingDir, err := os.Getwd()
		assert.Nil(err)

		_, err = Walk(Config{File: filepath.Join(workingDir, "..", "testdata/fake/fake.ori")})
		assert.Nil(err)
	})

	t.Run("success_directory", func(t *testing.T) {
		workingDir, err := os.Getwd()
		assert.Nil(err)

		_, err = Walk(Config{Directory: filepath.Join(workingDir, "..", "testdata/fake")})
		assert.Nil(err)
	})
}
