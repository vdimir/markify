package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/vdimir/markify/app"
)

//go:generate $GOPATH/bin/statik --src=./assets --dest=./ -p app -f

var revision = "local"

// Opts contains command line options (see go-flags for details)
type Opts struct {
	Hostname string `short:"h" long:"host" required:"false" description:"server host name" env:"SERVER_HOSTNAME"`
	Port     uint16 `short:"p" long:"port" required:"false" description:"server port" env:"SERVER_PORT" default:"8080"`
	DataDir  string `short:"d" long:"data" required:"false" description:"path to directory with data" env:"DB_DATA_PATH"`
	Debug    bool   `long:"debug" description:"debug mode"`
}

func main() {
	fmt.Printf("Running version %s\n", revision)
	var opts Opts

	_, err := flags.Parse(&opts)

	if err != nil {
		os.Exit(1)
	}

	appServer, err := app.NewApp(&app.Config{
		Debug:        opts.Debug,
		AssetsPrefix: "assets",
		DBPath:       opts.DataDir,
		StatusText:   fmt.Sprintf(`{"revision":"%s"}`, revision),
	}, nil)

	if err != nil {
		panic(err)
	}

	appServer.StartServer(opts.Hostname, opts.Port)
}
