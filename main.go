package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"ServerBoi/lib"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Authorized channels
var channels = map[string]bool{
	"242453642362355712": true,
	"713865223584481301": true,
}

var commandMap = map[string]func(s *discordgo.Session, m *discordgo.MessageCreate){
	"no u": test,
}

func test(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("Message from")
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

	// Authorized channels
	channels := map[string]bool{
		"242453642362355712": true,
		"713865223584481301": true,
	}

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if channels[m.ChannelID] {

		fmt.Println("Message from", m.Author.ID, "in", m.ChannelID, "| Message:", m.Content)

		// The classic
		if m.Content == "no u" {
			s.ChannelMessageSend(m.ChannelID, "no u")
		}

	}

}

func main() {

	lib.Butts()

	token := goDotEnvVariable("DISCORD_TOKEN")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
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
