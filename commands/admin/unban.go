package admin

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

var UnbanPlayerDoc = support.CommandDoc{
	Name:  "unban",
	Usage: "$unban <player>",
	Doc:   `El comando elimina al jugador de la lista de prohibición en el servidor.`,
}

// UnbanPlayer unbans a player on the server.
func UnbanPlayer(s *discordgo.Session, args string) {
	if strings.ContainsAny(args, " \n\t") {
		support.SendFormat(s, "Uso: "+UnbanPlayerDoc.Usage)
		return
	}
	command := "/unban " + args
	success := support.Factorio.Send(command)
	if success {
		support.Send(s, "Jugador "+args+" ha sido desbaneado!")
	} else {
		support.Send(s, "Lo sentimos, hubo un error al enviar el comando /unban")
	}
}
