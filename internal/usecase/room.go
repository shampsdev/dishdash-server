package usecase

import (
	"context"
	"fmt"
	"slices"
	"sync"

	log "github.com/sirupsen/logrus"

	"dishdash.ru/internal/domain"
	"dishdash.ru/pkg/filter"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type State string

type VoteOption int

var (
	VoteDislike VoteOption
	VoteLike    VoteOption = 1
)

type Match struct {
	ID    int           `json:"id"`
	Place *domain.Place `json:"card"`
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
	swipeUseCase     Swipe
	userUseCase      User
	placeRecommender *PlaceRecommender
	log              *log.Entry
}

func NewRoom(
	lobby *domain.Lobby,
	lobbyUseCase Lobby,
	placeUseCase Place,
	swipeUseCase Swipe,
	userUseCase User,
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
		swipeUseCase:     swipeUseCase,
		userUseCase:      userUseCase,
		placeRecommender: placeRecommender,
		log:              log.WithFields(log.Fields{"room": lobby.ID}),
	}
	r.state = lobby.State
	r.log.Debugf("state: %s", r.state)

	if r.state != domain.InLobby {
		// can't load swiping lobby
		err := r.setState(domain.Finished)
		if err != nil {
			return nil, fmt.Errorf("failed to set state: %w", err)
		}
	}

	if r.state == domain.Finished || r.state == domain.InLobby {
		err := r.loadFromDB()
		if err != nil {
			return nil, fmt.Errorf("failed to load from db: %w", err)
		}
	}

	return r, nil
}

func (r *Room) loadFromDB() error {
	var err error
	r.swipes, err = r.swipeUseCase.GetSwipesByLobbyID(context.Background(), r.ID)
	if err != nil {
		return fmt.Errorf("failed to get swipes: %w", err)
	}

	users, err := r.userUseCase.GetUsersByLobbyID(context.Background(), r.ID)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}
	for _, user := range users {
		r.usersMap[user.ID] = user
	}

	r.matches, err = r.mathesFromSwipes()
	if err != nil {
		return fmt.Errorf("failed to get matches: %w", err)
	}

	return nil
}

func (r *Room) mathesFromSwipes() ([]*Match, error) {
	r.log.WithFields(log.Fields{"swipeCount": len(r.swipes)}).Debug("restore matches from db")
	swipeCount := orderedmap.New[int64, int]()

	for _, swipe := range r.swipes {
		if swipe.Type != domain.LIKE {
			continue
		}
		c, ok := swipeCount.Get(swipe.PlaceID)
		if !ok {
			swipeCount.Set(swipe.PlaceID, 0)
		}
		swipeCount.Set(swipe.PlaceID, c+1)
	}

	matches := make([]*Match, 0)

	for e := swipeCount.Oldest(); e != nil; e = e.Next() {
		r.log.WithFields(log.Fields{"place": e.Key, "count": e.Value}).Debug("swipe count")
		if e.Value > len(r.usersMap)/2 {
			place, err := r.placeUseCase.GetPlaceByID(context.Background(), e.Key)
			if err != nil {
				return nil, fmt.Errorf("failed to get place: %w", err)
			}
			matches = append(matches, &Match{
				ID:    len(matches),
				Place: place,
			})
		}
	}

	return matches, nil
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

	if r.state != domain.InLobby {
		return fmt.Errorf("can't add user to lobby in state %s", r.state)
	}

	if _, has := r.usersMap[user.ID]; has {
		return nil
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
		return nil
	}
	delete(r.usersMap, id)

	if r.state == domain.InLobby {
		err := r.syncUsersWithBd()
		if err != nil {
			return fmt.Errorf("error while syncing users with bd: %w", err)
		}
	}

	return nil
}

func (r *Room) Empty() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return len(r.usersMap) == 0
}

func (r *Room) State() domain.LobbyState {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.state
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
	location domain.Coordinate,
	priceAvg int,
	tagIDs, placeIDs []int64,
) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.updateLobbySettings(ctx, location, priceAvg, tagIDs, placeIDs)
}

func (r *Room) updateLobbySettings(
	ctx context.Context,
	location domain.Coordinate,
	priceAvg int,
	tagIDs, placeIDs []int64,
) error {
	lobby, err := r.lobbyUseCase.SetLobbySettings(ctx, UpdateLobbySettingsInput{
		ID:       r.lobby.ID,
		PriceAvg: priceAvg,
		Location: location,
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
		err := r.updateLobbySettings(ctx, r.lobby.Location, 500, []int64{3}, nil)
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

	err = r.updateLobbySettings(ctx, r.lobby.Location, r.lobby.PriceAvg,
		filter.Map(r.lobby.Tags, func(t *domain.Tag) int64 {
			return t.ID
		}),
		filter.Map(r.places, func(p *domain.Place) int64 {
			return p.ID
		}),
	)
	if err != nil {
		return fmt.Errorf("error while updating lobby settings: %w", err)
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
	err := r.swipeUseCase.SaveSwipe(context.Background(), &domain.Swipe{
		LobbyID: r.lobby.ID,
		PlaceID: placeID,
		UserID:  userID,
		Type:    t,
	})
	if err != nil {
		return nil, fmt.Errorf("error while saving swipe: %w", err)
	}

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

func (r *Room) Matches() []*Match {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.matches
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
