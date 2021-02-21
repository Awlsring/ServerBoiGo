package commands

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"ServerBoi/cfg"
	"ServerBoi/services"

	"github.com/bwmarrin/discordgo"
)

func getTargetServer(messageSlice []string, servers map[int]cfg.Server) cfg.Server {
	serverID, err := strconv.Atoi(messageSlice[2])
	if err != nil {
		log.Fatalf("Cant convert given id to int")
	}

	targetServer := servers[serverID]

	return targetServer
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
