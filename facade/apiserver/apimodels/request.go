package apimodels

type UpdateGameServersRequest struct {
	Filter    string `json:"filter"`
	JsonPatch string `json:"jsonPatch"`
}

type UpdateGameServerSetsRequest struct {
	Filter    string `json:"filter"`
	JsonPatch string `json:"jsonPatch"`
}
