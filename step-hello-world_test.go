package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HelloWorld_StateMachine(t *testing.T) {
	state_machine, err := StateMachine()
	assert.NoError(t, err)

	output, err := state_machine.ExecuteToMap(&Hello{})
	assert.NoError(t, err)
	assert.Equal(t, "Giday Mate", output["Greeting"])

	assert.Equal(t, state_machine.ExecutionPath(), []string{
		"HelloFn",
		"Hello",
	})
}
