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
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
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

	pipeline := pipelines.NewCodePipeline(stack, jsii.String("Pipeline"), &pipelines.CodePipelineProps{
		SelfMutation:  jsii.Bool(false),
		UseChangeSets: jsii.Bool(false),
		PipelineName:  jsii.String("GoCdkTestPipeline"),
		Synth: pipelines.NewShellStep(jsii.String("Synth"), &pipelines.ShellStepProps{
			Input: pipelines.CodePipelineSource_Connection(jsii.String("AdrianMendez1199/go-cdk-test"), jsii.String("main"), &pipelines.ConnectionSourceOptions{
				ConnectionArn:        jsii.String("arn:aws:codeconnections:us-east-1:009160027850:connection/eff0aae5-4bbc-4626-b4de-ddd97ae50dcb"), // Created using the AWS console // Check best practices for creating CodeStar Connections
				ActionName:           jsii.String("Build"),
				CodeBuildCloneOutput: jsii.Bool(true),
				TriggerOnPush:        jsii.Bool(true),
			}),

			Commands: &[]*string{
				jsii.String("npm install -g aws-cdk"),
				jsii.String("go mod tidy"),
				jsii.String("cdk synth"),
			},
			PrimaryOutputDirectory: jsii.String("cdk.out"),
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

// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://aws.github.io/aws-sdk-go-v2/docs/getting-started/
