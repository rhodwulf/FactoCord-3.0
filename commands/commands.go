package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maxsupermanhd/FactoCord-3.0/commands/admin"
	"github.com/maxsupermanhd/FactoCord-3.0/commands/utils"
	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

// Command is a struct containing fields that hold command information.
type Command struct {
	Name string

	Command func(s *discordgo.Session, args string)

	Admin func(args string) bool
	Doc   *support.CommandDoc
	Desc  string
}

func alwaysAdmin(_ string) bool {
	return true
}

// Commands is a list of all available commands
var Commands = [...]Command{
	// Admin Commands
	{
		Name:    "server",
		Command: admin.ServerCommand,
		Admin:   admin.ServerCommandAdminPermission,
		Doc:     &admin.ServerCommandDoc,
		Desc:    "Administrar el servidor de factorio",
	},
	{
		Name:    "save",
		Command: admin.SaveServer,
		Admin:   alwaysAdmin,
		Doc:     &admin.SaveServerDoc,
		Desc:    "Guarda el juego",
	},
	{
		Name:    "kick",
		Command: admin.KickPlayer,
		Admin:   alwaysAdmin,
		Doc:     &admin.KickPlayerDoc,
		Desc:    "Explulsa a un jugador del servidor",
	},
	{
		Name:    "ban",
		Command: admin.BanPlayer,
		Admin:   alwaysAdmin,
		Doc:     &admin.BanPlayerDoc,
		Desc:    "Banea a un jugador del servidor",
	},
	{
		Name:    "unban",
		Command: admin.UnbanPlayer,
		Admin:   alwaysAdmin,
		Doc:     &admin.UnbanPlayerDoc,
		Desc:    "Desbanea a un usuario del servidor",
	},
	{
		Name:    "config",
		Command: admin.ConfigCommand,
		Admin:   alwaysAdmin,
		Doc:     &admin.ConfigCommandDoc,
		Desc:    "Administra config.json",
	},
	{
		Name:    "mod",
		Command: admin.ModCommand,
		Admin:   alwaysAdmin,
		Doc:     &admin.ModCommandDoc,
		Desc:    "Administra mod-list.json",
	},

	// Util Commands
	{
		Name:    "mods",
		Command: utils.ModsList,
		Admin:   nil,
		Doc:     &utils.ModListDoc,
		Desc:    "Crea una lista de los mods del servidor",
	},
	{
		Name:    "version",
		Command: utils.VersionString,
		Admin:   nil,
		Doc:     &utils.VersionDoc,
		Desc:    "Muestra la versión del servidor",
	},
	{
		Name:    "info",
		Command: utils.GameInfo,
		Admin:   nil,
		Doc:     &utils.InfoDoc,
		Desc:    "Muestra información del servidor",
	},
	{
		Name:    "online",
		Command: utils.GameOnline,
		Admin:   nil,
		Doc:     &utils.OnlineDoc,
		Desc:    "Muestra jugadores en línea",
	},
	{
		Name:  "help",
		Admin: nil,
		Doc: &support.CommandDoc{
			Name: "help",
			Usage: "$help\n" +
				"$help <command>\n" +
				"$help <command> <subcommand>",
			Doc: "comando devuelve una lista de todos los comandos y documentación sobre cualquier comando y sus subcomandos",
		},
		Desc: "Enumera los comandos para el servidor y obtenga documentación sobre comandos y subcomandos. Intenta usar `$help help`",
	},
}

func helpCommand(s *discordgo.Session, args string) {
	if args == "" {
		helpAllCommands(s)
		return
	}
	args = strings.ToLower(args)
	commandName, subcommand := support.SplitDivide(args, " ")
	for _, command := range Commands {
		if command.Name == commandName {
			helpOnCommand(s, command.Doc, subcommand)
			return
		}
	}
	support.Send(s, "No existe el comando \""+commandName+"\"")
}

