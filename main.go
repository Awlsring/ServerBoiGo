package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"ServerBoi/cfg"
	"ServerBoi/commands"
	"ServerBoi/services"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Globals Vars
var servers = cfg.LoadConfig()

var conversations = make(map[string]commands.Conversation)

// Disallowed channels. (TODO Load from config)
var channels = map[string]bool{
	"170364850256609280": false,
	"666432608539901963": false,
	"453802459379269633": false,
	"278255133891100673": false,
	"316003727506931713": false,
	"585951696753131520": false,
	"616679427979476994": false,
	"711488008351645758": false,
	"186263688603369473": false,
	"698658837447704707": false,
	"689705433174114353": false,
}

var commandMap = map[string]func(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string){
	"!server": commands.Server,
	"!list":   commands.List,
	"!ls":     commands.List,
	"!help":   commands.Help,
}

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	convoKey := fmt.Sprintf("%v-%v", m.Author.ID, m.ChannelID)

	if convo, ok := conversations[convoKey]; ok {
		// Resume convo logic
		fmt.Println("Resuming conversation")
		currentStage := convo.CommandTree.CurrentStage
		convo.CommandTree.Stages[currentStage]()
	}

	messageSlice := strings.Split(m.Content, " ")

	// If valid channel and has ! at the start.
	if !channels[m.ChannelID] && strings.HasPrefix(m.Content, "!") {

		fmt.Println("Message from", m.Author.ID, "in", m.ChannelID, "| Message:", m.Content)

		command := strings.ToLower(messageSlice[0])

		if command, ok := commandMap[command]; ok {
			go command(s, m, servers, messageSlice)
		} else {
			msg := fmt.Sprintf("`%v` is not a command.", m.Content)

			s.ChannelMessageSend(m.ChannelID, msg)
		}
	} else {
		fmt.Println("Else")
		commands.Fun(s, m)
	}

}

func main() {

	token := goDotEnvVariable("DISCORD_TOKEN")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Set Discord status for a little flex.
	dg.Identify.Presence = discordgo.GatewayStatusUpdate{
		Game: discordgo.Activity{
			Name: "you | Use !help to start",
			Type: discordgo.ActivityTypeListening,
			URL:  "https://github.com/Awlsring/ServerBoiGo",
		},
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Recieve Meassages from authorized channels and DMs.
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Start coprocess to continuely check for running servers.
	go checkServerActivity(servers)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

// Move this elsewhere
func checkServerActivity(serverList map[int]cfg.Server) {
	// Run every 15 minutes
	// For each server, check if running
	// If running with no active players, mark with a counter
	// If running with counter and no active players, shut down

	serverCounter := map[int]bool{}

	// Infinite for loop that'll run every 15 minutes
	for {
		log.Println("Starting server activity check.")
		// Loop through each server in servers
		for _, server := range serverList {
			log.Printf("Checking server %v", server.Name)

			// Check if server has auto shutdown available
			if server.ServerInfo.AutoShutdown {
				//Grab the server info
				info := services.GetInstanceInfo(server)
				//If the server is on...
				if info["state"] == "running" {
					log.Println("Checking if server is active.")

					ip := info["ip"]
					port := server.ServerInfo.Port
					fmt.Printf("%v:%v", ip, port)

					// ...Get player count
					resp, serverErr := commands.SteamA2SServerInfo(ip, port)
					// If error retrieving data...
					if serverErr != nil {
						logstr := fmt.Sprintf("Error getting player count. Error: %v", serverErr)
						log.Println(logstr)
					} else {
						pc := int(resp.Players)
						// If the player count is 0...
						if pc == 0 {
							if _, exists := serverCounter[server.ID]; exists {
								log.Println("Server has had no players twice in 15 minutes, shutting down.")
								// Save server. (TODO: Check to see if backup is enabled)
								services.RunServerBackup(server)
								// Stop the server
								services.StopServer(server)

								// and delete key from counter map
								delete(serverCounter, server.ID)

							} else {
								log.Println("Server has no players. Server will shutdown in 15 minutes if no players are active next check.")
								// Add server to counter map
								serverCounter[server.ID] = true
							}
						} else {
							// If player count isnt 0, remove tracking if in map
							if _, exists := serverCounter[server.ID]; exists {
								// Delete key from counter map
								delete(serverCounter, server.ID)
							}
						}
					}
				} else {
					log.Println("Server isn't running.")
				}
			} else {
				log.Println("Server doesn't have auto shutdown enabled.")
			}
		}

		//Check status again in 15 minutes
		time.Sleep(15 * time.Minute)
	}
}
