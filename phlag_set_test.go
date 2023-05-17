package phlags

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syke99/phlags/internal"
	"os"
	"strconv"
	"testing"
)

func TestNewSet(t *testing.T) {
	// Act
	greetingSet := NewSet(internal.GreetSet)

	// Assert
	assert.NotNil(t, greetingSet)
}

func TestPhlagSet_AddPhlag(t *testing.T) {
	// Arrange
	politeSet := NewSet(internal.PoliteSet)
	hello, _ := New(internal.HelloAbrv, internal.HelloFull, internal.TestUsage, internal.TestDefault)

	// Act
	politeSet.AddPhlag(hello)

	// Assert
	assert.Equal(t, 1, len(politeSet.set))
}

func TestPhlagSet_Parse_SingleSet(t *testing.T) {
	// Arrange
	politeSet := NewSet(internal.PoliteSet)
	hello, _ := New(internal.HelloAbrv, internal.HelloFull, internal.TestUsage, internal.TestDefault)
	goodbye, _ := New(internal.GoodbyeAbrv, internal.GoodbyeFull, internal.TestUsage, internal.TestDefault)

	politeSet.AddPhlag(hello).AddPhlag(goodbye)

	questionSet := NewSet(internal.QuestionsSet)
	who, _ := New(internal.WhoAbrv, internal.WhoFull, internal.TestUsage, internal.TestDefault)
	what, _ := New(internal.WhatAbrv, internal.WhatFull, internal.TestUsage, internal.TestDefault)

	questionSet.AddPhlag(who).AddPhlag(what)

	os.Args = []string{"", internal.PoliteSet, fmt.Sprintf("%s=%s", internal.HelloAbrv, internal.HelloWorld), strconv.Itoa(internal.One), strconv.Itoa(internal.Two), strconv.Itoa(internal.Three), fmt.Sprintf("%s=%s", internal.GoodbyeAbrv, internal.GoodbyeAll), strconv.Itoa(internal.Four), strconv.Itoa(internal.Five), strconv.Itoa(internal.Six)}

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

func TestPhlagSet_Parse_MultipleSets(t *testing.T) {
	// Arrange
	politeSet := NewSet(internal.PoliteSet)
	hello, _ := New(internal.HelloAbrv, internal.HelloFull, internal.TestUsage, internal.TestDefault)
	goodbye, _ := New(internal.GoodbyeAbrv, internal.GoodbyeFull, internal.TestUsage, internal.TestDefault)

	politeSet.AddPhlag(hello).AddPhlag(goodbye)

	questionSet := NewSet(internal.QuestionsSet)
	who, _ := New(internal.WhoAbrv, internal.WhoFull, internal.TestUsage, internal.TestDefault)
	what, _ := New(internal.WhatAbrv, internal.WhatFull, internal.TestUsage, internal.TestDefault)

	questionSet.AddPhlag(who).AddPhlag(what)

	os.Args = []string{"", internal.PoliteSet, fmt.Sprintf("%s=%s", internal.HelloAbrv, internal.HelloWorld), strconv.Itoa(internal.One), strconv.Itoa(internal.Two), strconv.Itoa(internal.Three), fmt.Sprintf("%s=%s", internal.GoodbyeAbrv, internal.GoodbyeAll), strconv.Itoa(internal.Four), strconv.Itoa(internal.Five), strconv.Itoa(internal.Six), internal.QuestionsSet, fmt.Sprintf("%s=%s", internal.WhoAbrv, internal.You), fmt.Sprintf("%s=%s", internal.WhatAbrv, internal.Testing)}

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
	assert.Equal(t, internal.You, *who.String())
	assert.Equal(t, internal.Testing, *what.String())
}

func TestPhlagSet_Parse_BasePhlags(t *testing.T) {
	// Arrange
	politeSet := NewSet(internal.PoliteSet)
	hello, _ := New(internal.HelloAbrv, internal.HelloFull, internal.TestUsage, internal.TestDefault)
	goodbye, _ := New(internal.GoodbyeAbrv, internal.GoodbyeFull, internal.TestUsage, internal.TestDefault)

	politeSet.AddPhlag(hello).AddPhlag(goodbye)

	questionSet := NewSet(internal.QuestionsSet)
	who, _ := New(internal.WhoAbrv, internal.WhoFull, internal.TestUsage, internal.TestDefault)
	what, _ := New(internal.WhatAbrv, internal.WhatFull, internal.TestUsage, internal.TestDefault)

	questionSet.AddPhlag(who).AddPhlag(what)

	base, _ := New(internal.BaseAbrv, internal.BaseFull, internal.TestUsage, internal.TestDefault)

	AddBaseSetPhlag(base)

	os.Args = []string{"", fmt.Sprintf("%s=%s", internal.BaseAbrv, internal.Success), strconv.Itoa(internal.Seven), strconv.Itoa(internal.Eight), strconv.Itoa(internal.Nine), internal.PoliteSet, fmt.Sprintf("%s=%s", internal.HelloAbrv, internal.HelloWorld), strconv.Itoa(internal.One), strconv.Itoa(internal.Two), strconv.Itoa(internal.Three), fmt.Sprintf("%s=%s", internal.GoodbyeAbrv, internal.GoodbyeAll), strconv.Itoa(internal.Four), strconv.Itoa(internal.Five), strconv.Itoa(internal.Six)}

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
	assert.Equal(t, internal.Success, *base.String())
	assert.Equal(t, internal.Seven, *base.args[0].Int())
	assert.Equal(t, internal.Eight, *base.args[1].Int())
	assert.Equal(t, internal.Nine, *base.args[2].Int())
}
