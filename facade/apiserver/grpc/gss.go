package service

import (
	"context"
	"encoding/json"
	"github.com/CloudNativeGame/kruise-game-api/facade/rest/apimodels"
	api "github.com/CloudNativeGame/kruise-game-api/facade/rest/proto"
	"github.com/CloudNativeGame/kruise-game-api/facade/rest/service"
	"github.com/CloudNativeGame/kruise-game-api/pkg/deleter"
	"github.com/CloudNativeGame/kruise-game-api/pkg/updater"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
	"log/slog"
)

type GameServerSetGrpcService struct {
	api.UnimplementedGameServerSetServiceServer
}

func (g GameServerSetGrpcService) GetGameServerSets(ctx context.Context, request *api.GetGameServerSetsRequest) (*api.GetGameServerSetsResponse, error) {
	gameServerSets, err := service.GetGssService().GetGameServerSets(request.Filter)
	if err != nil {
		slog.Error("get filtered GameServerSets failed", "error", err)
		return nil, err
	}

	marshaledGameServerSets, err := marshalGameServerSets(gameServerSets)
	if err != nil {
		slog.Error("marshal GameServerSets failed", "error", err)
		return nil, err
	}

	return &api.GetGameServerSetsResponse{
		GameServerSets: marshaledGameServerSets,
	}, nil
}

func marshalGameServerSets(gameServerSets []*v1alpha1.GameServerSet) ([]string, error) {
	marshaledGameServerSets := make([]string, 0, len(gameServerSets))
	for _, gameServerSet := range gameServerSets {
		marshaledGameServer, err := json.Marshal(gameServerSet)
		if err != nil {
			return nil, err
		}
		marshaledGameServerSets = append(marshaledGameServerSets, string(marshaledGameServer))
	}

	return marshaledGameServerSets, nil
}

func (g GameServerSetGrpcService) UpdateGameServerSets(ctx context.Context, request *api.UpdateGameServerSetsRequest) (*api.UpdateGameServerSetsResponse, error) {
	results, err := service.GetGssService().UpdateGameServerSets(&apimodels.UpdateGameServerSetsRequest{
		Filter:    request.Filter,
		JsonPatch: request.JsonPatch,
	})
	if err != nil {
		slog.Error("update GameServerSets failed", "error", err)
		return nil, err
	}

	updateGameServerSetsResults, err := toGrpcUpdateGameServerSetsResults(results)
	if err != nil {
		slog.Error("transform UpdateGameServerSetsResults failed", "error", err)
		return nil, err
	}

	return &api.UpdateGameServerSetsResponse{
		Results: updateGameServerSetsResults,
	}, nil
}

func toGrpcUpdateGameServerSetsResults(updateGameServerSetResults []updater.UpdateGssResult) ([]*api.UpdateGameServerSetResult, error) {
	results := make([]*api.UpdateGameServerSetResult, 0, len(updateGameServerSetResults))
	for _, updateGameServerSetResult := range updateGameServerSetResults {
		marshaledGameServerSet, err := json.Marshal(updateGameServerSetResult.Gss)
		if err != nil {
			return nil, err
		}
		marshaledUpdatedGameServerSet, err := json.Marshal(updateGameServerSetResult.UpdatedGss)
		if err != nil {
			return nil, err
		}
		result := &api.UpdateGameServerSetResult{
			GameServerSet:        string(marshaledGameServerSet),
			UpdatedGameServerSet: string(marshaledUpdatedGameServerSet),
		}
		if updateGameServerSetResult.Err != nil {
			*result.Error = updateGameServerSetResult.Err.Error()
		}
		results = append(results, result)
	}

	return results, nil
}

func (g GameServerSetGrpcService) DeleteGameServerSets(ctx context.Context, request *api.DeleteGameServerSetsRequest) (*api.DeleteGameServerSetsResponse, error) {
	results, err := service.GetGssService().DeleteGameServerSets(request.Filter)
	if err != nil {
		slog.Error("delete GameServerSets failed", "error", err)
		return nil, err
	}

	deleteGameServerSetsResults, err := toGrpcDeleteGameServerSetsResults(results)
	if err != nil {
		slog.Error("transform DeleteGameServerSetsResults failed", "error", err)
		return nil, err
	}

	return &api.DeleteGameServerSetsResponse{
		Results: deleteGameServerSetsResults,
	}, nil
}

func toGrpcDeleteGameServerSetsResults(deleteGameServerSetResults []deleter.DeleteGssResult) ([]*api.DeleteGameServerSetResult, error) {
	results := make([]*api.DeleteGameServerSetResult, 0, len(deleteGameServerSetResults))
	for _, deleteGameServerSetResult := range deleteGameServerSetResults {
		marshaledGameServerSet, err := json.Marshal(deleteGameServerSetResult.Gss)
		if err != nil {
			return nil, err
		}
		result := &api.DeleteGameServerSetResult{
			GameServerSet: string(marshaledGameServerSet),
		}
		if deleteGameServerSetResult.Err != nil {
			*result.Error = deleteGameServerSetResult.Err.Error()
		}
		results = append(results, result)
	}

	return results, nil
}
