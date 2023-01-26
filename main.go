package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

type App struct {
	Config Config
	Q      Queue
}

func (c Config) FindScript(repo string) Script {
	var s Script
	for _, script := range c.Scripts {
		if script.Repository == repo {
			s = script
		}
	}

	return s
}

type Config struct {
	Scripts []Script `json:"scripts"`
}

type Script struct {
	Name       string `json:"name"`
	Repository string `json:"repository"`
	Script     string `json:"script"`
}

func (a *App) loadConfig(path string) {
	viper.SetConfigType("json")
	viper.SetConfigFile(path)
	viper.ReadInConfig()
	viper.Unmarshal(&a.Config)
}

func main() {
	port := flag.String("port", ":9090", "Port to listen on")
	configPath := flag.String("config", "./.config.json", "Path to .config.json")
	flag.Parse()

	app := App{
		Q: Queue{
			Jobs: make(chan Job),
		},
	}
	app.loadConfig(*configPath)

	app.Q.Start()

	fmt.Printf("%+v\n", app.Config.Scripts)

	http.HandleFunc("/hook", app.handleHook)

	fmt.Printf("Server listening on %s\n", *port)
	http.ListenAndServe(*port, nil)
}

type HookData struct {
	Event      string `json:"event"`
	Repository string `json:"repository"`
	Commit     string `json:"commit"`
	Ref        string `json:"ref"`
	Head       string `json:"head"`
	Workflow   string `json:"workflow"`
	RequestID  string `json:"requestID"`
}

func (a *App) handleHook(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var data HookData
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err)
	}

	script := a.Config.FindScript(data.Repository)
	j := Job{
		Repository: script.Repository,
		Script:     script.Script,
	}
	a.Q.Jobs <- j
	w.Write([]byte("started"))
}
