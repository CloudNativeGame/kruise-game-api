package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/CloudNativeGame/kruise-game-api/facade/rest/apimodels"
	"github.com/CloudNativeGame/kruise-game-api/pkg/deleter"
	filterbuilder "github.com/CloudNativeGame/kruise-game-api/pkg/filter/builder"
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
	return getResources[*v1alpha1.GameServer](g.httpClient, filterBuilder.Build(), g.serverUrl+"/v1/gameservers")
}

func (g *KruiseGameApiHttpClient) GetGameServerSets(filterBuilder *filterbuilder.GssFilterBuilder) ([]*v1alpha1.GameServerSet, error) {
	return getResources[*v1alpha1.GameServerSet](g.httpClient, filterBuilder.Build(), g.serverUrl+"/v1/gameserversets")
}

func (g *KruiseGameApiHttpClient) UpdateGameServers(filterBuilder *filterbuilder.GsFilterBuilder, jsonPatchBuilder *jsonpatchbuilder.GsJsonPatchBuilder) ([]updater.UpdateGsResult, error) {
	return updateResources[apimodels.UpdateGameServersRequest, updater.UpdateGsResult](g.httpClient, apimodels.UpdateGameServersRequest{
		Filter:    filterBuilder.Build(),
		JsonPatch: string(jsonPatchBuilder.Build()),
	}, g.serverUrl+"/v1/gameservers")
}

func (g *KruiseGameApiHttpClient) UpdateGameServerSets(filterBuilder *filterbuilder.GsFilterBuilder, jsonPatchBuilder *jsonpatchbuilder.GsJsonPatchBuilder) ([]updater.UpdateGssResult, error) {
	return updateResources[apimodels.UpdateGameServerSetsRequest, updater.UpdateGssResult](g.httpClient, apimodels.UpdateGameServerSetsRequest{
		Filter:    filterBuilder.Build(),
		JsonPatch: string(jsonPatchBuilder.Build()),
	}, g.serverUrl+"/v1/gameserversets")
}

func (g *KruiseGameApiHttpClient) DeleteGameServers(filterBuilder *filterbuilder.GsFilterBuilder) ([]deleter.DeleteGsResult, error) {
	return deleteResources[deleter.DeleteGsResult](g.httpClient, filterBuilder.Build(), g.serverUrl+"/v1/gameservers")
}

func (g *KruiseGameApiHttpClient) DeleteGameServerSets(filterBuilder *filterbuilder.GssFilterBuilder) ([]deleter.DeleteGssResult, error) {
	return deleteResources[deleter.DeleteGssResult](g.httpClient, filterBuilder.Build(), g.serverUrl+"/v1/gameserversets")
}

type KruiseGameResource interface {
	*v1alpha1.GameServer | *v1alpha1.GameServerSet
}

type UpdateKruiseGameRequest interface {
	apimodels.UpdateGameServersRequest | apimodels.UpdateGameServerSetsRequest
}

type UpdateKruiseGameResult interface {
	updater.UpdateGsResult | updater.UpdateGssResult
}

type DeleteKruiseGameResult interface {
	deleter.DeleteGsResult | deleter.DeleteGssResult
}

func getResources[T KruiseGameResource](httpClient *http.Client, rawFilter, rawUrl string) ([]T, error) {
	params := url.Values{}
	params.Add("filter", rawFilter)
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()
	resp, err := httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var resources []T
	err = json.Unmarshal(respBody, &resources)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func updateResources[TReq UpdateKruiseGameRequest, TRes UpdateKruiseGameResult](httpClient *http.Client, request TReq, rawUrl string) ([]TRes, error) {
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Post(rawUrl, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var updateResults []TRes
	err = json.Unmarshal(respBody, &updateResults)
	if err != nil {
		return nil, err
	}

	return updateResults, nil
}

func deleteResources[TRes DeleteKruiseGameResult](httpClient *http.Client, rawFilter, rawUrl string) ([]TRes, error) {
	params := url.Values{}
	params.Add("filter", rawFilter)
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var deleteResults []TRes
	err = json.Unmarshal(respBody, &deleteResults)
	if err != nil {
		return nil, err
	}
	return deleteResults, nil
}
