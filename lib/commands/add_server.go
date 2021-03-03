package commands

import (
	"ServerBoi/cfg"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func AddServer(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string) {
	fmt.Println("Add Server")

	//Outline Stages
	var stages = map[string]func(m *discordgo.MessageCreate, c *Conversation) string{
		"0": stageAddServerNameGet,
	}

	//Create convo
	convo := Conversation{
		UserID: m.Author.ID,
		CommandTree: CommandTree{
			Name:     "AddServer",
			Function: AddServer,
			Stages:   stages,
		},
	}

	//Start First Stage
	convo.CommandTree.Stages["0"](m, &convo)

	// func moveStage()

}

func stageAddServerNameGet(m *discordgo.MessageCreate, c *Conversation) string {
	msg := "What is the servers name?"

	return msg
}

func stageAddServerNameSet(m *discordgo.MessageCreate) string {
	msg := "What is the servers name?"

	return msg
}
