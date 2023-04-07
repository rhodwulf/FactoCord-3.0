package utils

import (
	"fmt"
	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

var VersionDoc = support.CommandDoc{
	Name: "version",
	Doc: `el comando emite la versión del servidor de factorio y la versión de FactoCord.
Si dice que la versión de FactoCord es desconocida, busque en el error.log`,
}

func VersionString(s *discordgo.Session, _ string) {
	factorioVersion, err := support.FactorioVersion()
	if err != nil {
		support.Send(s, "Lo sentimos, hubo un error al verificar la versión de factorio")
		support.Panik(err, "... al correr `factorio --version`")
		return
	}
	res := "Versión del servidor: **" + factorioVersion + "**"

	res += fmt.Sprintf("\nVersión de FactoCord: **%s**", support.FactoCordVersion)

	support.Send(s, res)
}
