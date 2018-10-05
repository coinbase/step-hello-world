package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HelloWorld_StateMachine(t *testing.T) {
	stateMachine, err := StateMachine()
	assert.NoError(t, err)

	exec, err := stateMachine.Execute(&Hello{})

	assert.NoError(t, err)
	assert.Regexp(t, "Giday Mate", exec.LastOutputJSON)

	assert.Equal(t, []string{
		"Hello",
	}, exec.Path())
}
