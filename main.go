package main

import (
	"fmt"
	"os"

	"github.com/vdimir/markify/app"
	"github.com/jessevdk/go-flags"
)

//go:generate $GOPATH/bin/statik --src=./assets --dest=./ -p app -f

var version = "local"
var date = "local"
var commit = "local"

// Opts contains command line options (see go-flags for details)
type Opts struct {
	Hostname string `short:"h" long:"host" required:"false" description:"server host name" env:"SERVER_HOSTNAME"`
	Port     uint16 `short:"p" long:"port" required:"false" description:"server port" env:"SERVER_PORT" default:"8080"`
	Debug    bool   `long:"debug" description:"debug mode"`
}

func main() {
	fmt.Printf("Running version %s (%s)\n", version, commit)
	var opts Opts

	_, err := flags.Parse(&opts)

	if err != nil {
		os.Exit(1)
	}

	appServer, err := app.NewApp(app.Config{
		ServerAddrHost: opts.Hostname,
		ServerPort:     opts.Port,
		Debug:          opts.Debug,
		AssetsPrefix:   "assets",
		PageCachePath:  "page.db",
		KeyStorePath:   "keys.db",
		MdTextPath:     "mdtext.db",
	})

	if err != nil {
		panic(err)
	}

	appServer.StartServer()
}
