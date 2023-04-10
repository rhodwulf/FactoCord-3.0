# Este archivo contiene instrucciones de instalación para sistemas basados en Debian
En este tutorial se utiliza el servidor de factorio "headless" (de factorio.com) en `/home/factorio`.

Probado en Ubuntu 18.04.4 LTS, Ubuntu 18.04.4 LTS (Servidor), WSL1 Ubuntu

# Instalar binarios prediseñados

## Paso 0

Descargue el ejecutable de [un lanzamiento en github] (https://github.com/maxsupermanhd/FactoCord-3.0/releases)

## Step 1
Configurando

Ingrese el directorio de instalación `cd FactoCord/` (o cualquier nombre que desee)

Copie `config-example.json` a `config.json` (`cp config-example.json config.json`)

Abra `config.json` con cualquier editor (por ejemplo, `nano config.json`)

Luego en el editor de texto debes configurar:
1. Tu token de Discord para el bot (`discord_token`)
2. ID del canal de factorio para chatear (`factorio_channel_id`)
3. Parámetros de lanzamiento (flags al ejecutable de factorio) (`launch_parameters`)
4. Ruta del ejecutable (`ejecutable`)
5. Discord ID de todos los administradores (para comandos) (`admin_ids`)
6. Ubicación del archivo `mod-list.json` (incluido el nombre del archivo) (`mod_list_location`)

## Paso 2
Ejecutar

```
chmod +x ./FactoCord3
./FactoCord3
```


# Instalación desde fuentes

## Paso 0
Instalación de deps

Asegúrese de que el sistema esté actualizado `sudo apt-get update -y && sudo apt-get upgrade -y`

Descargue go 1.12+ (`sudo apt install golang-go git -y`) (es posible que deba obtenerlo del sitio web, los repositorios pueden estar desactualizados)

Obtener paquetes go: `go get`

## Paso 1
Repositorio de github

`git clone https://github.com/maxsupermanhd/FactoCord-3.0.git`

## Paso 2
Configurando

Introduzca el directorio creado `cd FactoCord-3.0/`

Copie `config-example.json` a `config.json` (`cp config-example.json config.json`)

Abra `config.json` con cualquier editor (por ejemplo, `nano config.json`)

Luego en el editor de texto debes configurar:
1. Tu token de Discord para el bot (discord_token)
2. ID del canal de factorio para chatear (factorio_channel_id)
3. Parámetros de lanzamiento (flags to factorio ejecutable) (launch_parameters)
4. Ruta ejecutable (ejecutable)
5. ID de administrador (para comandos) (admin_ids)
6. Ubicación del archivo .json de la lista de mods (incluido el nombre del archivo) (mod_list_location)

# Paso 3
Compilar

`go build`

## Etapa 4
Ejecutar

`./FactoCord-3.0`

# Usar el soporte de escenarios
... eventualmente deshabilitará los logros, pero tendrás un chat agradable y claro en Discord.
Para que se muestren personas expulsadas y poder personalizar/modificar potencialmente los mensajes
utilice control.lua **addition** desde la raíz del repositorio.
Si no desea que control.lua se modifique con tanta fuerza, puede colocarlo cerca y utilizar
un `require` para obtener la función de envío de discordia (envoltorio) y tener la funcionalidad completa de FactoCord.
