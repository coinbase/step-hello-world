# GeoEngineer Resources For Step Hello World
# geo apply resources/step-hello-world.rb

########################################
###           ENVIRONMENT            ###
########################################

env = environment('step') {
  region      ENV.fetch('AWS_REGION')
  account_id  ENV.fetch('AWS_ACCOUNT_ID')
}

########################################
###            PROJECT               ###
########################################

project = project('coinbase', 'step-hello-world') {
  environments 'step'
  tags {
    ProjectName "coinbase/step-hello-world"
    ConfigName "development"
    DeployWith "step-deployer"
    self[:org] = "coinbase"
    self[:project] = "step-hello-world"
  }
}

step_name = "#{project.org}-#{project.name}"
step_role_name =  "#{step_name}-step-function-role"
lambda_role_name = "#{step_name}-lambda-role"

########################################
###         Step Function            ###
########################################

step_role = project.resource('aws_iam_role', step_role_name) {
  name step_role_name
  path "/step/#{project.full_name}/development/"
  assume_role_policy(
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": {
            "Service": "states.us-east-1.amazonaws.com"
          },
          "Action": "sts:AssumeRole"
        }
      ]
    }.to_json
  )
}

project.resource('aws_iam_role_policy', step_role_name) {
  depends_on [step_role.terraform_name]
  name step_role_name
  role step_role_name
  policy(
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "lambda:InvokeFunction"
          ],
          "Resource": "*"
        }
    ]
    }.to_json
  )
}

project.resource("aws_sfn_state_machine","sfn_state_machine") {
  depends_on [step_role.terraform_name]
  name       step_name
  role_arn   step_role.to_ref("arn")
  definition '{
    "StartAt": "Noop",
    "States": {
      "Noop": {
        "Type": "Pass",
        "End": true
      }
    }
  }'

  lifecycle {
    ignore_changes ["definition"] # Ignore changes to state machine
  }
}

########################################
###         IAM  Role                ###
########################################

lambda_role = project.resource('aws_iam_role', lambda_role_name) {
  name lambda_role_name
  path "/"
  assume_role_policy(
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Sid": "",
          "Effect": "Allow",
          "Principal": {
            "Service": "lambda.amazonaws.com"
          },
          "Action": "sts:AssumeRole"
        }
      ]
    }.to_json
  )
}

project.resource('aws_iam_role_policy', lambda_role_name) {
  depends_on [lambda_role.terraform_name]
  name lambda_role_name
  role lambda_role_name
  policy(
    {
      "Version": "2012-10-17",
      "Statement": [
        # WRITE TO LOGS
        {
          "Effect": "Allow",
          "Action": [
            "logs:CreateLogGroup",
            "logs:CreateLogStream",
            "logs:PutLogEvents"
          ],
          "Resource": "arn:aws:logs:*:*:log-group:/aws/lambda/*"
        }
      ]
    }.to_json
  )
}

########################################
###            Lambda                ###
########################################

lambda_function = project.resource("aws_lambda_function", step_name) {
  function_name step_name
  description step_name

  role lambda_role.to_ref('arn')

  lifecycle {
    ignore_changes ["environment", "filename", "source_code_hash"]
  }

  filename File.expand_path(File.dirname(__FILE__)) + '/lambda.zip'
  handler "lambda"
  memory_size 128
  runtime "go1.x"
  timeout "30"
  publish "true"

}
