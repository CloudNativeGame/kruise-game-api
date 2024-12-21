package service

import (
	"context"
	"encoding/json"
	"github.com/CloudNativeGame/kruise-game-api/facade/apiserver/apimodels"
	api "github.com/CloudNativeGame/kruise-game-api/facade/apiserver/proto"
	"github.com/CloudNativeGame/kruise-game-api/facade/apiserver/service"
	"github.com/CloudNativeGame/kruise-game-api/pkg/deleter"
	"github.com/CloudNativeGame/kruise-game-api/pkg/updater"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
	"log/slog"
)

type GameServerGrpcService struct {
	api.UnimplementedGameServerServiceServer
}

func (g *GameServerGrpcService) GetGameServers(ctx context.Context, request *api.GetGameServersRequest) (*api.GetGameServersResponse, error) {
	gameServers, err := service.GetGsService().GetGameServers(request.Filter)
	if err != nil {
		slog.Error("get filtered GameServers failed", "error", err)
		return nil, err
	}

	marshaledGameServers, err := marshalGameServers(gameServers)
	if err != nil {
		slog.Error("marshal GameServers failed", "error", err)
		return nil, err
	}

	return &api.GetGameServersResponse{
		GameServers: marshaledGameServers,
	}, nil
}

func marshalGameServers(gameServers []*v1alpha1.GameServer) ([]string, error) {
	marshaledGameServers := make([]string, 0, len(gameServers))
	for _, gameServer := range gameServers {
		marshaledGameServer, err := json.Marshal(gameServer)
		if err != nil {
			return nil, err
		}
		marshaledGameServers = append(marshaledGameServers, string(marshaledGameServer))
	}

	return marshaledGameServers, nil
}

func (g *GameServerGrpcService) UpdateGameServers(ctx context.Context, request *api.UpdateGameServersRequest) (*api.UpdateGameServersResponse, error) {
	results, err := service.GetGsService().UpdateGameServers(&apimodels.UpdateGameServersRequest{
		Filter:    request.Filter,
		JsonPatch: request.JsonPatch,
	})
	if err != nil {
		slog.Error("update GameServers failed", "error", err)
		return nil, err
	}

	updateGameServersResults, err := toGrpcUpdateGameServersResults(results)
	if err != nil {
		slog.Error("transform UpdateGameServersResults failed", "error", err)
		return nil, err
	}

	return &api.UpdateGameServersResponse{
		Results: updateGameServersResults,
	}, nil
}

func toGrpcUpdateGameServersResults(updateGameServerResults []updater.UpdateGsResult) ([]*api.UpdateGameServerResult, error) {
	results := make([]*api.UpdateGameServerResult, 0, len(updateGameServerResults))
	for _, updateGameServerResult := range updateGameServerResults {
		marshaledGameServer, err := json.Marshal(updateGameServerResult.Gs)
		if err != nil {
			return nil, err
		}
		marshaledUpdatedGameServer, err := json.Marshal(updateGameServerResult.UpdatedGs)
		if err != nil {
			return nil, err
		}
		result := &api.UpdateGameServerResult{
			GameServer:        string(marshaledGameServer),
			UpdatedGameServer: string(marshaledUpdatedGameServer),
		}
		if updateGameServerResult.Err != nil {
			*result.Error = updateGameServerResult.Err.Error()
		}
		results = append(results, result)
	}

	return results, nil
}

func (g *GameServerGrpcService) DeleteGameServers(ctx context.Context, request *api.DeleteGameServersRequest) (*api.DeleteGameServersResponse, error) {
	results, err := service.GetGsService().DeleteGameServers(request.Filter)
	if err != nil {
		slog.Error("delete GameServers failed", "error", err)
		return nil, err
	}

	deleteGameServersResults, err := toGrpcDeleteGameServersResults(results)
	if err != nil {
		slog.Error("transform DeleteGameServersResults failed", "error", err)
		return nil, err
	}

	return &api.DeleteGameServersResponse{
		Results: deleteGameServersResults,
	}, nil
}

func toGrpcDeleteGameServersResults(deleteGameServerResults []deleter.DeleteGsResult) ([]*api.DeleteGameServerResult, error) {
	results := make([]*api.DeleteGameServerResult, 0, len(deleteGameServerResults))
	for _, deleteGameServerResult := range deleteGameServerResults {
		marshaledGameServer, err := json.Marshal(deleteGameServerResult.Gs)
		if err != nil {
			return nil, err
		}
		result := &api.DeleteGameServerResult{
			GameServer: string(marshaledGameServer),
		}
		if deleteGameServerResult.Err != nil {
			*result.Error = deleteGameServerResult.Err.Error()
		}
		results = append(results, result)
	}

	return results, nil
}
