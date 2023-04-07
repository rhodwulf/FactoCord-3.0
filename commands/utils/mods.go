package utils

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"path"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

// ModJson is struct containing a slice of Mod.
type ModJson struct {
	Mods []Mod
}

// Mod is a struct containing info about a mod.
type Mod struct {
	Name    string
	Enabled bool
}

var ModListDoc = support.CommandDoc{
	Name:  "mods",
	Doc:   `el comando genera información sobre los mods actuales`,
	Usage: "$mods [on | off | all | files]",
	Subcommands: []support.CommandDoc{
		{Name: "on", Doc: `el comando muestra mods habilitados actualmente`},
		{Name: "off", Doc: `el comando muestra mods actualmente deshabilitados`},
		{Name: "all", Doc: `el comando muestra todas las modificaciones en mod-list.json`},
		{Name: "files", Doc: `el comando muestra los nombres de archivo de todos los mods descargados`},
	},
}

func modList(ModList *ModJson, returnEnabled bool, returnDisabled bool) string {
	var enabled, disabled int
	var S = "mod"
	if len(ModList.Mods) > 1 {
		S = "mods"
	}
	for _, mod := range ModList.Mods {
		if mod.Enabled {
			enabled += 1
		} else {
			disabled += 1
		}
	}

	res := fmt.Sprintf("%d total %s (%d activados, %d desactivados)", len(ModList.Mods), S, enabled, disabled)

	if returnEnabled {
		res += "\n**Activados:**"
		any := false
		for _, mod := range ModList.Mods {
			if mod.Enabled {
				any = true
				res += "\n    " + mod.Name
			}
		}
		if !any {
			res += " **Ninguno**"
		}
	}
	if returnDisabled {
		if returnEnabled {
			res += "\n"
		}
		res += "\n**Desactivados:**"
		any := false
		for _, mod := range ModList.Mods {
			if !mod.Enabled {
				any = true
				res += "\n    " + mod.Name
			}
		}
		if !any {
			res += " **Ninguno**"
		}
	}

	return res
}

func modsFiles() string {
	res := ""
	baseDir := path.Dir(support.Config.ModListLocation)
	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		support.Critical(err, "wtf")
	}
	for _, file := range files {
		re := support.ModFileRegexp.FindString(file.Name())
		if re != "" {
			res += "\n    " + file.Name()
		}
	}
	if res == "" {
		return "**Sin mods**"
	} else {
		return "**Mods instalados:**" + res
	}
}

// ModsList returns the list of mods running on the server.
func ModsList(s *discordgo.Session, args string) {
	returnEnabled := true
	returnDisabled := false
	if args == "on" || args == "" {
		returnEnabled = true
	} else if args == "off" {
		returnEnabled = false
		returnDisabled = true
	} else if args == "all" {
		returnDisabled = true
	} else if args == "files" {
		support.Send(s, modsFiles())
		return
	} else {
		support.SendFormat(s, "Usage: "+ModListDoc.Usage)
		return
	}
	ModList := &ModJson{}
	Json, err := ioutil.ReadFile(support.Config.ModListLocation)
	if err != nil {
		support.Send(s, "Lo sentimos, hubo un error al leer tu lista de mods.")
		support.Panik(err, "hubo un error al leer la lista de mods, ¿lo especificó en el archivo config.json?")
		return
	}

	err = json.Unmarshal(Json, &ModList)
	if err != nil {
		support.Send(s, "Lo sentimos, hubo un error al leer tu lista de mods.")
		support.Panik(err, "hubo un error al leer la lista de mods")
		return
	}
	support.ChunkedMessageSend(s, modList(ModList, returnEnabled, returnDisabled))
	return
}
