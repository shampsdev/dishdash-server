package framework

import (
	"bytes"
	"encoding/json"
	"fmt"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase"
)

func (fw *Framework) postUserWithID(user *domain.User) (*domain.User, error) {
	b, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user: %w", err)
	}

	resp, err := fw.HttpCli.Post(fmt.Sprintf("%s/users/with_id", fw.ApiHost), "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed to post user: %w", err)
	}

	respUser := &domain.User{}
	err = json.NewDecoder(resp.Body).Decode(respUser)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return respUser, nil
}

func (fw *Framework) MustFindLobby() *domain.Lobby {
	lobby, err := fw.FindLobby()
	if err != nil {
		panic(err)
	}
	return lobby
}

func (fw *Framework) FindLobby() (*domain.Lobby, error) {
	findLobbyInput := usecase.FindLobbyInput{
		Dist: 0,
		// ИТМО - Кронверкский проспект, 49
		Location: domain.Coordinate{Lon: 30.310011, Lat: 59.956363},
	}
	b, err := json.Marshal(findLobbyInput)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal find lobby input: %w", err)
	}

	resp, err := fw.HttpCli.Post(fmt.Sprintf("%s/lobbies/find", fw.ApiHost), "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed to post find lobby: %w", err)
	}

	lobby := &domain.Lobby{}
	err = json.NewDecoder(resp.Body).Decode(lobby)
	if err != nil {
		return nil, fmt.Errorf("failed to decode lobby: %w", err)
	}
	return lobby, nil
}
