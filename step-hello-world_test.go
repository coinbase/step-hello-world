package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HelloWorld_StateMachine(t *testing.T) {
	stateMachine, err := StateMachine()
	assert.NoError(t, err)

	err = stateMachine.SetTaskFnHandlers(CreateTaskHandlers())
	assert.NoError(t, err)


	exec, err := stateMachine.Execute(&Hello{})
	assert.NoError(t, err)
	assert.Equal(t, "Giday Mate", exec.Output["Greeting"])

	assert.Equal(t, exec.Path(), []string{
		"Hello",
	})
}
