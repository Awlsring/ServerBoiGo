package commands

import (
	"github.com/bwmarrin/discordgo"
)

type Conversation struct {
	UserID      string
	CommandTree CommandTree
}

type CommandTree struct {
	Name         string
	Id           string
	Stages       map[int]func(s *discordgo.Session, m *discordgo.MessageCreate, c *CommandTree)
	CurrentStage int  `default:"0"`
	Locked       bool `default:"false"`
	Correction   bool `default:"false"`
	CommandCache map[string]string
}

func (c CommandTree) NextStage(s *discordgo.Session, m *discordgo.MessageCreate, co *CommandTree) {
	co.CurrentStage = c.CurrentStage + 1
	c.Stages[c.CurrentStage](s, m, &c)
}
