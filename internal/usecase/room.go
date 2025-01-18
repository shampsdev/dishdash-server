package usecase

import (
	"context"
	"fmt"
	"slices"
	"sync"

	log "github.com/sirupsen/logrus"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase/state"
	"dishdash.ru/pkg/algo"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type State string

type VoteType string

const (
	VoteTypeMatch  VoteType = "match"
	VoteTypeFinish VoteType = "finish"
)

type Match struct {
	ID    int           `json:"id"`
	Place *domain.Place `json:"card"`
}

type OptionID int64

type VoteOption struct {
	ID   OptionID `json:"id"`
	Desc string   `json:"description"`
}

type Vote interface {
	isVote()
}

type BaseVote struct {
	ID      int64        `json:"id"`
	Options []VoteOption `json:"options"`
	Type    VoteType     `json:"type"`
	votes   map[string]OptionID
}

type FinishVote struct {
	BaseVote
}

func (v *FinishVote) isVote() {}

type MatchVote struct {
	BaseVote
	Place *domain.Place `json:"card"`
}

func (v *MatchVote) isVote() {}

const (
	OptionIDLike OptionID = iota
	OptionIDDislike

	OptionIDFinish
	OptionIDContinue
)

type VoteResult struct {
	Type     VoteType `json:"type"`
	VoteID   int64    `json:"voteId"`
	OptionID OptionID `json:"optionId"`
}

type Room struct {
	id   string
	lock sync.RWMutex

	lobby *domain.Lobby

	recommendationOpts *domain.RecommendationOpts

	state      domain.LobbyState
	usersMap   map[string]*domain.User
	usersPlace map[string]*domain.Place
	places     []*domain.Place
	swipes     []*domain.Swipe
	matches    []*Match
	result     *domain.Place

	votes map[int64]Vote

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
		id:               lobby.ID,
		lobby:            lobby,
		usersMap:         make(map[string]*domain.User),
		places:           make([]*domain.Place, 0),
		usersPlace:       make(map[string]*domain.Place),
		swipes:           make([]*domain.Swipe, 0),
		matches:          make([]*Match, 0),
		votes:            make(map[int64]Vote),
		lobbyUseCase:     lobbyUseCase,
		placeUseCase:     placeUseCase,
		swipeUseCase:     swipeUseCase,
		userUseCase:      userUseCase,
		placeRecommender: placeRecommender,
		log:              log.WithFields(log.Fields{"room": lobby.ID}),
	}

	r.votes[0] = &FinishVote{
		BaseVote{
			ID:      0,
			Options: []VoteOption{{ID: OptionIDFinish, Desc: "Finish"}, {ID: OptionIDContinue, Desc: "Continue"}},
			Type:    VoteTypeFinish,
			votes:   make(map[string]OptionID),
		},
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

func (r *Room) ID() string {
	return r.id
}

func (r *Room) loadFromDB() error {
	var err error
	r.swipes, err = r.swipeUseCase.GetSwipesByLobbyID(context.Background(), r.id)
	if err != nil {
		return fmt.Errorf("failed to get swipes: %w", err)
	}

	users, err := r.userUseCase.GetUsersByLobbyID(context.Background(), r.id)
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
		if e.Value == len(r.usersMap) {
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
		ID:       r.id,
		PriceAvg: r.lobby.PriceAvg,
		Location: r.lobby.Location,
		Tags: algo.Map(r.lobby.Tags, func(t *domain.Tag) int64 {
			return t.ID
		}),
	}
}

func (r *Room) OnJoin(c *state.Context[*Room]) error {
	r.usersMap[c.User.ID] = c.User

	// c.BroadcastToOthers(event.UserJoined, event.UserJoinedEvent{
	// 	ID:     c.User.ID,
	// 	Name:   c.User.Name,
	// 	Avatar: c.User.Avatar,
	// })

	return nil
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
		algo.Map(r.users(), func(u *domain.User) string {
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
	recommendationOpts *domain.RecommendationOpts,
) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.updateLobbySettings(ctx, location, priceAvg, tagIDs, placeIDs, recommendationOpts)
}

func (r *Room) updateLobbySettings(
	ctx context.Context,
	location domain.Coordinate,
	priceAvg int,
	tagIDs, placeIDs []int64,
	recommendationOpts *domain.RecommendationOpts,
) error {
	r.recommendationOpts = recommendationOpts
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
		err := r.updateLobbySettings(ctx, r.lobby.Location, 500, []int64{3}, nil, r.recommendationOpts)
		if err != nil {
			r.log.WithError(err).Error("Action 'UpdateLobby' failed")
			return err
		}
		r.log.Debug("Request successful: Action 'UpdateLobby' completed with default tags")
	}

	var err error
	r.places, err = r.placeRecommender.RecommendPlaces(ctx,
		r.recommendationOpts,
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
		algo.Map(r.lobby.Tags, func(t *domain.Tag) int64 {
			return t.ID
		}),
		algo.Map(r.places, func(p *domain.Place) int64 {
			return p.ID
		}),
		r.recommendationOpts,
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

func (r *Room) Swipe(userID string, placeID int64, t domain.SwipeType) (*MatchVote, error) {
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

	likes := algo.Count(r.swipes, func(swipe *domain.Swipe) bool {
		return swipe.PlaceID == placeID && swipe.Type == domain.LIKE
	})
	if likes == len(r.usersMap) {
		match := &Match{Place: r.places[slices.IndexFunc(r.places, func(place *domain.Place) bool {
			return place.ID == placeID
		})]}
		match.ID = len(r.matches)
		r.matches = append(r.matches, match)
		r.result = match.Place

		vote := &MatchVote{
			BaseVote: BaseVote{
				ID: int64(len(r.votes)),
				Options: []VoteOption{
					{ID: OptionIDLike, Desc: "Like"},
					{ID: OptionIDDislike, Desc: "Dislike"},
				},
				Type:  VoteTypeMatch,
				votes: make(map[string]OptionID),
			},
			Place: match.Place,
		}
		r.votes[vote.ID] = vote
		return vote, nil
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
	if r.result != nil {
		return r.result
	}
	if len(r.matches) > 0 {
		return r.matches[len(r.matches)-1].Place
	}
	return nil
}

func (r *Room) Matches() []*Match {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.matches
}

func (r *Room) Votes() []Vote {
	r.lock.RLock()
	defer r.lock.RUnlock()
	votes := make([]Vote, 0, len(r.votes))
	for _, vote := range r.votes {
		votes = append(votes, vote)
	}
	return votes
}

func (r *Room) Vote(userID string, voteID int64, optionID OptionID) (*VoteResult, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	vote, ok := r.votes[voteID]
	if !ok {
		return nil, fmt.Errorf("vote with id %d not found", voteID)
	}

	switch v := vote.(type) {
	case *MatchVote:
		if optionID != OptionIDLike && optionID != OptionIDDislike {
			return nil, fmt.Errorf("invalid option id: %d", optionID)
		}
		v.votes[userID] = optionID

		log.Debugf("<User %s> voted for <Option %d> in %s", userID, optionID, v.Type)
		res := r.voteMatchResult(v)

		if res != nil && res.OptionID == OptionIDLike {
			r.result = v.Place
			err := r.setState(domain.Finished)
			if err != nil {
				return nil, fmt.Errorf("error while setting lobby state: %w", err)
			}
			return res, nil
		}

		return res, nil
	case *FinishVote:
		if optionID != OptionIDFinish && optionID != OptionIDContinue {
			return nil, fmt.Errorf("invalid option id: %d", optionID)
		}
		v.votes[userID] = optionID

		log.Debugf("<User %s> voted for <Option %d> in %s", userID, optionID, v.Type)
		res := r.voteFinishResult(v)
		if res != nil && res.OptionID == OptionIDFinish {
			err := r.setState(domain.Finished)
			if err != nil {
				return nil, fmt.Errorf("error while setting lobby state: %w", err)
			}
			return res, nil
		}

		return res, nil
	}

	log.Warnf("Unknown voting type: %T", vote)
	return nil, fmt.Errorf("unknown voting type: %T", vote)
}

func (r *Room) voteMatchResult(vote *MatchVote) *VoteResult {
	if len(vote.votes) != len(r.usersMap) {
		return nil
	}
	likes := 0
	for _, vote := range vote.votes {
		if vote == OptionIDLike {
			likes++
		}
	}
	if likes == len(r.usersMap) {
		return &VoteResult{
			VoteID:   vote.ID,
			Type:     VoteTypeMatch,
			OptionID: OptionIDLike,
		}
	}

	return &VoteResult{
		VoteID:   vote.ID,
		Type:     VoteTypeMatch,
		OptionID: OptionIDDislike,
	}
}

func (r *Room) voteFinishResult(vote *FinishVote) *VoteResult {
	if len(vote.votes) != len(r.usersMap) {
		return nil
	}
	likes := 0
	for _, vote := range vote.votes {
		if vote == OptionIDFinish {
			likes++
		}
	}
	if likes == len(r.usersMap) {
		return &VoteResult{
			VoteID:   vote.ID,
			Type:     VoteTypeFinish,
			OptionID: OptionIDFinish,
		}
	}
	return nil
}
