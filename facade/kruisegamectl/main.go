package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/CloudNativeGame/kruise-game-api/pkg/filter"
	"github.com/CloudNativeGame/kruise-game-api/pkg/options"
	"github.com/CloudNativeGame/kruise-game-api/pkg/updater"
	"log/slog"
	"os"
	"time"
)

var opts = &slog.HandlerOptions{AddSource: true, Level: slog.LevelError}
var logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))

func init() {
	slog.SetDefault(logger)
}

func main() {
	var rawFilter string
	var jsonPatch string
	var kubeConfigPath string
	var pretty bool

	flag.StringVar(&rawFilter, "filter", "", "filter for the game servers")
	flag.StringVar(&jsonPatch, "jsonpatch", "", "jsonPatch for the game servers")
	flag.StringVar(&kubeConfigPath, "kubeconfig", "", "path of the kube config")
	flag.BoolVar(&pretty, "pretty", false, "whether to prettify the response JSON")

	flag.Parse()

	kubeOption := options.KubeOption{
		KubeConfigPath:          kubeConfigPath,
		InformersReSyncInterval: time.Second * 30,
	}

	f := filter.NewFilter(&filter.FilterOption{
		KubeOption: kubeOption,
	})

	// TODO: fix need sleep
	time.Sleep(time.Second)

	gameServers, err := f.GetFilteredGameServers(rawFilter)
	if err != nil {
		slog.Error("get filtered GameServers failed", "error", err)
		os.Exit(1)
	}

	if jsonPatch != "" {
		u := updater.NewUpdater(&updater.UpdaterOption{
			KubeOption: kubeOption,
		})

		results := u.Update(gameServers, []byte(jsonPatch))
		var resultsJson []byte
		if pretty {
			resultsJson, err = json.MarshalIndent(results, "", "    ")
		} else {
			resultsJson, err = json.Marshal(results)
		}
		if err != nil {
			slog.Error("marshal GameServers update results failed", "error", err)
			os.Exit(1)
		}
		fmt.Println(string(resultsJson))
	} else {
		var resultsJson []byte
		if pretty {
			resultsJson, err = json.MarshalIndent(gameServers, "", "    ")
		} else {
			resultsJson, err = json.Marshal(gameServers)
		}
		if err != nil {
			slog.Error("marshal GameServers failed", "error", err)
			os.Exit(1)
		}

		fmt.Println(string(resultsJson))
	}
}
