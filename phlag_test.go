package phlags

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syke99/phlags/internal"
	"os"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	// Act
	flg, _ := New("-h", "--hello", internal.TestUsage, internal.TestDefault)

	// Assert
	assert.NotNil(t, flg)
}

func TestNew_MissingNames(t *testing.T) {
	// Act
	flg, _ := New("", "", internal.TestUsage, internal.TestDefault)

	// Assert
	assert.Nil(t, flg)
}

func TestParse(t *testing.T) {
	// Arrange
	hello, _ := New(internal.HelloAbrv, internal.HelloFull, internal.TestUsage, internal.TestDefault)
	goodbye, _ := New(internal.GoodbyeAbrv, internal.GoodbyeFull, internal.TestUsage, internal.TestDefault)

	os.Args = []string{"", fmt.Sprintf("%s=%s", internal.HelloAbrv, internal.HelloWorld), strconv.Itoa(internal.One), strconv.Itoa(internal.Two), strconv.Itoa(internal.Three), fmt.Sprintf("%s=%s", internal.GoodbyeAbrv, internal.GoodbyeAll), strconv.Itoa(internal.Four), strconv.Itoa(internal.Five), strconv.Itoa(internal.Six)}

	// Act
	err := Parse()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, internal.HelloWorld, *hello.String())
	assert.Equal(t, internal.One, *hello.args[0].Int())
	assert.Equal(t, internal.Two, *hello.args[1].Int())
	assert.Equal(t, internal.Three, *hello.args[2].Int())
	assert.Equal(t, internal.GoodbyeAll, *goodbye.String())
	assert.Equal(t, internal.Four, *goodbye.args[0].Int())
	assert.Equal(t, internal.Five, *goodbye.args[1].Int())
	assert.Equal(t, internal.Six, *goodbye.args[2].Int())
}
