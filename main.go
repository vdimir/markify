package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/vdimir/markify/app"
)

var revision = "local"

// Opts contains command line options (see go-flags for details)
type Opts struct {
	Hostname      string `short:"h" long:"host" required:"false" description:"server host name" env:"MARKIFY_SERVER_HOSTNAME"`
	Port          uint16 `short:"p" long:"port" required:"false" description:"server port" env:"MARKIFY_SERVER_PORT" default:"8080"`
	Storage       string `short:"s" long:"storage" required:"false" description:"storage specification '<type_of_storage>:<config>', see storage code for details" env:"MARKIFY_STORAGE" default:"local:./"`
	AdminPassword string `long:"admin_secret" required:"false" description:"Admin credential to access /_admin endpoint" env:"MARKIFY_ADMIN_PWD"`
	SecretSeed    string `long:"seed_secret" required:"false" description:"Secret seed to generate tokens" env:"MARKIFY_SEED"`
	Debug         bool   `long:"debug" description:"debug mode"`
}

func main() {
	log.Printf("[DEBUG] Starting app version %s\n", revision)
	var opts Opts

	_, err := flags.Parse(&opts)

	if err != nil {
		os.Exit(1)
	}

	appServer, err := app.NewApp(&app.Config{
		Debug:         opts.Debug,
		AssetsPrefix:  "app/assets",
		StorageSpec:   opts.Storage,
		StatusText:    fmt.Sprintf(`{"revision":"%s"}`, revision),
		AdminPassword: opts.AdminPassword,
		UIDSecret:     opts.SecretSeed,
	})

	if err != nil {
		log.Printf("[ERROR] App can't be created: %s", err)
		panic(err)
	}

	appServer.StartServer(opts.Hostname, opts.Port)

	log.Printf("[DEBUG] App closed")

}
