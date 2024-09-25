package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/pipelines"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GoCdkTestStackProps struct {
	awscdk.StackProps
}

func NewGoCdkTestStack(scope constructs.Construct, id string, props *GoCdkTestStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	lambdaFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("go-cdk-test"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("get-user-lambda"),
		Entry:        jsii.String("./cmd/api"),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Environment: &map[string]*string{
			"PRODUCTION": jsii.String("true"),
		},
	})

	api := awsapigateway.NewRestApi(stack, jsii.String("ApiGateway"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("GoCdkTestApi"),
		DeployOptions: &awsapigateway.StageOptions{
			StageName: jsii.String("dev"),
		},
	})

	v1 := api.Root().AddResource(jsii.String("v1"), nil)
	userResource := v1.AddResource(jsii.String("user"), nil)
	userResource.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(lambdaFunction, nil), nil)

	return stack
}

// Define a new stage that includes the GoCdkTestStack
func NewGoCdkTestStage(scope constructs.Construct, id string, env *awscdk.Environment) awscdk.Stage {
	stage := awscdk.NewStage(scope, &id, &awscdk.StageProps{
		Env: env,
	})

	NewGoCdkTestStack(stage, "GoCdkTestStack", &GoCdkTestStackProps{
		StackProps: awscdk.StackProps{
			Env: env,
		},
	})

	return stage
}

func NewPipelineStack(scope constructs.Construct, id string, env *awscdk.Environment) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &awscdk.StackProps{
		Env: env,
	})

	githubSource := pipelines.CodePipelineSource_GitHub(jsii.String("AdrianMendez1199/go-cdk-test"), jsii.String("main"), &pipelines.GitHubSourceOptions{
		Authentication: awscdk.SecretValue_SecretsManager(jsii.String("githubAccessToken"), nil),
	})

	pipeline := pipelines.NewCodePipeline(stack, jsii.String("Pipeline"), &pipelines.CodePipelineProps{
		PipelineName: jsii.String("GoCdkTestPipeline"),
		Synth: pipelines.NewShellStep(jsii.String("Synth"), &pipelines.ShellStepProps{
			Input: githubSource,
			Commands: &[]*string{
				jsii.String("npm install -g aws-cdk"),
				jsii.String("go mod tidy"),
				jsii.String("cdk synth"),
			},
		}),
	})

	// Add the GoCdkTestStage to the pipeline
	pipeline.AddStage(NewGoCdkTestStage(stack, "DeployStage", env), nil)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewPipelineStack(app, "PipelineStack", env())

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed.
func env() *awscdk.Environment {
	return nil
	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }
}
