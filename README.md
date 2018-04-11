# step-hello-world

`step-hello-world` is an example application for the [Step](https://github.com/coinbase/step) library. It is a AWS Step Function and Lambda that takes an input, and returns a greeting.


### Getting started

Building with `go build .`, which will create the binary `step-hello-world`.

You can test its execution with:

```bash
./step-hello-world exec
{
  "Greeting": "Hello World"
}

./step-hello-world exec '{"Greeting": "Hi"}'
{
 "Greeting": "Hi"
}
```

Looking at the State Machine output:

```bash
./step-hello-world json
{
 "Comment": "Hello World",
 "StartAt": "HelloFn",
 "States": {
  "Hello": {
   "Type": "Task",
   "Comment": "Deploy Step Function",
   "Resource": "coinbase-step-hello-world",
   "End": true
  },
  "HelloFn": {
   "Type": "Pass",
   "ResultPath": "$.Task",
   "Result": "Hello",
   "Next": "Hello"
  }
 }
}
```

### Deploying

To create the AWS resources we use [GeoEngineer](https://github.com/coinbase/geoengineer). This requires `ruby`, `bundler` and `terraform`. 

Create the resources with:

```
bundle
./geo apply resources/step-hello-world.rb
```

Once the resources are created you can bootstrap or deploy `step-hello-world` with the `step` binary from [github.com/coinbase/step](https://github.com/coinbase/step). 

Bootstrap (directly upload to the Step Function and Lambda):

```bash
# Use AWS credentials or assume-role into AWS
# Build linux zip for lambda
GOOS=linux go build -o lambda
zip lambda.zip lambda

# Tell step to bootstrap this lambda
step bootstrap                        \
  -lambda "coinbase-step-hello-world" \
  -step "coinbase-step-hello-world"   \
  -states "$(./step-hello-world json)"
```

Deploy (via the step-deployer step function in github.com/coinbase/step):

```bash
GOOS=linux go build -o lambda
zip lambda.zip lambda

# Tell step-deployer to deploy this lambda
step deploy                           \
  -lambda "coinbase-step-hello-world" \
  -step "coinbase-step-hello-world"   \
  -states "$(./step-hello-world json)"
```

To invoke the deployed step-hello-world in AWS requires the `aws-cli`:

```bash

ARN="arn:aws:states:${AWS_REGION}:${AWS_ACCOUNT_ID}:stateMachine:coinbase-step-hello-world"

EXECUTION_ARN=$(aws stepfunctions start-execution \
                     --state-machine-arn $ARN \
                     --input '{"Greeting": "Hi"}' | jq -r ".executionArn" )

sleep 1

aws stepfunctions describe-execution --execution-arn $EXECUTION_ARN
{
    "executionArn": "...coinbase-step-hello-world:execution",
    "stateMachineArn": "...stateMachine:coinbase-step-hello-world",
    "name": "...",
    "status": "SUCCEEDED",
    "startDate": 1520050165.321,
    "stopDate": 1520050165.597,
    "input": "{\"Greeting\": \"Hi\"}",
    "output": "{\"Greeting\":\"Hi\"}"
}

aws stepfunctions get-execution-history --execution-arn $EXECUTION_ARN

{
    "events": [
        {
            "timestamp": 1520050165.321,
            "type": "ExecutionStarted",
            "id": 1,
            "previousEventId": 0,
            "executionStartedEventDetails": {
                "input": "{\"Greeting\": \"Hi\"}",
                "roleArn": "...role/coinbase-step-hello-world-step-function-role"
            }
        },
        ...
}
```

