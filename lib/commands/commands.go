package commands

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"ServerBoi/cfg"
	"ServerBoi/services"

	"github.com/bwmarrin/discordgo"
)

type Conversation struct {
	UserID      string
	CommandTree CommandTree
}

type CommandTree struct {
	Stages       map[string]func()
	CurrentStage string `default:"0"`
	Locked       bool   `default:"false"`
	Correction   bool   `default:"false"`
}

func getTargetServer(messageSlice []string, servers map[int]cfg.Server) cfg.Server {
	serverID, err := strconv.Atoi(messageSlice[2])
	if err != nil {
		log.Fatalf("Cant convert given id to int")
	}

	targetServer := servers[serverID]

	return targetServer
}

func AddServer(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string) {
	fmt.Println("Add Server")

	stages := map[string]func(){
		"0": AddServerStageNameSet,
	}

	convo := Conversation{
		UserID: m.Author.ID,
		CommandTree: CommandTree{
			Stages: stages,
		},
	}

	//Start First Stage
	convo.CommandTree.Stages["0"]()

}

func AddServerStageNameSet() {
	fmt.Println("Add Server")
}

func Fun(s *discordgo.Session, m *discordgo.MessageCreate) {
	nou := []string{"no u", "nou", "n0 u", "no you", "noyou", "n o u", "n 0 u", "no, u"}
	thx := []string{"thanks", "thx", "thank"}
	sorry := []string{"i'm sorry", "sorry", "my bad", "sorry"}

	arrays := [...][]string{nou, thx, sorry}
	arrayNames := [...]string{"nou", "thx", "sorry"}
	var contained bool
	var containedIn string

	fun := map[string]func(s *discordgo.Session, m *discordgo.MessageCreate){
		"nou":   nouFunc,
		"thx":   thxFunc,
		"sorry": sorryFunc,
	}

	for i, array := range arrays {
		for _, opt := range array {
			contained = strings.HasPrefix(m.Content, opt)
			if contained {
				containedIn = arrayNames[i]
				fun[containedIn](s, m)
				break
			}
		}
	}
}

func nouFunc(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("NOU")
	randNum := rand.Intn(50)

	var msg string

	switch randNum {
	case 1:
		msg = "Fuck you"
	case 2, 8, 9, 10:
		msg = fmt.Sprintf("No u, %v", m.Author.Username)
	case 3, 11, 12, 13, 14:
		msg = "No u buddy"
	case 4:
		msg = "Nerd"
	case 5, 15, 16, 17, 18:
		msg = fmt.Sprintf("No, fuk u %v", m.Author.Username)
	case 6:
		msg = "Wow, you're right. I've never though about it that way before"
	case 7:
		msg = "Yeah u rite"
	default:
		msg = "No u"
	}

	s.ChannelMessageSend(m.ChannelID, msg)

}

func thxFunc(s *discordgo.Session, m *discordgo.MessageCreate) {
	np := [...]string{"Np", "No prob", "Gotchu"}

	num := rand.Intn(len(np))

	msg := np[num]

	s.ChannelMessageSend(m.ChannelID, msg)

}

func sorryFunc(s *discordgo.Session, m *discordgo.MessageCreate) {
	app := [...]string{
		"Don't appologize it shows weakness",
		"Its okay just don't let it happen again",
		"Good",
	}

	num := rand.Intn(len(app))

	msg := app[num]

	s.ChannelMessageSend(m.ChannelID, msg)

}

