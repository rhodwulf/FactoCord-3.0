package admin

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"github.com/maxsupermanhd/FactoCord-3.0/support"
)

var ConfigCommandDoc = support.CommandDoc{
	Name:  "config",
	Usage: "$config save | load | get <path> | set <path> <value>?",
	Doc:   "el comando administra la configuración de FactoCord",
	Subcommands: []support.CommandDoc{
		{
			Name: "save",
			Doc: "El comando guarda la configuración de FactoCord de la memoria para `config.json`.\n" +
				"También agrega las opciones que faltan en config.json",
		},
		{
			Name: "load",
			Doc: "El comando carga la configuración desde `config.json`.\n" +
				"Cualquier cambio no guardado después del último `$config save` comando se perderá.",
		},
		{
			Name:  "get",
			Usage: "$config get <path>?",
			Doc: "El comando genera el valor de un ajuste de configuración especificado por <path>.\n" +
				"Todos los miembros de la ruta están separados por un punto '.'\n" +
				"Si la ruta está vacía, genera la configuración completa.\n" +
				"Secretos como discord_token se mantienen en secreto.\n" +
				"Ejemplos:\n" +
				"```\n" +
				"$config get\n" +
				"$config get admin_ids\n" +
				"$config get admin_ids.0\n" +
				"$config get command_roles\n" +
				"$config get command_roles.mod\n" +
				"$config get messages\n" +
				"```",
		},
		{
			Name: "set",
			Usage: "$config set <path>\n" +
				"$config set <path> <value>",
			Doc: "El comando establece el valor de un ajuste de configuración especificado por <path>.\n" +
				"Este comando puede establecer solo tipos simples como cadenas, números y booleanos.\n" +
				"Si no se especifica ningún valor, este comando elimina el valor si es posible; de lo contrario, lo establece en un zero-value (0, \"\", false).\n" +
				"Para agregar un valor a una matriz o un objeto, especifique su índice como '*' (e.g. `$config set admin_ids.* 1234`).\n" +
				"Los cambios realizados por este comando no se guardan automáticamente. Usar `$config save` para hacerlo.\n" +
				"Ejemplos:" +
				"```\n" +
				"$config set prefix !\n" +
				"$config set game_name \"Factorio 1.0\"\n" +
				"$config set ingame_discord_user_colors true\n" +
				"$config set admin_ids.0 123456789\n" +
				"$config set admin_ids.* 987654321\n" +
				"$config set command_roles.mod 55555555\n" +
				"$config set messages.server_save **:mango: Game saved!**\n" +
				"```",
		},
	},
}

// ModCommand returns the list of mods running on the server.
func ConfigCommand(s *discordgo.Session, args string) {
	if args == "" {
		support.SendFormat(s, "Uso: "+ConfigCommandDoc.Usage)
		return
	}
	action, args := support.SplitDivide(args, " ")
	args = strings.TrimSpace(args)
	if _, ok := commands[action]; !ok {
		support.SendFormat(s, "Uso: "+ConfigCommandDoc.Usage)
		return
	}
	res := commands[action](args)
	support.Send(s, res)
}

var commands = map[string]func(string) string{
	"save": save,
	"load": load,
	"get":  get,
	"set":  set,
}

func save(args string) string {
	if args != "" {
		return "Guardar no acepta argumentos"
	}
	s, err := json.MarshalIndent(support.Config, "", "    ")
	if err != nil {
		support.Panik(err, "... al convertir la configuración a json")
		return "Error al convertir config a json"
	}
	err = ioutil.WriteFile(support.ConfigPath, s, 0666)
	if err != nil {
		support.Panik(err, "... al guardar config.json")
		return "Error al guardar config.json"
	}
	return "configuración guardada"
}

func load(args string) string {
	if args != "" {
		return "La carga no acepta argumentos"
	}
	err := support.Config.Load()
	if err != nil {
		return err.Error()
	}
	return "Configuración recargada"
}

func get(args string) string {
	if strings.Contains(args, " \n\t") {
		return "¿Por qué hay espacios en la ruta?"
	}
	var value interface{}
	if args == "" {
		config := support.Config // copy
		config.DiscordToken = "mi precioso"
		config.Username = "mi precioso"
		config.ModPortalToken = "mi precioso"
		value = config
	} else {
		path := strings.Split(args, ".")
		if path[0] == "discord_token" {
			return "Shhhh, es un secreto"
		}
		x, err := walk(&support.Config, path)
		if err != nil {
			return err.Error()
		}
		value = x.Interface()
	}
	res, err := json.MarshalIndent(value, "", "    ")
	if err != nil {
		support.Panik(err, "... al convertir a json")
		return "Error al convertir a json"
	}
	return fmt.Sprintf("```json\n%s\n```", string(res))
}

