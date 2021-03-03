package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"ServerBoi/cfg"
	"ServerBoi/services"

	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/go-a2s"
)

func Server(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string) {

	if len(messageSlice) >= 3 {

		subcommand := strings.ToLower(messageSlice[2])

		serverFunctions := map[string]func(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string){
			"start":     Start,
			"stop":      Stop,
			"reboot":    Reboot,
			"info":      Info,
			"authorize": authorizeOnServer,
			"stats":     stats,
			"backup":    backup,
			"players":   currentPlayerCount,
		}

		if serverFunc, ok := serverFunctions[subcommand]; ok {
			serverFunc(s, m, servers, messageSlice)
		} else {
			msg := fmt.Sprintf("`%v` is not a valid option for !server", subcommand)

			s.ChannelMessageSend(m.ChannelID, msg)
		}
	} else {
		msg := fmt.Sprintf("`%v` is not a valid `!server` command. `!server` commands should be structured like `!server <server_id> <action>`", m.Content)

		s.ChannelMessageSend(m.ChannelID, msg)
	}

}

func currentPlayerCount(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string) {
	fmt.Println("PlayerCount")

	var msg string

	server, er := getTargetServer(messageSlice[1], servers)
	if er != "" {
		msg := fmt.Sprintf("%v", er)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}

	instanceInfo := services.GetInstanceInfo(server)
	state := instanceInfo["state"]
	ip := instanceInfo["ip"]

	if ip != "" && state == "running" {

		port := server.ServerInfo.Port

		serverInfo, serverErr := SteamA2SServerInfo(ip, port)
		if serverErr != nil {
			log.Println(serverErr)
			msg = "Unable to recieve server metadata."
		} else {
			pcString := strconv.Itoa(int(serverInfo.Players))
			msg = fmt.Sprintf("Current player count is %v.", pcString)
		}
	} else {
		msg = fmt.Sprintf("The server must be running to get player count. The server is currently %v.", state)
	}

	s.ChannelMessageSend(m.ChannelID, msg)

}

// SteamA2SServerInfo | Query steam server via A2S protocol to recieve server metadata.
func SteamA2SServerInfo(ip string, port string) (*a2s.ServerInfo, error) {
	clientString := fmt.Sprintf("%v:%v", ip, port)

	client, err := a2s.NewClient(clientString)
	if err != nil {
		fmt.Println(err)
	}

	defer client.Close()

	info, queryErr := client.QueryInfo()
	if queryErr != nil {
		fmt.Println(queryErr)
	}

	client.Close()

	return info, queryErr
}

func backup(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string) {
	fmt.Printf("Backup")

	var msg string

	server, er := getTargetServer(messageSlice[1], servers)
	if er != "" {
		msg := fmt.Sprintf("%v", er)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}

	instanceInfo := services.GetInstanceInfo(server)
	state := instanceInfo["state"]

	if state == "running" {

		premsg := "Attempting to back up save..."
		s.ChannelMessageSend(m.ChannelID, premsg)

		msg = services.RunServerBackup(server)

	} else {
		msg = fmt.Sprintf("The server must be running to get save. The server is currently %v.", state)
	}

	s.ChannelMessageSend(m.ChannelID, msg)

}

func stats(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string) {
	fmt.Println("CPU")

	var msg string

	server, er := getTargetServer(messageSlice[1], servers)
	if er != "" {
		msg := fmt.Sprintf("%v", er)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}

	instanceInfo := services.GetInstanceInfo(server)
	state := instanceInfo["state"]

	if state == "running" {
		premsg := "Getting server stats..."
		s.ChannelMessageSend(m.ChannelID, premsg)

		msg = services.GetServerCPU(server)

	} else {
		msg = fmt.Sprintf("The server must be running to get stats. The server is currently %v.", state)
	}

	s.ChannelMessageSend(m.ChannelID, msg)

}

func authorizeOnServer(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string) {
	fmt.Println("Authorize")

	server, er := getTargetServer(messageSlice[1], servers)
	if er != "" {
		msg := fmt.Sprintf("%v", er)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}

	fmt.Println("getting user from ID")

	newUser, err := s.User("96089969990336512")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(servers)

	fmt.Println(server.Authorized)

	server.Authorized[newUser.ID] = true

}

//Start | Calls a services.StartServer to run the correct API call to start the target server.
func Start(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string) {
	fmt.Println("Start")

	var msg string

	server, err := getTargetServer(messageSlice[1], servers)
	if err != "" {
		msg := fmt.Sprintf("%v", err)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}

	//Check if user is authorized to interact with
	if server.Authorized[m.Author.ID] {

		msg = fmt.Sprintf("Starting server %v. Waiting for IP to be assigned...", server.ID)

		s.ChannelMessageSend(m.ChannelID, msg)

		success := services.StartServer(server)

		if success {
			//Wait for server to be assigned ip
			ip := ""
			for ip == "" {
				instanceInfo := services.GetInstanceInfo(server)
				ip = instanceInfo["ip"]
				time.Sleep(1 * time.Second)
			}

			if server.ServerInfo.Port != "" {
				port := server.ServerInfo.Port
				ip = fmt.Sprintf("%v:%v", ip, port)
			}

			msg = fmt.Sprintf("Server has been started on %v", ip)
		} else {
			msg = "Server is currently in a state where it can be started. Try again in a few minutes."
		}

	} else {
		msg = "Only admin or the server owner may peform this action"

	}
	s.ChannelMessageSend(m.ChannelID, msg)
}

func Stop(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string) {
	fmt.Println("Stop")

	server, err := getTargetServer(messageSlice[1], servers)
	if err != "" {
		msg := fmt.Sprintf("%v", err)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}

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
	server, err := getTargetServer(messageSlice[1], servers)
	if err != "" {
		msg := fmt.Sprintf("%v", err)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}

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

		if server.ServerInfo.Port != "" {
			port := server.ServerInfo.Port
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

	targetServer, err := getTargetServer(messageSlice[1], servers)
	if err != "" {
		msg := fmt.Sprintf("%v", err)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}

	instanceInfo := services.GetInstanceInfo(targetServer)

	ip := instanceInfo["ip"]
	state := instanceInfo["state"]

	msg := fmt.Sprintf("The server is currently %v", state)

	if ip != "" {
		msg = fmt.Sprintf("%v and its IP is %v", msg, ip)
		if targetServer.ServerInfo.Port != "" {
			port := targetServer.ServerInfo.Port
			msg = fmt.Sprintf("%v:%v.", msg, port)
		}
	} else {

		msg = fmt.Sprintf("%v.", msg)
	}

	s.ChannelMessageSend(m.ChannelID, msg)

}
