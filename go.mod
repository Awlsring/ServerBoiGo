module ServerBoi

go 1.14

require (
	ServerBoi/cfg v0.0.0
	ServerBoi/commands v0.0.0
	ServerBoi/services v0.0.0
	github.com/aws/aws-sdk-go-v2/service/ssm v1.1.1 // indirect
	github.com/bwmarrin/discordgo v0.23.2
	github.com/joho/godotenv v1.3.0
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f // indirect
)

replace ServerBoi/services => ./lib/services

replace ServerBoi/cfg => ./lib/cfg

replace ServerBoi/commands => ./lib/commands
