package phlags

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	// Act
	flg, _ := New("-h", "--hello", "test usage", "test default")

	// Assert
	assert.NotNil(t, flg)
}

func TestNew_MissingNames(t *testing.T) {
	// Act
	flg, _ := New("", "", "test usage", "test default")

	// Assert
	assert.Nil(t, flg)
}

func TestParse(t *testing.T) {
	// Arrange
	hello, _ := New("-h", "--hello", "test usage", "test default")
	goodbye, _ := New("-g", "--goodbye", "test usage", "test default")

	os.Args = []string{"", "-h=helloWorld", "1", "2", "3", "-g=goodbyeAll", "4", "5", "6"}

	// Act
	err := Parse()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "helloWorld", *hello.String())
	assert.Equal(t, 1, *hello.args[0].Int())
	assert.Equal(t, 2, *hello.args[1].Int())
	assert.Equal(t, 3, *hello.args[2].Int())
	assert.Equal(t, "goodbyeAll", *goodbye.String())
	assert.Equal(t, 4, *goodbye.args[0].Int())
	assert.Equal(t, 5, *goodbye.args[1].Int())
	assert.Equal(t, 6, *goodbye.args[2].Int())
}