func set(args string) string {
	pathS, valueS := support.SplitDivide(args, " ")
	if pathS == "" {
		return support.FormatUsage("Uso: $config set <path> <value>?")
	}
	path := strings.Split(pathS, ".")
	if len(path) == 0 {
		return "wtf??"
	}
	if path[0] == "discord_token" {
		return "¿Están tratando de lavarme el cerebro?"
	}
	name := path[len(path)-1]
	pathTo := strings.Join(path[:len(path)-1], ".")
	if pathTo == "" {
		pathTo = "."
	}
	current, err := walk(&support.Config, path[:len(path)-1])
	if err != nil {
		return err.Error()
	}
	switch current.Kind() {
	case reflect.Slice:
		if name == "*" {
			value, errs := createValue(current.Type().Elem(), valueS)
			if errs != "" {
				return pathS + errs
			}
			current.Set(reflect.Append(current, value))
		} else {
			num, err := strconv.ParseUint(name, 10, 0)
			if err != nil {
				return fmt.Sprintf("%s es matriz pero \"%s\" no es un int", pathS, name)
			}
			if current.Len() <= int(num) {
				return fmt.Sprintf("%d es mayor que %s's tamaño (%d)", num, pathS, current.Len())
			}
			if valueS == "" {
				sliceRemove(current, int(num))
			} else {
				value, errs := createValue(current.Type().Elem(), valueS)
				if errs != "" {
					return pathS + errs
				}
				current.Index(int(num)).Set(value)
			}
		}
	case reflect.Struct:
		fieldName := getFieldByTag(name, "json", current.Type())
		if fieldName == "" {
			if pathTo == "." {
				return fmt.Sprintf("config no tiene una opción llamada \"%s\"", name)
			} else {
				return fmt.Sprintf("estructura %s no tiene un campo llamado \"%s\"", pathTo, name)
			}
		}
		field := current.FieldByName(fieldName)
		value, errs := createValue(field.Type(), valueS)
		if errs != "" {
			return pathS + errs
		}
		field.Set(value)
	case reflect.Map:
		key, errs := createValue(current.Type().Key(), name)
		if errs != "" {
			return pathS + errs
		}
		var value reflect.Value
		if valueS == "" {
			value = reflect.Value{}
		} else {
			value, errs = createValue(current.Type().Elem(), valueS)
			if errs != "" {
				return pathS + errs
			}
		}
		current.SetMapIndex(key, value)
	default:
		return fmt.Sprintf("%s's tipo (%s) no es soprotad", pathS, current.Type().String())
	}
	return "Value set"
}

func walk(v interface{}, path []string) (reflect.Value, error) {
	var current = reflect.ValueOf(v)
	if current.Type().Kind() != reflect.Ptr {
		panic("walk: v should be pointer")
	}
	current = current.Elem()
	for i, name := range path {
		walkedPath := strings.Join(path[:i], ".")
		switch current.Kind() {
		case reflect.Slice:
			num, err := strconv.ParseUint(name, 10, 0)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("%s es matriz pero \"%s\" is not an int", walkedPath, name)
			}
			if current.Len() <= int(num) {
				return reflect.Value{}, fmt.Errorf("tamaño de %s[%d] es menos que %d", walkedPath, current.Len(), num)
			}
			current = current.Index(int(num))
		case reflect.Struct:
			field := getFieldByTag(name, "json", current.Type())
			if field == "" {
				return reflect.Value{}, fmt.Errorf("estructura %s no tiene un campo llamado \"%s\"", walkedPath, name)
			}
			current = current.FieldByName(field)
		case reflect.Map:
			key, errs := createValue(current.Type().Key(), name)
			if errs != "" {
				return reflect.Value{}, fmt.Errorf(walkedPath + "." + name + errs)
			}
			current = current.MapIndex(key)
			if !current.IsValid() {
				return reflect.Value{}, fmt.Errorf("%s no tiene llave \"%s\"", walkedPath, name)
			}
		default:
			return reflect.Value{}, fmt.Errorf("%s's tipo (%s) no es soportado", walkedPath, current.Type().String())
		}
	}
	return current, nil
}

func createValue(t reflect.Type, value string) (reflect.Value, string) {
	if value == "" {
		return reflect.New(t).Elem(), ""
	}
	switch t.Kind() {
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return reflect.Value{}, fmt.Sprintf(" requiere bool pero \"%s\" no se puede convertir a bool", value)
		}
		return reflect.ValueOf(val), ""
	case reflect.Int:
		num, err := strconv.ParseUint(value, 10, 0)
		if err != nil {
			return reflect.Value{}, fmt.Sprintf(" requiere int pero \"%s\" no es un int", value)
		}
		return reflect.ValueOf(int(num)), ""
	case reflect.String:
		if value[0] == '"' && value[len(value)-1] == '"' {
			return reflect.ValueOf(value[1 : len(value)-1]), ""
		}
		return reflect.ValueOf(value), ""
	default:
		return reflect.Value{}, fmt.Sprintf("'s tipo (%s) no es soportado", t.String())
	}
}

func getFieldByTag(tag, key string, t reflect.Type) (fieldname string) {
	if t.Kind() != reflect.Struct {
		panic("bad type")
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		v := strings.Split(f.Tag.Get(key), ",")[0] // use split to ignore tag "options" like omitempty, etc.
		if v == tag {
			return f.Name
		}
	}
	return ""
}

func sliceRemove(v reflect.Value, index int) {
	for i := index; i < v.Len()-1; i++ {
		v.Index(i).Set(v.Index(i + 1))
	}
	v.SetLen(v.Len() - 1)
}
