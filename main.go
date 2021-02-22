package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"ServerBoi/cfg"
	"ServerBoi/commands"

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
}

var commandMap = map[string]func(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string){
	"!start":  commands.Start,
	"!stop":   commands.Stop,
	"!reboot": commands.Reboot,
	"!info":   commands.Info,
	"!server": commands.Server,
	"!list":   commands.List,
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

	// convoKey := fmt.Sprintf("%v-%v", m.Author.ID, m.ChannelID)

	// if convo, ok := conversations[convoKey]; ok {
	// 	// Resume convo logic
	// 	fmt.Println("Resuming conversation")
	// 	convo.CommandTree.Command()
	// }

	messageSlice := strings.Split(m.Content, " ")

	// If valid channel and has ! at the start.
	if !channels[m.ChannelID] && strings.HasPrefix(m.Content, "!") {

		fmt.Println("Message from", m.Author.ID, "in", m.ChannelID, "| Message:", m.Content)

		command := messageSlice[0]

		if command, ok := commandMap[command]; ok {
			go command(s, m, servers, messageSlice)
		}
	} else {
		fmt.Println("Else")
		commands.Fun(s, m)
	}

}

/// Maybe take all config stuff to its own file?

func main() {

	token := goDotEnvVariable("DISCORD_TOKEN")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.Identify.Presence = discordgo.GatewayStatusUpdate{
		Game: discordgo.Activity{
			Name: "you | Use !help to start",
			Type: discordgo.ActivityTypeListening,
			URL:  "https://github.com/Awlsring/ServerBoiGo",
		},
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Recieve Meassages from authorized channels ans DMs
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}
