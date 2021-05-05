package apiserver

import (
	"apigateway/api/apigateway"
	api "apigateway/api/playlistservice"
	"golang.org/x/net/context"
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

func (server *apiGatewayServer) GetUserPlaylists(ctx context.Context, req *apigateway.GetUserPlaylistsRequest) (*apigateway.GetUserPlaylistsResponse, error) {
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
