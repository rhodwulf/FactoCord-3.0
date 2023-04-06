package admin

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

// KickPlayerUsage comment...

var KickPlayerDoc = support.CommandDoc{
	Name:  "kick",
	Usage: "$kick <player> <reason>",
	Doc:   `el comando expulsa al jugador del servidor por una razón específica`,
}

// KickPlayer kicks a player from the server.
func KickPlayer(s *discordgo.Session, args string) {
	if len(args) == 0 {
		support.SendFormat(s, "Uso: "+KickPlayerDoc.Usage)
		return
	}
	args2 := strings.SplitN(args+" ", " ", 2)
	player := strings.TrimSpace(args2[0])
	reason := strings.TrimSpace(args2[1])

	if len(player) == 0 || len(reason) == 0 {
		support.SendFormat(s, "Uso: "+KickPlayerDoc.Usage)
		return
	}
	command := "/kick " + player + " " + reason
	success := support.Factorio.Send(command)
	if success {
		support.Send(s, "Jugador "+player+" ha sido expulsado por la razón "+reason+"!")
	} else {
		support.Send(s, "Lo sentimos, hubo un error al enviar el comando /kick")
	}
}
