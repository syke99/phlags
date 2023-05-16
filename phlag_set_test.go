package phlags

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewSet(t *testing.T) {
	// Act
	greetingSet := NewSet("greet")

	// Assert
	assert.NotNil(t, greetingSet)
}

func TestPhlagSet_AddPhlag(t *testing.T) {
	// Arrange
	politeSet := NewSet("polite")
	hello, _ := New("-h", "--hello", "test usage", "test default")

	// Act
	politeSet.AddPhlag(hello)

	// Assert
	assert.Equal(t, 1, len(politeSet.set))
}

func TestPhlagSet_Parse(t *testing.T) {
	// Arrange
	politeSet := NewSet("polite")
	hello, _ := New("-h", "--hello", "test usage", "test default")
	goodbye, _ := New("-g", "--goodbye", "test usage", "test default")

	politeSet.AddPhlag(hello).AddPhlag(goodbye)

	questionSet := NewSet("questions")
	who, _ := New("-wh", "--who", "test usage", "test default")
	what, _ := New("-wt", "--what", "test usage", "test default")

	questionSet.AddPhlag(who).AddPhlag(what)

	os.Args = []string{"", "polite", "-h=helloWorld", "1", "2", "3", "-g=goodbyeAll", "4", "5", "6"}

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
