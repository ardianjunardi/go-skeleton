package main

import (
	"fmt"
	"go-skeleton/bootstrap"
	"go-skeleton/lib/psql"
	"go-skeleton/lib/utils"
	"go-skeleton/services/api"
	"log"
	"os"

	"path/filepath"
	"runtime"

	"github.com/urfave/cli/v2"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
	config     utils.Config
	debug      = false

	// app the base of skeleton
	app *bootstrap.App
)

// EnvConfigPath environtment variable that set the config path
const EnvConfigPath = "REBEL_CLI_CONFIG_PATH"

// setup initialize the used variable and dependencies
func setup() {
	configFile := os.Getenv(EnvConfigPath)
	if configFile == "" {
		configFile = "./config.json"
	}

	log.Println(configFile)

	config = utils.NewViperConfig(basepath, configFile)

	debug = config.GetBool("app.debug")
	validator := bootstrap.SetupValidator(config)
	cLog := bootstrap.SetupLogger(config)

	// connect to redis cache
	rdCache, err := bootstrap.SetupRedis(
		config.GetString("db.redis.addr"),
		config.GetString("db.redis.password"),
		1,
	)
	if err != nil {
		fmt.Println("[redis-cache] " + err.Error())
	}

	// connect to database
	db, err := psql.Connect(config.GetString("db.psql_dsn"))
	if err != nil {
		panic(err)
	}

	app = &bootstrap.App{
		Debug:     debug,
		Config:    config,
		Validator: validator,
		Log:       cLog,
		DB:        db,
		Redis:     rdCache,
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	setup()

	// add new service
	app.AddService(api.Booting(app), "api", "API service")

	cmd := &cli.App{
		Name:     "Verein Core",
		Usage:    "Verein Core, cli",
		Commands: app.ServiceCmd,
		Action: func(cli *cli.Context) error {
			fmt.Printf("%s version@%s\n", cli.App.Name, "2.1")
			return nil
		},
	}

	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
