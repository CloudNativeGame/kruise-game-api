package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CloudNativeGame/kruise-game-api/pkg/deleter"
	"github.com/CloudNativeGame/kruise-game-api/pkg/filter"
	"github.com/CloudNativeGame/kruise-game-api/pkg/options"
	"github.com/CloudNativeGame/kruise-game-api/pkg/updater"
	"github.com/jessevdk/go-flags"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
	"log/slog"
	"os"
	"time"
)

type CmdOption struct {
	ResourceKind   string `short:"r" long:"resourcekind" description:"resource kind: 'gameserver' or 'gameserverset'" required:"true"`
	Filter         string `short:"f" long:"filter" description:"filter for the resources;\nupdate or delete the resources that match the filter when 'jsonpatch' or 'deletion' parameter is provided" required:"false"`
	JsonPatch      string `short:"j" long:"jsonpatch" description:"JsonPatch for the resources; fail and do nothing when the 'deletion' parameter appears at the same time" required:"false"`
	Deletion       bool   `short:"d" long:"deletion" description:"whether to delete resources; fail and do nothing when the 'jsonpatch' parameter appears at the same time" required:"false"`
	IsPretty       bool   `short:"p" long:"pretty" description:"whether to prettify the response JSON" required:"false"`
	KubeConfigPath string `short:"k" long:"kubeconfig" description:"path of the kube config" required:"false"`
}

func (c *CmdOption) Check() error {
	if c.ResourceKind != "gameserverset" && c.ResourceKind != "gameserver" {
		return errors.New("resource kind must be gameserver or gameserverset")
	}
	if c.Deletion && c.JsonPatch != "" {
		return errors.New("--jsonpatch and --deletion cannot be used together")
	}
	return nil
}

var opts = &slog.HandlerOptions{AddSource: false, Level: slog.LevelError}
var logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))

func init() {
	slog.SetDefault(logger)
}

func main() {
	var cmdOption CmdOption
	_, err := flags.Parse(&cmdOption)
	if err != nil {
		slog.Error("parse parameters failed", "error", err)
		os.Exit(1)
	}

	err = cmdOption.Check()
	if err != nil {
		slog.Error("check parameters failed", "error", err)
		os.Exit(1)
	}

	kubeOption := options.KubeOption{
		KubeConfigPath:          cmdOption.KubeConfigPath,
		InformersReSyncInterval: time.Second * 30,
	}

	if cmdOption.ResourceKind == "gameserverset" {
		if cmdOption.JsonPatch != "" {
			result, code, err := patchGameServerSets(&kubeOption, &cmdOption)
			if err != nil {
				slog.Error("patch GameServerSets failed", "error", err.Error())
			} else {
				fmt.Println(result)
			}
			os.Exit(code)

		} else if cmdOption.Deletion {
			result, code, err := deleteGameServerSets(&kubeOption, &cmdOption)
			if err != nil {
				slog.Error("delete GameServerSets failed", "error", err.Error())
			} else {
				fmt.Println(result)
			}
			os.Exit(code)
		} else {
			result, code, err := getGameServerSets(&kubeOption, &cmdOption)
			if err != nil {
				slog.Error("get GameServerSets failed", "error", err.Error())
			} else {
				fmt.Println(result)
			}
			os.Exit(code)
		}
	} else if cmdOption.ResourceKind == "gameserver" {
		if cmdOption.JsonPatch != "" {
			result, code, err := patchGameServers(&kubeOption, &cmdOption)
			if err != nil {
				slog.Error("patch GameServers failed", "error", err.Error())
			} else {
				fmt.Println(result)
			}
			os.Exit(code)
		} else if cmdOption.Deletion {
			result, code, err := deleteGameServers(&kubeOption, &cmdOption)
			if err != nil {
				slog.Error("delete GameServers failed", "error", err.Error())
			} else {
				fmt.Println(result)
			}
			os.Exit(code)
		} else {
			result, code, err := getGameServers(&kubeOption, &cmdOption)
			if err != nil {
				slog.Error("get GameServers failed", "error", err.Error())
			} else {
				fmt.Println(result)
			}
			os.Exit(code)
		}
	}
}

func getGameServerSets(kubeOption *options.KubeOption, cmdOption *CmdOption) (string, int, error) {
	gameServerSets, err := getFilteredGameServerSets(kubeOption, cmdOption)
	if err != nil {
		return "", 1, err
	}

	var resultsJson []byte
	if cmdOption.IsPretty {
		resultsJson, err = json.MarshalIndent(gameServerSets, "", "    ")
	} else {
		resultsJson, err = json.Marshal(gameServerSets)
	}
	if err != nil {
		return "", 1, err
	}

	return string(resultsJson), 0, nil
}

