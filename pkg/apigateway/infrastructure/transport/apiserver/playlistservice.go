package apiserver

import (
	"apigateway/api/apigateway"
	api "apigateway/api/playlistservice"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *apiGatewayServer) CreatePlaylist(ctx context.Context, req *apigateway.CreatePlaylistRequest) (*apigateway.CreatePlaylistResponse, error) {
	userToken, err := server.authenticateUser(ctx)
	if err != nil {
		return nil, err
	}

	serializedToken, err := server.userDescriptorSerializer.Serialize(userToken)
	if err != nil {
		return nil, err
	}

	resp, err := server.playlistServiceClient.CreatePlaylist(ctx, &api.CreatePlaylistRequest{
		Name:      req.Name,
		UserToken: serializedToken,
	})
	if err != nil {
		return nil, err
	}

	return &apigateway.CreatePlaylistResponse{
		PlaylistID: resp.PlaylistID,
	}, nil
}

func (server *apiGatewayServer) GetPlaylist(ctx context.Context, req *apigateway.GetPlaylistRequest) (*apigateway.GetPlaylistResponse, error) {
	userToken, err := server.authenticateUser(ctx)
	if err != nil {
		return nil, err
	}

	serializedToken, err := server.userDescriptorSerializer.Serialize(userToken)
	if err != nil {
		return nil, err
	}

	resp, err := server.playlistServiceClient.GetPlaylist(ctx, &api.GetPlaylistRequest{
		PlaylistID: req.PlaylistID,
		UserToken:  serializedToken,
	})
	if err != nil {
		return nil, err
	}

	return &apigateway.GetPlaylistResponse{
		Name:               resp.Name,
		OwnerID:            resp.OwnerID,
		CreatedAtTimestamp: resp.CreatedAtTimestamp,
		UpdatedAtTimestamp: resp.UpdatedAtTimestamp,
		PlaylistItems:      convertToPlaylistItemsApiGateway(resp.PlaylistItems),
	}, nil
}

func (server *apiGatewayServer) GetUserPlaylists(ctx context.Context, _ *apigateway.GetUserPlaylistsRequest) (*apigateway.GetUserPlaylistsResponse, error) {
	userToken, err := server.authenticateUser(ctx)
	if err != nil {
		return nil, err
	}

	serializedToken, err := server.userDescriptorSerializer.Serialize(userToken)
	if err != nil {
		return nil, err
	}

	resp, err := server.playlistServiceClient.GetUserPlaylists(ctx, &api.GetUserPlaylistsRequest{
		UserToken: serializedToken,
	})
	if err != nil {
		return nil, err
	}

	return &apigateway.GetUserPlaylistsResponse{
		Playlists: convertToPlaylistsApiGateway(resp.Playlists),
	}, nil
}

func (server *apiGatewayServer) AddToPlaylist(ctx context.Context, req *apigateway.AddToPlaylistRequest) (*apigateway.AddToPlaylistResponse, error) {
	userToken, err := server.authenticateUser(ctx)
	if err != nil {
		return nil, err
	}

	serializedToken, err := server.userDescriptorSerializer.Serialize(userToken)
	if err != nil {
		return nil, err
	}

	resp, err := server.playlistServiceClient.AddToPlaylist(ctx, &api.AddToPlaylistRequest{
		PlaylistID: req.PlaylistID,
		UserToken:  serializedToken,
		ContentID:  req.ContentID,
	})
	if err != nil {
		return nil, err
	}

	return &apigateway.AddToPlaylistResponse{
		PlaylistItemID: resp.PlaylistItemID,
	}, nil
}

func (server *apiGatewayServer) SetPlaylistName(ctx context.Context, req *apigateway.SetPlaylistNameRequest) (*emptypb.Empty, error) {
	userToken, err := server.authenticateUser(ctx)
	if err != nil {
		return nil, err
	}

	serializedToken, err := server.userDescriptorSerializer.Serialize(userToken)
	if err != nil {
		return nil, err
	}

	_, err = server.playlistServiceClient.SetPlaylistName(ctx, &api.SetPlaylistNameRequest{
		PlaylistID: req.PlaylistID,
		NewName:    req.NewName,
		UserToken:  serializedToken,
	})

	return &emptypb.Empty{}, err
}

func (server *apiGatewayServer) RemoveFromPlaylist(ctx context.Context, req *apigateway.RemoveFromPlaylistRequest) (*emptypb.Empty, error) {
	userToken, err := server.authenticateUser(ctx)
	if err != nil {
		return nil, err
	}

	serializedToken, err := server.userDescriptorSerializer.Serialize(userToken)
	if err != nil {
		return nil, err
	}

	_, err = server.playlistServiceClient.RemoveFromPlaylist(ctx, &api.RemoveFromPlaylistRequest{
		PlaylistItemID: req.PlaylistItemID,
		UserToken:      serializedToken,
	})

	return &emptypb.Empty{}, err
}

func (server *apiGatewayServer) RemovePlaylist(ctx context.Context, req *apigateway.RemovePlaylistRequest) (*emptypb.Empty, error) {
	userToken, err := server.authenticateUser(ctx)
	if err != nil {
		return nil, err
	}

	serializedToken, err := server.userDescriptorSerializer.Serialize(userToken)
	if err != nil {
		return nil, err
	}

	_, err = server.playlistServiceClient.RemovePlaylist(ctx, &api.RemovePlaylistRequest{
		PlaylistID: req.PlaylistID,
		UserToken:  serializedToken,
	})

	return &emptypb.Empty{}, err
}

func convertToPlaylistsApiGateway(playlists []*api.Playlist) []*apigateway.Playlist {
	result := make([]*apigateway.Playlist, len(playlists))
	for i, playlist := range playlists {
		result[i] = convertToPlaylistApiGateway(playlist)
	}
	return result
}

func convertToPlaylistApiGateway(playlist *api.Playlist) *apigateway.Playlist {
	return &apigateway.Playlist{
		PlaylistID:         playlist.PlaylistID,
		Name:               playlist.Name,
		OwnerID:            playlist.OwnerID,
		CreatedAtTimestamp: playlist.CreatedAtTimestamp,
		UpdatedAtTimestamp: playlist.UpdatedAtTimestamp,
		PlaylistItems:      convertToPlaylistItemsApiGateway(playlist.PlaylistItems),
	}
}

func convertToPlaylistItemsApiGateway(playlistItems []*api.PlaylistItem) []*apigateway.PlaylistItem {
	result := make([]*apigateway.PlaylistItem, len(playlistItems))
	for i, item := range playlistItems {
		result[i] = &apigateway.PlaylistItem{
			PlaylistItemID:     item.PlaylistItemID,
			ContentID:          item.ContentID,
			CreatedAtTimestamp: item.CreatedAtTimestamp,
		}
	}
	return result
}
