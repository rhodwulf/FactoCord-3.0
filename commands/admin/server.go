package admin

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/maxsupermanhd/FactoCord-3.0/support"
	"io"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ServerCommandDoc = support.CommandDoc{
	Name: "server",
	Usage: "$server\n" +
		"$server [stop|start|restart|update <version>?]",
	Doc: "el comando administra el servidor de factorio.\n" +
		"`$server` muestra el estado actual del servidor. cualquiera puede ejecutarlo.`",
	Subcommands: []support.CommandDoc{
		{Name: "stop", Doc: `el comando detiene el servidor`},
		{Name: "start", Doc: `comando inicia el servidor`},
		{Name: "restart", Doc: `el comando reinicia el servidor`},
		{
			Name: "update",
			Doc:  `el comando actualiza el servidor a la versión más reciente o a la versión especificada`,
			Usage: "$server update\n" +
				"$server update <version>",
		},
	},
}

func ServerCommandAdminPermission(args string) bool {
	return strings.TrimSpace(args) != ""
}

func ServerCommand(s *discordgo.Session, args string) {
	action, arg := support.SplitDivide(args, " ")
	switch action {
	case "":
		if support.Factorio.IsRunning() {
			support.Send(s, "El servidor de Factorio esta **en línea**")
		} else {
			support.Send(s, "El servidor de Factorio esta **detenido**")
		}
	case "stop":
		support.Factorio.Stop(s)
	case "start":
		support.Factorio.Start(s)
	case "restart":
		support.Factorio.Stop(s)
		support.Factorio.Start(s)
	case "update":
		serverUpdate(s, arg)
	default:
		support.SendFormat(s, "Uso: "+ServerCommandDoc.Usage)
	}
}

func serverUpdate(s *discordgo.Session, version string) {
	if support.Factorio.IsRunning() {
		support.Send(s, "Primero debe detener el servidor.")
		return
	}
	//username := support.Config.Username
	//token := support.Config.ModPortalToken
	//if username == "" {
	//	support.Send(s, "Username is required for update")
	//	return
	//}
	//if token == "" {
	//	support.Send(s, "Token is required for update")
	//	return
	//}
	factorioVersion, err := support.FactorioVersion()
	if err != nil {
		support.Panik(err, "... comprobando la versión de factorio")
		support.Send(s, "Error al comprobar la versión de factorio")
		return
	}

	if version == "" {
		version, err = getLatestVersion()
		if err != nil {
			support.Panik(err, "Error al obtener la información de la última versión")
			support.Send(s, "Error al obtener la información de la última versión")
			return
		}
		if version == factorioVersion {
			support.Send(s, "El servidor ya está actualizado a la última versión.")
			return
		}
	} else if version == factorioVersion {
		support.Send(s, "El servidor ya está actualizado a esa versión.")
		return
	}

	resp, err := http.Get(fmt.Sprintf("https://updater.factorio.com/get-download/%s/headless/linux64", version))
	if err != nil {
		support.Panik(err, "Error de conexion al descargar factorio")
		support.Send(s, "Se produjo algún error de conexión")
		return
	}
	if resp.StatusCode == 404 {
		support.Send(s, fmt.Sprintf("Version %s no encontrada\n"+
			"Referirse a <https://factorio.com/download/archive> para ver las versiones disponibles", version))
		return
	}
	if resp.ContentLength <= 0 {
		support.Send(s, "Error con content-length")
		return
	}
	if resp.Header.Get("Content-Disposition") == "" {
		support.Send(s, "Error con content-disposition")
		return
	}
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		support.Send(s, "Error con content-disposition")
		return
	}
	filename, ok := params["filename"]
	if !ok {
		support.Send(s, "Error con content-disposition")
		return
	}
	path := "/tmp/" + filename

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	if err != nil {
		support.Panik(err, "Error abriendo "+path+" for write")
		support.Send(s, path+": error al abrir el archivo para escribir")
		return
	}

	message := support.Send(s, support.FormatNamed(support.Config.Messages.DownloadStart, "file", filename))
	counter := &support.WriteCounter{Total: uint64(resp.ContentLength)}
	progress := support.ProgressUpdate{
		WriteCounter: counter,
		Message:      message,
		Progress:     support.FormatNamed(support.Config.Messages.DownloadProgress, "file", filename),
		Finished:     support.FormatNamed(support.Config.Messages.Unpacking, "file", filename),
	}
	go support.DownloadProgressUpdater(s, &progress)

	_, err = io.Copy(io.MultiWriter(file, counter), resp.Body)
	resp.Body.Close()
	file.Close()
	if err != nil {
		counter.Error = true
		message.Edit(s, ":interrobang: Error descargando "+filename)
		support.Panik(err, "Error descargando archivo")
		return
	}

	dir, err := filepath.Abs(support.Config.Executable)
	if err != nil {
		support.Panik(err, "Error al obtener la ruta absoluta del ejecutable")
		support.Send(s, "Error al obtener la ruta absoluta del ejecutable")
		return
	}
	dir = filepath.Dir(dir) // x64
	dir = filepath.Dir(dir) // bin
	dir = filepath.Dir(dir) // factorio
	cmd := exec.Command("tar", "-C", dir, "--strip-components=1", "-xf", path)
	err = cmd.Run()
	if err != nil {
		support.Panik(err, "Error al ejecutar tar para descomprimir el archivo")
		support.Send(s, "Error running tar to unpack the archive")
		return
	}

	message.Edit(s, support.FormatNamed(support.Config.Messages.UnpackingComplete, "version", version))
	_ = os.Remove(path)
}

type latestVersions struct {
	Stable, Experimental struct {
		Alpha, Demo, Headless string
	}
}

func getLatestVersion() (string, error) {
	resp, err := http.Get("https://factorio.com/api/latest-releases")
	if err != nil {
		return "", err
	}
	var versions latestVersions
	err = json.NewDecoder(resp.Body).Decode(&versions)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	return versions.Experimental.Headless, nil
}
