# kruise-game-api

Filtering query and update operation API for kruise-game resources.

## Overview

* This repository provides three request ways: command line, REST interface, and go language package to implement filtering and updating of kruise-game resources.
* The filter syntax is based on this repository: [structured-filter-go](https://github.com/CloudNativeGame/structured-filter-go).
* GameServer implemented the following Scene filters:

| Key                 | Value Type    | Filtered Field                            | Filter Examples                                    |
|---------------------|---------------|-------------------------------------------|----------------------------------------------------|
| namespace           | String Filter | /metadata/namespace                       | `{"namespace": {"$in":["minecraft", "terraria"]}}` |
| opsState            | String Filter | /spec/opsState                            | `{"opsState": {"$eq": "None"}}`                    |
| updatePriority      | Number Filter | /spec/updatePriority                      | `{"updatePriority": {"$ne": 0}}`                   |
| deletionPriority    | Number Filter | /spec/deletionPriority                    | `{"deletionPriority": {"$ne": 0}}`                 |
| currentState        | String Filter | /status/currentState                      | `{"currentState": {"$eq": "Ready"}}`               |
| currentNetworkState | String Filter | /status/networkStatus/currentNetworkState | `{"currentNetworkState": {"$eq": "Ready"}}`        |

* GameServerSet implemented the following Scene filters:

| Key            | Value Type    | Filtered Field       | Filter Examples                                    |
|----------------|---------------|----------------------|----------------------------------------------------|
| namespace      | String Filter | /metadata/namespace  | `{"namespace": {"$in":["minecraft", "terraria"]}}` |

* Use [JSON Patch](https://datatracker.ietf.org/doc/html/rfc6902) as update syntax.

## Usage

### REST server

* deploy:

```shell
kubectl apply -f ./deploy
```

* interfaces:

Get filtered game servers:

```
GET /v1/gameservers

Query parameters:
filter string

Response body:
[]v1alpha1.GameServer
```

Get filtered game server sets:

```
GET /v1/gameserversets

Query parameters:
filter string

Response body:
[]v1alpha1.GameServerSet
```

Update game servers:

```
POST /v1/gameservers

Request body:
UpdateGameServersRequest

Response body:
[]UpdateGsResult

type UpdateGameServersRequest struct {
	Filter    string `json:"filter"`
	JsonPatch string `json:"jsonPatch"`
}

type UpdateGsResult struct {
	Gs        *v1alpha1.GameServer `json:"gs"`
	UpdatedGs *v1alpha1.GameServer `json:"updatedGs"`
	Err       error                `json:"err"`
}
```

Update game server sets:

```
POST /v1/gameserversets

Request body:
UpdateGameServerSetsRequest

Response body:
[]UpdateGssResult

type UpdateGameServerSetsRequest struct {
	Filter    string `json:"filter"`
	JsonPatch string `json:"jsonPatch"`
}

type UpdateGssResult struct {
	Gss        *v1alpha1.GameServerSet `json:"gss"`
	UpdatedGss *v1alpha1.GameServerSet `json:"updatedGss"`
	Err        error                   `json:"err"`
}
```

* Use the HTTP client to call the REST interface. See [examples](https://github.com/CloudNativeGame/kruise-game-api/blob/main/examples/rest/curl/curl.sh).

* Also, can use the `KruiseGameApiHttpClient` provided by package `github.com/CloudNativeGame/kruise-game-api/facade/rest/client` in golang code. See [examples](https://github.com/CloudNativeGame/kruise-game-api/blob/main/examples/rest/main.go).

### Command line

* See [examples](https://github.com/CloudNativeGame/kruise-game-api/tree/main/examples/kruisegamectl/kruisegamectl.sh).

```shell
Usage of kruisegamectl:
  -filter string
        filter for the game servers
  -kubeconfig string
        path of the kube config
  -patch string
        jsonPatch for the game servers
  -pretty bool
        whether to prettify the response JSON
```