func patchGameServerSets(kubeOption *options.KubeOption, cmdOption *CmdOption) (string, int, error) {
	u := updater.NewUpdater(&updater.UpdaterOption{
		KubeOption: *kubeOption,
	})

	gameServerSets, err := getFilteredGameServerSets(kubeOption, cmdOption)
	if err != nil {
		return "", 1, err
	}

	results := u.UpdateGameServerSets(gameServerSets, []byte(cmdOption.JsonPatch))
	var resultsJson []byte
	if cmdOption.IsPretty {
		resultsJson, err = json.MarshalIndent(results, "", "    ")
	} else {
		resultsJson, err = json.Marshal(results)
	}
	if err != nil {
		return "", 1, err
	}

	for _, result := range results {
		if result.Err != nil {
			return string(resultsJson), 2, nil
		}
	}

	return string(resultsJson), 0, nil
}

func deleteGameServerSets(kubeOption *options.KubeOption, cmdOption *CmdOption) (string, int, error) {
	d := deleter.NewDeleter(&deleter.DeleterOption{
		KubeOption: *kubeOption,
	})

	gameServerSets, err := getFilteredGameServerSets(kubeOption, cmdOption)
	if err != nil {
		return "", 1, err
	}

	results := d.DeleteGameServerSets(gameServerSets)
	var resultsJson []byte
	if cmdOption.IsPretty {
		resultsJson, err = json.MarshalIndent(results, "", "    ")
	} else {
		resultsJson, err = json.Marshal(results)
	}
	if err != nil {
		return "", 1, err
	}

	for _, result := range results {
		if result.Err != nil {
			return string(resultsJson), 2, nil
		}
	}

	return string(resultsJson), 0, nil
}

func getFilteredGameServerSets(kubeOption *options.KubeOption, cmdOption *CmdOption) ([]*v1alpha1.GameServerSet, error) {
	f := filter.NewGssFilter(&filter.FilterOption{
		KubeOption: *kubeOption,
	})

	// TODO: fix need sleep
	time.Sleep(time.Second)

	gameServerSets, err := f.GetFilteredGameServerSets(cmdOption.Filter)
	if err != nil {
		return nil, err
	}

	return gameServerSets, nil
}

func getGameServers(kubeOption *options.KubeOption, cmdOption *CmdOption) (string, int, error) {
	gameServers, err := getFilteredGameServers(kubeOption, cmdOption)
	if err != nil {
		return "", 1, err
	}

	var resultsJson []byte
	if cmdOption.IsPretty {
		resultsJson, err = json.MarshalIndent(gameServers, "", "    ")
	} else {
		resultsJson, err = json.Marshal(gameServers)
	}
	if err != nil {
		return "", 1, err
	}

	return string(resultsJson), 1, nil
}

func patchGameServers(kubeOption *options.KubeOption, cmdOption *CmdOption) (string, int, error) {
	u := updater.NewUpdater(&updater.UpdaterOption{
		KubeOption: *kubeOption,
	})

	gameServers, err := getFilteredGameServers(kubeOption, cmdOption)
	if err != nil {
		return "", 1, err
	}

	results := u.UpdateGameServers(gameServers, []byte(cmdOption.JsonPatch))
	var resultsJson []byte
	if cmdOption.IsPretty {
		resultsJson, err = json.MarshalIndent(results, "", "    ")
	} else {
		resultsJson, err = json.Marshal(results)
	}
	if err != nil {
		return "", 1, err
	}

	for _, result := range results {
		if result.Err != nil {
			return string(resultsJson), 2, nil
		}
	}

	return string(resultsJson), 0, nil
}

func deleteGameServers(kubeOption *options.KubeOption, cmdOption *CmdOption) (string, int, error) {
	d := deleter.NewDeleter(&deleter.DeleterOption{
		KubeOption: *kubeOption,
	})

	gameServers, err := getFilteredGameServers(kubeOption, cmdOption)
	if err != nil {
		return "", 1, err
	}

	results := d.DeleteGameServers(gameServers)
	var resultsJson []byte
	if cmdOption.IsPretty {
		resultsJson, err = json.MarshalIndent(results, "", "    ")
	} else {
		resultsJson, err = json.Marshal(results)
	}
	if err != nil {
		return "", 1, err
	}

	for _, result := range results {
		if result.Err != nil {
			return string(resultsJson), 2, nil
		}
	}

	return string(resultsJson), 0, nil
}

func getFilteredGameServers(kubeOption *options.KubeOption, cmdOption *CmdOption) ([]*v1alpha1.GameServer, error) {
	f := filter.NewGsFilter(&filter.FilterOption{
		KubeOption: *kubeOption,
	})

	// TODO: fix need sleep
	time.Sleep(time.Second)

	gameServers, err := f.GetFilteredGameServers(cmdOption.Filter)
	if err != nil {
		slog.Error("get filtered GameServers failed", "error", err)
		return nil, err
	}

	return gameServers, nil
}
