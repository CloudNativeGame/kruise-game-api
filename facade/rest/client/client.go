package client

import (
	"bytes"
	"encoding/json"
	"github.com/CloudNativeGame/kruise-game-api/facade/rest/apimodels"
	filterbuilder "github.com/CloudNativeGame/kruise-game-api/pkg/filter"
	jsonpatchbuilder "github.com/CloudNativeGame/kruise-game-api/pkg/jsonpatches/builder"
	"github.com/CloudNativeGame/kruise-game-api/pkg/updater"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type KruiseGameApiHttpClient struct {
	httpClient *http.Client
	serverUrl  string
}

func NewKruiseGameApiHttpClient() *KruiseGameApiHttpClient {
	serverUrl := os.Getenv("SERVER_URL")
	if serverUrl == "" {
		serverUrl = "http://kruise-game-api.kruise-game-system.svc.cluster.local"
	}
	return &KruiseGameApiHttpClient{
		httpClient: &http.Client{
			Timeout: time.Duration(30) * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        0,
				MaxIdleConnsPerHost: 50000,
				MaxConnsPerHost:     0,
				IdleConnTimeout:     300 * time.Second,
			},
		},
		serverUrl: serverUrl,
	}
}

func (g *KruiseGameApiHttpClient) GetGameServers(filterBuilder *filterbuilder.GsFilterBuilder) ([]*v1alpha1.GameServer, error) {
	params := url.Values{}
	params.Add("filter", filterBuilder.Build())
	u, err := url.Parse(g.serverUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()
	resp, err := g.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var gameServers []*v1alpha1.GameServer
	err = json.Unmarshal(respBody, &gameServers)
	if err != nil {
		return nil, err
	}
	return gameServers, nil
}

func (g *KruiseGameApiHttpClient) UpdateGameServers(filterBuilder *filterbuilder.GsFilterBuilder, jsonPatchBuilder *jsonpatchbuilder.GsJsonPatchBuilder) ([]updater.UpdateResult, error) {
	request := apimodels.UpdateGameServersRequest{
		Filter:    filterBuilder.Build(),
		JsonPatch: string(jsonPatchBuilder.Build()),
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := g.httpClient.Post(g.serverUrl, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var updateResults []updater.UpdateResult
	err = json.Unmarshal(respBody, &updateResults)
	if err != nil {
		return nil, err
	}

	return updateResults, nil
}
