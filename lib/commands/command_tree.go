package commands

import (
	"ServerBoi/cfg"

	"github.com/bwmarrin/discordgo"
)

type Conversation struct {
	UserID      string
	CommandTree CommandTree
}

type CommandTree struct {
	Name         string
	Function     func(s *discordgo.Session, m *discordgo.MessageCreate, servers map[int]cfg.Server, messageSlice []string)
	Stages       map[string]func(m *discordgo.MessageCreate, c *Conversation) string
	CurrentStage string `default:"0"`
	Locked       bool   `default:"false"`
	Correction   bool   `default:"false"`
}
