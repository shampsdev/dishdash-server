package usecase

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"

	log "github.com/sirupsen/logrus"

	"dishdash.ru/internal/domain"
	"dishdash.ru/pkg/filter"
)

type State string

type VoteOption int

var (
	VoteDislike VoteOption
	VoteLike    VoteOption = 1
)

type Match struct {
	ID    int
	Place *domain.Place
}

type Room struct {
	ID   string
	lock sync.RWMutex

	lobby *domain.Lobby

	state      domain.LobbyState
	usersMap   map[string]*domain.User
	usersPlace map[string]*domain.Place
	places     []*domain.Place
	swipes     []*domain.Swipe
	matches    []*Match
	matchVotes map[string]VoteOption

	lobbyUseCase     Lobby
	placeUseCase     Place
	placeRecommender *PlaceRecommender
	log              *log.Entry
}

func NewRoom(
	lobby *domain.Lobby,
	lobbyUseCase Lobby,
	placeUseCase Place,
	placeRecommender *PlaceRecommender,
) (*Room, error) {
	r := &Room{
		ID:               lobby.ID,
		lobby:            lobby,
		usersMap:         make(map[string]*domain.User),
		places:           make([]*domain.Place, 0),
		usersPlace:       make(map[string]*domain.Place),
		swipes:           make([]*domain.Swipe, 0),
		matchVotes:       make(map[string]VoteOption),
		lobbyUseCase:     lobbyUseCase,
		placeUseCase:     placeUseCase,
		placeRecommender: placeRecommender,
		log:              log.WithFields(log.Fields{"room": lobby.ID}),
	}
	r.state = lobby.State

	if lobby.State != domain.InLobby {
		return nil, errors.New("can't connect to started lobby")
	}

	return r, nil
}

func (r *Room) GetNextPlaceForUser(id string) *domain.Place {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.usersPlace[id]
}

func (r *Room) Users() []*domain.User {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.users()
}

func (r *Room) users() []*domain.User {
	users := make([]*domain.User, 0)
	for _, user := range r.usersMap {
		users = append(users, user)
	}
	return users
}

func (r *Room) Settings() UpdateLobbySettingsInput {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.settings()
}

func (r *Room) settings() UpdateLobbySettingsInput {
	return UpdateLobbySettingsInput{
		ID:       r.ID,
		PriceAvg: r.lobby.PriceAvg,
		Location: r.lobby.Location,
		Tags: filter.Map(r.lobby.Tags, func(t *domain.Tag) int64 {
			return t.ID
		}),
	}
}

func (r *Room) AddUser(user *domain.User) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, has := r.usersMap[user.ID]; has {
		return fmt.Errorf("user %s already exists", user.ID)
	}
	r.usersMap[user.ID] = user

	if r.state == domain.Swiping {
		r.usersPlace[user.ID] = r.places[0]
	}

	return r.syncUsersWithBd()
}

func (r *Room) syncUsersWithBd() error {
	var err error
	r.lobby.Users, err = r.lobbyUseCase.SetLobbyUsers(
		context.Background(),
		r.lobby.ID,
		filter.Map(r.users(), func(u *domain.User) string {
			return u.ID
		}),
	)
	if err != nil {
		return fmt.Errorf("error while setting lobby users: %w", err)
	}
	return nil
}

func (r *Room) RemoveUser(id string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, has := r.usersMap[id]
	if !has {
		return fmt.Errorf("user %s not found", id)
	}
	delete(r.usersMap, id)

	if r.state == domain.Swiping {
		err := r.syncUsersWithBd()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Room) Empty() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return len(r.usersMap) == 0
}

func (r *Room) InLobby() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.state == domain.InLobby
}

func (r *Room) Voting() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.state == domain.Voting
}

func (r *Room) Swiping() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.state == domain.Swiping
}

func (r *Room) Finished() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.state == domain.Finished
}

func (r *Room) UpdateLobbySettings(
	ctx context.Context,
	priceAvg int,
	tagIDs, placeIDs []int64,
) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.updateLobbySettings(ctx, priceAvg, tagIDs, placeIDs)
}

