package utils

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

var OnlineDoc = support.CommandDoc{
	Name: "online",
	Doc:  `muestra los jugadores en línea (y el número máximo de jugadores si está configurado)`,
}

func getOnline(info *gameInfo) *support.TextListT {
	if len(info.Players) == 0 {
		return &support.TextListT{
			Heading: "**nadie está en línea**",
			None:    "",
		}
	}
	maxPlayers := ""
	if info.MaxPlayers != 0 {
		maxPlayers = fmt.Sprintf("/%d", info.MaxPlayers)
	}
	online := support.DefaultTextList(
		fmt.Sprintf("**%d%s jugador%es en línea:**", len(info.Players), maxPlayers, support.PluralS(len(info.Players))),
	)
	for _, player := range info.Players {
		online.Append(player)
	}
	return &online
}

func GameOnline(s *discordgo.Session, _ string) {
	if !support.Factorio.IsRunning() {
		support.Send(s, "El servidor no está funcionando.")
		return
	}
	if support.Factorio.GameID == "" {
		support.Send(s, "El servidor no registró un juego en el servidor de factorio")
		return
	}

	resp, err := http.Get("https://multiplayer.factorio.com/get-game-details/" + support.Factorio.GameID)
	if err != nil {
		support.Panik(err, "Connection error to /get-game-details")
		support.Send(s, "Some connection error occurred")
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		support.Panik(err, "Error reading /get-game-details")
		support.Send(s, "Some connection error occurred")
		return
	}

	info := gameInfo{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		support.Panik(err, "Error unmarshalling /get-game-details")
		support.Send(s, "Some json error occurred")
		return
	}
	if info.Message != "" {
		support.Send(s, "The server reports: "+info.Message)
		return
	}

	support.Send(s, getOnline(&info).Render())
}
