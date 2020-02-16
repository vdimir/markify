package main

import (
	"fmt"
	"os"
	"path"

	"github.com/jessevdk/go-flags"
	"github.com/vdimir/markify/app"
)

//go:generate $GOPATH/bin/statik --src=./assets --dest=./ -p app -f

var version = "local"
var date = "local"
var commit = "local"

// Opts contains command line options (see go-flags for details)
type Opts struct {
	Hostname string `short:"h" long:"host" required:"false" description:"server host name" env:"SERVER_HOSTNAME"`
	Port     uint16 `short:"p" long:"port" required:"false" description:"server port" env:"SERVER_PORT" default:"8080"`
	DataDir  string `short:"d" long:"data" required:"false" description:"path to directory with data"`
	Debug    bool   `long:"debug" description:"debug mode"`
}

func main() {
	fmt.Printf("Running version %s (%s) -- %s\n", version, commit, date)
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
		PageCachePath:  path.Join(opts.DataDir, "page.db"),
		URLHashPath:    path.Join(opts.DataDir, "url_hash.db"),
		MdTextPath:     path.Join(opts.DataDir, "mdtext.db"),
		StatusText: fmt.Sprintf(`{`+
			`"commit":"%s",`+
			`"version":"%s",`+
			`"date":"%s"`+
			`}`, commit, version, date),
	})

	if err != nil {
		panic(err)
	}

	appServer.StartServer()
}