func (r *Room) updateLobbySettings(
	ctx context.Context,
	priceAvg int,
	tagIDs, placeIDs []int64,
) error {
	lobby, err := r.lobbyUseCase.SetLobbySettings(ctx, UpdateLobbySettingsInput{
		ID:       r.lobby.ID,
		PriceAvg: priceAvg,
		Location: r.lobby.Location,
		Tags:     tagIDs,
		Places:   placeIDs,
	})
	if err != nil {
		return err
	}
	r.lobby.PriceAvg = lobby.PriceAvg
	r.lobby.Tags = lobby.Tags
	r.lobby.Places = lobby.Places
	return nil
}

func (r *Room) StartSwipes(ctx context.Context) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.log.Debug("Request started: Action 'StartSwipes' initiated")

	if len(r.lobby.Tags) == 0 {
		r.log.Warn("Action 'StartSwipes' encountered an issue. Reason: 'No tags found, using default tags'")
		err := r.updateLobbySettings(ctx, 500, []int64{3}, nil)
		if err != nil {
			r.log.WithError(err).Error("Action 'UpdateLobby' failed")
			return err
		}
		r.log.Debug("Request successful: Action 'UpdateLobby' completed with default tags")
	}

	var err error
	r.places, err = r.placeRecommender.RecommendPlaces(ctx,
		domain.RecommendData{
			Location: r.lobby.Location,
			PriceAvg: r.lobby.PriceAvg,
			Tags:     r.lobby.TagNames(),
		},
	)
	if err != nil {
		r.log.WithError(err).Error("Action 'GetPlacesForLobby' failed")
		return err
	}
	r.log.Debug("Request successful: Action 'GetPlacesForLobby' completed")

	err = r.updateLobbySettings(ctx, r.lobby.PriceAvg,
		filter.Map(r.lobby.Tags, func(t *domain.Tag) int64 {
			return t.ID
		}),
		filter.Map(r.places, func(p *domain.Place) int64 {
			return p.ID
		}),
	)
	if err != nil {
		r.log.WithError(err).Error("Action 'UpdateLobby' failed")
		return err
	}
	r.log.Debug("Request successful: Action 'UpdateLobby' completed")

	r.log.Debug("Request successful: Action 'UpdateLobby' completed")
	for id := range r.usersMap {
		if len(r.places) > 0 {
			r.usersPlace[id] = r.places[0]
			r.log.Debugf("<User %s> is assigned to <Place %d>", id, r.places[0].ID)
		} else {
			log.Warnf("No places available to assign to <User %s>", id)
		}
	}

	log.Info("Swipes successfully started")
	err = r.setState(domain.Swiping)
	return err
}

func (r *Room) Swipe(userID string, placeID int64, t domain.SwipeType) (*Match, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.swipes = append(r.swipes, &domain.Swipe{
		LobbyID: r.lobby.ID,
		PlaceID: placeID,
		UserID:  userID,
		Type:    t,
	})

	pIdx := slices.IndexFunc(r.places, func(place *domain.Place) bool {
		return place.ID == placeID
	})
	r.usersPlace[userID] = r.places[(pIdx+1)%len(r.places)]

	likes := filter.Count(r.swipes, func(swipe *domain.Swipe) bool {
		return swipe.PlaceID == placeID && swipe.Type == domain.LIKE
	})
	if likes > len(r.usersMap)/2 {
		match := &Match{Place: r.places[slices.IndexFunc(r.places, func(place *domain.Place) bool {
			return place.ID == placeID
		})]}
		match.ID = len(r.matches)
		r.matches = append(r.matches, match)
		if err := r.setState(domain.Voting); err != nil {
			return nil, err
		}
		return match, nil
	}

	return nil, nil
}

func (r *Room) setState(state domain.LobbyState) error {
	r.state = state
	err := r.lobbyUseCase.SetLobbyState(context.Background(), r.lobby.ID, state)
	return err
}

func (r *Room) Result() *domain.Place {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.matches[len(r.matches)-1].Place
}

func (r *Room) Vote(userID string, option VoteOption) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.matchVotes[userID] = option
	if len(r.matchVotes) == len(r.usersMap) {
		if r.allVotedLike() {
			if err := r.setState(domain.Finished); err != nil {
				return err
			}
		} else {
			if err := r.setState(domain.Swiping); err != nil {
				return err
			}
		}
		r.matchVotes = make(map[string]VoteOption)
	}
	return nil
}

func (r *Room) allVotedLike() bool {
	for _, vote := range r.matchVotes {
		if vote == VoteDislike {
			return false
		}
	}
	return true
}
