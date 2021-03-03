package commands

import (
	"fmt"
	"log"
	"strconv"

	"ServerBoi/cfg"
	"ServerBoi/services"

	"github.com/bwmarrin/discordgo"
)

func getTargetServer(target string, servers map[int]cfg.Server) (cfg.Server, string) {
	serverID, err := strconv.Atoi(target)
	if err != nil {
		log.Fatalf("Cant convert given id to int")
	}

	if server, ok := servers[serverID]; ok {
		return server, ""
	} else {
		return cfg.Server{}, fmt.Sprintf("Server ID `%v` doesn't exist in tracked servers.", serverID)
	}

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

			fmt.Println(instanceInfo)

			ip := instanceInfo["ip"]
			state := instanceInfo["state"]

			comsg := fmt.Sprintf(" **-** ID: %v | Name: %v | Game: %v", server.ID, server.Name, server.Game)

			if ip != "" {
				comsg = fmt.Sprintf("%v | IP: %v", comsg, ip)
				if server.ServerInfo.Port != "" {
					port := server.ServerInfo.Port
					comsg = fmt.Sprintf("%v:%v", comsg, port)
				}
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
	start := "`!server <server_id> start`"
	stop := "`!server <server_id> stop`"
	reboot := "`!server <server_id> reboot`"
	info := "`!server <server_id> info`"
	stats := "`!server <server_id> stats`"
	save := "`!server <server_id> backup`"
	players := "`!server <server_id> players`"

	msg := fmt.Sprintf(`
Here are my current commands:

**General Commands**
%v | lists all currently managed servers.

**Server Commands**
%v | Starts target server. Admin or owner only.
%v | Stops target server. Admin or owner only.
%v | Reboots target server. Admin or owner only.
%v | Returns servers info.
%v | Returns CPU, Mem, and Disk stats for instance.
%v | Runs a back up on the game world.
%v | Returns current server player count.

View my code at https://github.com/Awlsring/ServerBoiGo
	`, list, start, stop, reboot, info, stats, save, players)

	s.ChannelMessageSend(m.ChannelID, msg)
}
