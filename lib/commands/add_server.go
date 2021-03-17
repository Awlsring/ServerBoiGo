package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// func (c CommandTree)

// func AddServer() *CommandTree {

// 	stages := map[string]func(){
// 		"0": stageAddServerNameGet,
// 	}

// 	return &CommandTree{
// 		Name: "AddServer",
// 		Stages: stages,

// 	}
// }

func AddServer(s *discordgo.Session, m *discordgo.MessageCreate) *Conversation {
	fmt.Println("Add Server")

	//Outline Stages
	var stages = map[int]func(s *discordgo.Session, m *discordgo.MessageCreate, c *CommandTree){
		1: stageAddServerNameGet,
		2: stageAddServerNameSet,
		3: stageAddServerGameGet,
	}

	//Create convo
	convo := Conversation{
		UserID: m.Author.ID,
		CommandTree: CommandTree{
			Name:   "AddServer",
			Stages: stages,
		},
	}

	//Start First Stage
	// convo.CommandTree.Stages["0"](m, &convo)

	return &convo

}

func stageAddServerNameGet(s *discordgo.Session, m *discordgo.MessageCreate, c *CommandTree) {
	msg := "What is the servers name?"

	s.ChannelMessageSend(m.ChannelID, msg)
}

func stageAddServerNameSet(s *discordgo.Session, m *discordgo.MessageCreate, c *CommandTree) {
	c.CommandCache["Name"] = m.Content

	c.NextStage(s, m, c)

}

func stageAddServerGameGet(s *discordgo.Session, m *discordgo.MessageCreate, c *CommandTree) {
	msg := "What game is this for?"

	s.ChannelMessageSend(m.ChannelID, msg)
}
