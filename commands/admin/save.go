package admin

import (
	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

var SaveServerDoc = support.CommandDoc{
	Name: "save",
	Doc:  `comando envía un comando para guardar el juego en el servidor`,
}

// SaveServer executes the save command on the server.
func SaveServer(s *discordgo.Session, args string) {
	if len(args) != 0 {
		support.Send(s, "Save no acepta argumentos")
		return
	}
	success := support.Factorio.Send("/save")
	if success {
		support.Factorio.SaveRequested = true
		//support.Send(s, "Server saved successfully!")
	} else {
		support.Send(s, "Lo sentimos, hubo un error al enviar /save command")
	}
}
