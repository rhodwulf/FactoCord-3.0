package admin

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

var BanPlayerDoc = support.CommandDoc{
	Name:  "ban",
	Usage: "$ban <player> <reason>",
	Doc:   `el comando prohíbe al jugador en el servidor por un motivo específico`,
}

// BanPlayer bans a player on the server.
func BanPlayer(s *discordgo.Session, args string) {
	if len(args) == 0 {
		support.SendFormat(s, "Usage: "+BanPlayerDoc.Usage)
		return
	}
	args2 := strings.SplitN(args+" ", " ", 2)
	player := strings.TrimSpace(args2[0])
	reason := strings.TrimSpace(args2[1])

	if len(player) == 0 || len(reason) == 0 {
		support.SendFormat(s, "Usage: "+BanPlayerDoc.Usage)
		return
	}

	command := "/ban " + player + " " + reason
	success := support.Factorio.Send(command)
	if success {
		support.Send(s, "Jugador "+player+" prohibido con razon \""+reason+"\"!")
	} else {
		support.Send(s, "Lo siento, hubo un error al enviar el comando /ban")
	}
}
