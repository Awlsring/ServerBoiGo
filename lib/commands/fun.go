package commands

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Fun(s *discordgo.Session, m *discordgo.MessageCreate) {
	message := strings.ToLower(m.Content)

	nou := []string{"no u", "nou", "n0 u", "no you", "noyou", "n0u", "no you", "n0you", "n o u", "n 0 u", "no, u"}
	thx := []string{"thanks", "thx", "thank", "ty"}
	sorry := []string{"i'm sorry", "sorry", "my bad", "sorry"}

	arrays := [...][]string{nou, thx, sorry}
	arrayNames := [...]string{"nou", "thx", "sorry"}
	var contained bool
	var containedIn string

	fun := map[string]func(s *discordgo.Session, m *discordgo.MessageCreate){
		"nou":   nouFunc,
		"thx":   thxFunc,
		"sorry": sorryFunc,
	}

	for i, array := range arrays {
		for _, opt := range array {
			contained = strings.HasPrefix(message, opt)
			if contained {
				containedIn = arrayNames[i]
				fun[containedIn](s, m)
				break
			}
		}
	}
}

func nouFunc(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("NOU")
	randNum := rand.Intn(50)

	var msg string

	switch randNum {
	case 1:
		msg = "Fuck you"
	case 2, 8, 9, 10:
		msg = fmt.Sprintf("No u, %v", m.Author.Username)
	case 3, 11, 12, 13, 14:
		msg = "No u buddy"
	case 4:
		msg = "Nerd"
	case 5, 15, 16, 17, 18:
		msg = fmt.Sprintf("No, fuk u %v", m.Author.Username)
	case 6:
		msg = "Wow, you're right. I've never though about it that way before"
	case 7:
		msg = "Yeah u rite"
	default:
		msg = "No u"
	}

	s.ChannelMessageSend(m.ChannelID, msg)

}

func thxFunc(s *discordgo.Session, m *discordgo.MessageCreate) {
	np := [...]string{"Np", "No prob", "Gotchu"}

	num := rand.Intn(len(np))

	msg := np[num]

	s.ChannelMessageSend(m.ChannelID, msg)

}

func sorryFunc(s *discordgo.Session, m *discordgo.MessageCreate) {
	app := [...]string{
		"Don't appologize it shows weakness",
		"Its okay just don't let it happen again",
		"Good",
	}

	num := rand.Intn(len(app))

	msg := app[num]

	s.ChannelMessageSend(m.ChannelID, msg)

}
