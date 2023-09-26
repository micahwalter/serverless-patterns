package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type FirehoseTransformationCdkGoStackProps struct {
	awscdk.StackProps
}

func NewFirehoseTransformationCdkGoStack(scope constructs.Construct, id string, props *FirehoseTransformationCdkGoStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// lambda function for data transforation processing
	// lambda function needs to be compiled to the bin/data-transformation/bootsrtap before deploying with the CDK
	awslambda.NewFunction(stack, jsii.String("DataTransformationFunction"), &awslambda.FunctionProps{
		Code:         awslambda.Code_FromAsset(jsii.String("bin/data-transformation"), nil), //folder where bootstrap executable is located
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"), // Handler named bootstrap
		Architecture: awslambda.Architecture_ARM_64(),
	})

	// s3 bucket for data destination from firehose
	awss3.NewBucket(stack, jsii.String("DataDestinationBucket"), &awss3.BucketProps{
		Encryption: awss3.BucketEncryption_S3_MANAGED,
	})

	// kinesis firehose delivery stream
	// awskinesisfirehose.CfnDeliveryStream(stack, jsii.String("KinesisFirehoseDataStream"), &awskinesisfirehose.CfnDeliveryStreamProps{
	// 	ExtendedS3DestinationConfiguration: awskinesisfirehose.CfnDeliveryStream.,
	// })

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewFirehoseTransformationCdkGoStack(app, "FirehoseTransformationCdkGoStack", &FirehoseTransformationCdkGoStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