func helpOnCommand(s *discordgo.Session, command *support.CommandDoc, subcommandName string) {
	path := support.Config.Prefix + command.Name
	if subcommandName != "" {
		found := false
		for _, subcommand := range command.Subcommands {
			if subcommand.Name == subcommandName {
				command = &subcommand
				path += " " + subcommandName
				found = true
				break
			}
		}
		if !found {
			support.Send(s, fmt.Sprintf(`El comando "%s" no tiene subcomando "%s"`, command.Name, subcommandName))
			return
		}
	}
	quoted := "`" + path + "`"
	embed := &discordgo.MessageEmbed{
		Type:        "rich",
		Color:       0x6289FF,
		Title:       fmt.Sprintf("Documentación sobre `%s` comando", path),
		Description: support.FormatUsage(quoted + " " + command.Doc),
	}
	usage := command.Usage
	if usage == "" {
		usage = path
	}
	usage = support.FormatUsage(usage)
	if strings.Contains(usage, "\n") {
		usage = "```\n" + usage + "\n```"
	} else {
		usage = "`" + usage + "`"
	}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "**Usage:**",
		Value: usage,
	})
	if len(command.Subcommands) > 0 {
		subcommands := ""
		for _, subcommand := range command.Subcommands {
			if subcommands != "" {
				subcommands += "\n"
			}
			subcommands += path + " " + subcommand.Name
		}
		subcommands = "```\n" + subcommands + "\n```"
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "**Subcommands:**",
			Value: support.FormatUsage(subcommands),
		})
	}
	support.SendEmbed(s, embed)
}

func helpAllCommands(s *discordgo.Session) {
	fields := make([]*discordgo.MessageEmbedField, 0, len(Commands))

	for _, command := range Commands {
		desc := support.FormatUsage(command.Desc)
		if roleID, exists := support.Config.CommandRoles[strings.ToLower(command.Name)]; exists {
			roles, err := s.GuildRoles(support.GuildID)
			if err != nil {
				support.Panik(err, "... al consultar los roles del gremio")
				return
			}
			found := false
			for _, role := range roles {
				if role.ID == roleID {
					found = true
					desc = "[Role \"" + role.Name + "\"] " + desc
					break
				}
			}
			if !found {
				desc = "[Rol no encontrado en el gremio] " + desc
			}
		} else if command.Admin != nil {
			desc = "[Admin] " + desc
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  support.Config.Prefix + command.Name,
			Value: desc,
		})
	}
	embed := &discordgo.MessageEmbed{
		Type:        "rich",
		Color:       52,
		Description: "Lista de todos los comandos actualmente disponibles",
		Title:       "Comandos del servidor",
		Fields:      fields,
	}
	support.SendEmbed(s, embed)
}

// RunCommand runs a specified command.
func RunCommand(input string, s *discordgo.Session, m *discordgo.Message) {
	inputvars := strings.SplitN(input+" ", " ", 2)
	commandName := strings.ToLower(inputvars[0])
	args := strings.TrimSpace(inputvars[1])

	if commandName == strings.ToLower("Help") {
		helpCommand(s, args)
		return
	}

	for _, command := range Commands {
		if strings.ToLower(command.Name) == commandName {
			execute := false
			err := ""

			if command.Admin != nil && command.Admin(args) {
				if CheckAdmin(m.Author.ID) {
					execute = true
				} else {
					err = "¡No eres un administrador!"
				}
			} else {
				execute = true
			}
			if roleID, exists := support.Config.CommandRoles[commandName]; exists {
				// TODO? role name
				err = "No tienes el rol requerido"
				for _, memberRoleID := range m.Member.Roles {
					if memberRoleID == roleID {
						execute = true
					}
				}
			}
			if execute {
				command.Command(s, args)
			} else {
				support.Send(s, err)
			}
			return
		}
	}
	support.SendFormat(s, "Comando no encontrado. Intenta usar \"$help\"")
}

// CheckAdmin checks if the user attempting to run an admin command is an admin
func CheckAdmin(ID string) bool {
	for _, adminID := range support.Config.AdminIDs {
		if ID == adminID {
			return true
		}
	}
	return false
}
