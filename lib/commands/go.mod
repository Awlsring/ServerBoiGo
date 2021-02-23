module ServerBoi/commands

go 1.14

require (
	ServerBoi/cfg v0.0.0
	ServerBoi/services v0.0.0
	github.com/aws/aws-sdk-go-v2/service/ssm v1.1.1 // indirect
	github.com/bwmarrin/discordgo v0.23.2
)

replace ServerBoi/cfg => ../cfg

replace ServerBoi/services => ../services
