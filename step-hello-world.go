package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coinbase/step/handler"
	"github.com/coinbase/step/machine"
	"github.com/coinbase/step/utils/run"
)

func main() {
	var arg, command string
	switch len(os.Args) {
	case 1:
		fmt.Println("Starting Lambda")
		run.LambdaTasks(CreateTaskHandlers())
	case 2:
		command = os.Args[1]
		arg = ""
	case 3:
		command = os.Args[1]
		arg = os.Args[2]
	default:
		printUsage() // Print how to use and exit
	}

	switch command {
	case "json":
		run.JSON(StateMachine())
	case "exec":
		run.Exec(StateMachine())(&arg)
	default:
		printUsage() // Print how to use and exit
	}
}

func printUsage() {
	fmt.Println("Usage: step-hello-world <json|exec> <arg> (No args starts Lambda)")
	os.Exit(0)
}

// StateMachine returns the StateMachine
// replacing the `Resource` in Task states with the lambdaArn
func StateMachine() (*machine.StateMachine, error) {
	state_machine, err := machine.FromJSON([]byte(`{
    "Comment": "Hello World",
    "StartAt": "Hello",
    "States": {
      "Hello": {
				"Type": "TaskFn",
				"Resource": "arn:aws:lambda:{{aws_region}}:{{aws_account}}:function:{{lambda_name}}",
        "Comment": "Deploy Step Function",
        "End": true
      }
    }
  }`))

	if err != nil {
		return nil, err
	}

	return state_machine, nil
}

// CreateTaskHandlers returns
func CreateTaskHandlers() *handler.TaskHandlers {
	tm := handler.TaskHandlers{}
	tm["Hello"] = HelloHandler

	return &tm
}

////////////
// HANDLERS
////////////

type Hello struct {
	Greeting string
}

// Handlers must conform to function type
// func (context.Context, interface{}) (interface{}, error)
// The input is auto-unmarshalled and marsheled from and to JSON
func HelloHandler(_ context.Context, hello *Hello) (interface{}, error) {
	if hello.Greeting == "" {
		hello.Greeting = "Giday Mate"
	}
	return hello, nil
}
