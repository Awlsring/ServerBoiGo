module ServerBoi/services

go 1.14

require (
	ServerBoi/cfg v0.0.0
	github.com/aws/aws-sdk-go-v2 v1.2.0
	github.com/aws/aws-sdk-go-v2/config v1.1.1
	github.com/aws/aws-sdk-go-v2/credentials v1.1.1
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.1.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.1.1
	github.com/aws/aws-sdk-go-v2/service/sts v1.1.1
)

replace ServerBoi/cfg => ../cfg
