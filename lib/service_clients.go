package lib

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// Butts FUCK YOU
func Butts() {
	fmt.Println("Toes")

	client := CreateEc2Client("us-west-2")

	fmt.Println(client)
}

// CreateEc2Client creates a client to communicate with EC2
func CreateEc2Client(region string) ec2.Client {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	instanceID := "i-044e820ec4e8b0335"

	roleArn := "arn:aws:iam::576187103674:role/ServerBoiRole"
	roleSession := "Toes"

	stsClient := sts.NewFromConfig(cfg)

	input := &sts.AssumeRoleInput{
		RoleArn:         &roleArn,
		RoleSessionName: &roleSession,
	}

	result, err := TakeRole(context.TODO(), stsClient, input)
	if err != nil {
		fmt.Println("Got an error assuming the role:")
		fmt.Println(err)
	}

	accessKey := result.Credentials.AccessKeyId
	secretKey := result.Credentials.SecretAccessKey
	sessionToken := result.Credentials.SessionToken

	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(*accessKey, *secretKey, *sessionToken))

	// value, err := creds.Retrieve(context.TODO())

	ec2cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	client := ec2.NewFromConfig(ec2cfg)

	ec2input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{
			instanceID,
		},
	}

	resp, err := client.DescribeInstances(context.TODO(), ec2input)

	fmt.Println(resp)

	return *client
}

//TakeRole is
func TakeRole(c context.Context, api STSAssumeRoleAPI, input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	return api.AssumeRole(c, input)
}

// Test Buuts.
func Test() {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// instanceID := "i-345345345345345345"

	client := ec2.NewFromConfig(cfg)

	// input := &ec2.StartInstancesInput{
	// 	InstanceIds: []string{
	// 		instanceID,
	// 	},
	// }

	fmt.Println(cfg, client)

	// client.StartInstances(context.TODO(), input)

	// _, err = StartInstance(context.TODO(), client, input)
}

// EC2StartInstancesAPI defines the interface for the StartInstances function.
// We use this interface to test the function using a mocked service.
type EC2StartInstancesAPI interface {
	StartInstances(ctx context.Context,
		params *ec2.StartInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.StartInstancesOutput, error)
}

type STSAssumeRoleAPI interface {
	AssumeRole(ctx context.Context,
		params *sts.AssumeRoleInput,
		optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error)
}

// StartInstance starts an Amazon Elastic Compute Cloud (Amazon EC2) instance.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a StartInstancesOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to StartInstances.
func StartInstance(c context.Context, api EC2StartInstancesAPI, input *ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error) {
	resp, err := api.StartInstances(c, input)

	return resp, err
}