func Start(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string) {
	fmt.Println("Start")

	server := getTargetServer(messageSlice, servers)

	//Check if user is authorized to interact with
	if server.Authorized[m.Author.ID] {

		msg := fmt.Sprintf("Starting server %v. Waiting for IP to be assigned...", server.ID)

		s.ChannelMessageSend(m.ChannelID, msg)

		services.StartServer(server)

		//Wait for server to be assigned ip
		ip := ""
		for ip == "" {
			instanceInfo := services.GetInstanceInfo(server)
			ip = instanceInfo["ip"]
			time.Sleep(1 * time.Second)
		}

		if port, ok := server.ServerInfo["Port"]; ok {
			ip = fmt.Sprintf("%v:%v", ip, port)
		}

		msg = fmt.Sprintf("Server has been started on %v", ip)

		s.ChannelMessageSend(m.ChannelID, msg)

	} else {
		msg := "Only admin or the server owner may peform this action"

		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

func Stop(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string) {
	fmt.Println("Stop")

	server := getTargetServer(messageSlice, servers)

	//Check if user is authorized to interact with
	if server.Authorized[m.Author.ID] {

		services.StopServer(server)

		msg := fmt.Sprintf("Stopping server %v", server.ID)

		s.ChannelMessageSend(m.ChannelID, msg)
	} else {
		msg := "Only admin or the server owner may peform this action"

		s.ChannelMessageSend(m.ChannelID, msg)
	}

}

func Reboot(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string) {
	fmt.Println("Reboot")
	server := getTargetServer(messageSlice, servers)

	//Check if user is authorized to interact with
	if server.Authorized[m.Author.ID] {

		services.RebootServer(server)

		msg := fmt.Sprintf("Rebooting server %v... waiting for new IP te be assigned.", server.ID)

		s.ChannelMessageSend(m.ChannelID, msg)

		//Wait for server to be assigned ip
		ip := ""
		state := ""
		for ip == "" && state == "running" {
			instanceInfo := services.GetInstanceInfo(server)
			ip = instanceInfo["ip"]
			state = instanceInfo["state"]
			time.Sleep(1 * time.Second)
		}

		if port, ok := server.ServerInfo["Port"]; ok {
			ip = fmt.Sprintf("%v:%v", ip, port)
		}

		msg = fmt.Sprintf("Server has been started on %v", ip)

		s.ChannelMessageSend(m.ChannelID, msg)

	} else {
		msg := "Only admin or the server owner may peform this action"

		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

func Info(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string) {

	targetServer := getTargetServer(messageSlice, servers)

	instanceInfo := services.GetInstanceInfo(targetServer)

	ip := instanceInfo["ip"]
	state := instanceInfo["state"]

	msg := fmt.Sprintf("The server is currently %v", state)

	if ip != "" {
		msg = fmt.Sprintf("%v and its IP is %v", msg, ip)
	}

	if port, ok := targetServer.ServerInfo["Port"]; ok {
		msg = fmt.Sprintf("%v:%v.", msg, port)
	} else {
		msg = fmt.Sprintf("%v.", msg)
	}

	s.ChannelMessageSend(m.ChannelID, msg)

}

func List(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string) {
	serverAmount := len(servers)

	fmt.Println("List")

	messageChannel := make(chan string)

	msg := "**Current managed servers:**\n"

	for _, server := range servers {
		go func(server cfg.Server, messageChannel chan string) {

			fmt.Println(server)

			instanceInfo := services.GetInstanceInfo(server)
			ip := instanceInfo["ip"]
			state := instanceInfo["state"]

			comsg := fmt.Sprintf(" **-** ID: %v | Name: %v | Game: %v", server.ID, server.Name, server.Game)

			if ip != "" {
				comsg = fmt.Sprintf("%v | IP: %v", comsg, ip)
			}

			if port, ok := server.ServerInfo["Port"]; ok {
				comsg = fmt.Sprintf("%v:%v", comsg, port)
			}

			comsg = fmt.Sprintf("%v | Status: %v\n", comsg, state)

			messageChannel <- comsg
		}(server, messageChannel)

	}

	for i := 0; i < serverAmount; i++ {
		msg = fmt.Sprintf("%v%v", msg, <-messageChannel)
	}

	fmt.Printf("Sending message")

	s.ChannelMessageSend(m.ChannelID, msg)

}

func Help(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string) {
	list := "`!list`"
	start := "`!start server <server_id>`"
	stop := "`!stop server <server_id>`"
	reboot := "`!reboot server <server_id>`"
	info := "`!info server <server_id>`"

	msg := fmt.Sprintf(`
Here are my current commands:

**General Commands**
%v | lists all currently managed servers.

**Server Commands**
%v | Starts target server. Admin or owner only.
%v | Stops target server. Admin or owner only.
%v | Reboots target server. Admin or owner only.
%v | Returns servers info.

View my code at https://github.com/Awlsring/ServerBoiGo
	`, list, start, stop, reboot, info)

	s.ChannelMessageSend(m.ChannelID, msg)
}
