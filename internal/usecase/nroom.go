package usecase

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase/nevent"
	"dishdash.ru/internal/usecase/state"
	algo "dishdash.ru/pkg/algo"
	log "github.com/sirupsen/logrus"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type NRoom struct {
	id   string
	lock sync.RWMutex

	lobby *domain.Lobby

	recommendationOpts *domain.RecommendationOpts

	state      domain.LobbyState
	usersMap   map[string]*domain.User
	usersPlace map[string]*domain.Place
	places     []*domain.Place
	swipes     []*domain.Swipe
	matches    []*nevent.Match
	result     *domain.Place

	votes map[int64]nevent.VoteAnnounce

	lobbyUseCase     Lobby
	placeUseCase     Place
	swipeUseCase     Swipe
	userUseCase      User
	placeRecommender *PlaceRecommender
	log              *log.Entry
}

func NewNRoom(
	lobby *domain.Lobby,
	lobbyUseCase Lobby,
	placeUseCase Place,
	swipeUseCase Swipe,
	userUseCase User,
	placeRecommender *PlaceRecommender,
) (*NRoom, error) {
	r := &NRoom{
		id:               lobby.ID,
		lobby:            lobby,
		lobbyUseCase:     lobbyUseCase,
		placeUseCase:     placeUseCase,
		swipeUseCase:     swipeUseCase,
		userUseCase:      userUseCase,
		placeRecommender: placeRecommender,

		usersMap:   make(map[string]*domain.User),
		places:     make([]*domain.Place, 0),
		usersPlace: make(map[string]*domain.Place),
		swipes:     make([]*domain.Swipe, 0),
		matches:    make([]*nevent.Match, 0),
		votes:      make(map[int64]nevent.VoteAnnounce),
		log:        log.WithFields(log.Fields{"room": lobby.ID}),
	}

	r.votes[0] = nevent.FinishVote{
		BaseVote: nevent.BaseVote{
			ID: 0,
			Options: []nevent.VoteOption{
				{ID: nevent.OptionIDFinish, Desc: "Finish"},
				{ID: nevent.OptionIDContinue, Desc: "Continue"},
			},
			Type:  nevent.VoteTypeFinish,
			Votes: make(map[string]nevent.OptionID),
		},
	}

	r.state = lobby.State

	if r.state != domain.InLobby {
		// can't load swiping lobby
		err := r.setState(domain.Finished)
		if err != nil {
			return nil, fmt.Errorf("failed to set state: %w", err)
		}
	}

	r.log.Debugf("state: %s", r.state)

	err := r.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load room: %w", err)
	}

	return r, nil
}

func (r *NRoom) ID() string {
	return r.id
}

func (r *NRoom) Users() []*domain.User {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.users()
}

func (r *NRoom) users() []*domain.User {
	users := make([]*domain.User, 0)
	for _, user := range r.usersMap {
		users = append(users, user)
	}
	return users
}

func (r *NRoom) Load() error {
	r.lock.Lock()
	defer r.lock.Unlock()

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

func (r *NRoom) mathesFromSwipes() ([]*nevent.Match, error) {
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

	matches := make([]*nevent.Match, 0)

	for e := swipeCount.Oldest(); e != nil; e = e.Next() {
		r.log.WithFields(log.Fields{"place": e.Key, "count": e.Value}).Debug("swipe count")
		if e.Value == len(r.usersMap) {
			place, err := r.placeUseCase.GetPlaceByID(context.Background(), e.Key)
			if err != nil {
				return nil, fmt.Errorf("failed to get place: %w", err)
			}
			matches = append(matches, &nevent.Match{
				ID:    len(matches),
				Place: place,
			})
		}
	}

	return matches, nil
}

func (r *NRoom) setState(state domain.LobbyState) error {
	r.state = state
	err := r.lobbyUseCase.SetLobbyState(context.Background(), r.lobby.ID, state)
	if err != nil {
		return fmt.Errorf("failed to set lobby state: %w", err)
	}
	return nil
}

func (r *NRoom) OnJoin(c *state.Context[*NRoom]) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	switch r.state {
	case domain.InLobby:
		r.usersMap[c.User.ID] = c.User

		err := r.syncUsersWithBd(c.Ctx)
		if err != nil {
			return fmt.Errorf("failed to sync users with db: %w", err)
		}

		c.BroadcastToOthers(nevent.UserJoined{
			ID:     c.User.ID,
			Name:   c.User.Name,
			Avatar: c.User.Avatar,
		})
		for _, u := range r.users() {
			c.Emit(nevent.UserJoined{
				ID:     u.ID,
				Name:   u.Name,
				Avatar: u.Avatar,
			})
		}
		for _, v := range r.votes {
			c.Emit(v)
		}
	case domain.Finished:
		r.emitFinish(c)
	default:
		return fmt.Errorf("cannot join room in state %s", r.state)
	}

	c.Emit(nevent.SettingsUpdate{
		Location:    r.lobby.Location,
		PriceMin:    r.lobby.PriceAvg - 300,
		PriceMax:    r.lobby.PriceAvg + 300,
		MaxDistance: 4000,
		Tags: algo.Map(r.lobby.Tags, func(t *domain.Tag) int64 {
			return t.ID
		}),
		RecommendationOpts: r.recommendationOpts,
	})

	return nil
}

func (r *NRoom) OnLeave(c *state.Context[*NRoom]) error {
	_, has := r.usersMap[c.User.ID]
	if !has {
		return nil
	}
	delete(r.usersMap, c.User.ID)

	if r.state == domain.InLobby {
		err := r.syncUsersWithBd(c.Ctx)
		if err != nil {
			return fmt.Errorf("error while syncing users with bd: %w", err)
		}
	}

	c.BroadcastToOthers(nevent.UserLeft{
		ID: c.User.ID,
	})

	return nil
}

func (r *NRoom) Empty() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return len(r.usersMap) == 0
}

func (r *NRoom) OnLeaveLobby(c *state.Context[*NRoom], _ nevent.LeaveLobby) error {
	err := c.Close()
	if err != nil {
		return fmt.Errorf("error while closing connection: %w", err)
	}
	return nil
}

func (r *NRoom) syncUsersWithBd(ctx context.Context) error {
	var err error
	r.lobby.Users, err = r.lobbyUseCase.SetLobbyUsers(ctx,
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

func (r *NRoom) OnSettingsUpdate(c *state.Context[*NRoom], ev nevent.SettingsUpdate) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	err := r.updateLobbySettings(
		c.Ctx,
		ev.Location,
		(ev.PriceMax+ev.PriceMin)/2,
		ev.Tags,
		nil,
		ev.RecommendationOpts,
	)
	if err != nil {
		return fmt.Errorf("error while updating lobby settings: %w", err)
	}

	ev.UserID = c.User.ID
	c.Broadcast(ev)

	return nil
}

func (r *NRoom) updateLobbySettings(
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

func (r *NRoom) OnStartSwipes(c *state.Context[*NRoom], ev nevent.StartSwipes) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	c.Log.Debug("Request started: Action 'StartSwipes' initiated")

	if len(r.lobby.Tags) == 0 {
		c.Log.Warn("Action 'StartSwipes' encountered an issue. Reason: 'No tags found, using default tags'")
		err := r.updateLobbySettings(c.Ctx, r.lobby.Location, 500, []int64{3}, nil, r.recommendationOpts)
		if err != nil {
			c.Log.WithError(err).Error("Action 'UpdateLobby' failed")
			return err
		}
		c.Log.Debug("Request successful: Action 'UpdateLobby' completed with default tags")
	}

	var err error
	r.places, err = r.placeRecommender.RecommendPlaces(c.Ctx,
		r.recommendationOpts,
		domain.RecommendData{
			Location: r.lobby.Location,
			PriceAvg: r.lobby.PriceAvg,
			Tags:     r.lobby.TagNames(),
		},
	)
	if err != nil {
		c.Log.WithError(err).Error("Action 'GetPlacesForLobby' failed")
		return err
	}
	c.Log.Debug("Request successful: Action 'GetPlacesForLobby' completed")

	err = r.updateLobbySettings(c.Ctx, r.lobby.Location, r.lobby.PriceAvg,
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
	c.Log.Debug("Request successful: Action 'UpdateLobby' completed")

	c.Log.Debug("Request successful: Action 'UpdateLobby' completed")
	for id := range r.usersMap {
		if len(r.places) > 0 {
			r.usersPlace[id] = r.places[0]
			c.Log.Debugf("<User %s> is assigned to <Place %d>", id, r.places[0].ID)
		} else {
			log.Warnf("No places available to assign to <User %s>", id)
		}
	}

	log.Info("Swipes successfully started")
	err = r.setState(domain.Swiping)
	if err != nil {
		return fmt.Errorf("error while setting state: %w", err)
	}

	c.Broadcast(ev)
	c.ForEach(func(cc *state.Context[*NRoom]) {
		p := r.usersPlace[cc.User.ID]
		cc.Emit(nevent.Place{
			ID:   p.ID,
			Card: p,
		})
	})

	return nil
}

func (r *NRoom) OnSwipe(c *state.Context[*NRoom], ev nevent.Swipe) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	placeID := r.usersPlace[c.User.ID].ID

	swipe := &domain.Swipe{
		LobbyID: r.lobby.ID,
		PlaceID: placeID,
		UserID:  c.User.ID,
		Type:    ev.SwipeType,
	}

	r.swipes = append(r.swipes, swipe)
	err := r.swipeUseCase.SaveSwipe(context.Background(), swipe)
	if err != nil {
		return fmt.Errorf("error while saving swipe: %w", err)
	}

	pIdx := slices.IndexFunc(r.places, func(place *domain.Place) bool {
		return place.ID == placeID
	})
	r.usersPlace[c.User.ID] = r.places[(pIdx+1)%len(r.places)]

	likes := algo.Count(r.swipes, func(swipe *domain.Swipe) bool {
		return swipe.PlaceID == placeID && swipe.Type == domain.LIKE
	})

	if likes == len(r.usersMap) {
		match := &nevent.Match{Place: r.places[slices.IndexFunc(r.places, func(place *domain.Place) bool {
			return place.ID == placeID
		})]}
		match.ID = len(r.matches)
		r.matches = append(r.matches, match)
		r.result = match.Place

		vote := nevent.MatchVote{
			BaseVote: nevent.BaseVote{
				ID: int64(len(r.votes)),
				Options: []nevent.VoteOption{
					{ID: nevent.OptionIDLike, Desc: "Like"},
					{ID: nevent.OptionIDDislike, Desc: "Dislike"},
				},
				Type:  nevent.VoteTypeMatch,
				Votes: make(map[string]nevent.OptionID),
			},
			Place: match.Place,
		}
		r.votes[vote.ID] = vote
		c.Broadcast(vote)
	}

	c.Emit(nevent.Place{
		ID:   r.usersPlace[c.User.ID].ID,
		Card: r.usersPlace[c.User.ID],
	})

	return nil
}

func (r *NRoom) OnVote(c *state.Context[*NRoom], ev nevent.Vote) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	vote, ok := r.votes[ev.VoteID]
	if !ok {
		return fmt.Errorf("vote with id %d not found", ev.VoteID)
	}

	var res *nevent.VoteResult

	switch v := vote.(type) {
	case nevent.MatchVote:
		if ev.OptionID != nevent.OptionIDLike && ev.OptionID != nevent.OptionIDDislike {
			return fmt.Errorf("invalid option id: %d", ev.OptionID)
		}
		v.Votes[c.User.ID] = ev.OptionID

		log.Debugf("<User %s> voted for <Option %d> in %s", c.User.ID, ev.OptionID, v.Type)
		res = r.voteMatchResult(v)

		if res != nil && res.OptionID == nevent.OptionIDLike {
			r.result = v.Place
			err := r.setState(domain.Finished)
			if err != nil {
				return fmt.Errorf("error while setting lobby state: %w", err)
			}
		}

	case nevent.FinishVote:
		if ev.OptionID != nevent.OptionIDFinish && ev.OptionID != nevent.OptionIDContinue {
			return fmt.Errorf("invalid option id: %d", ev.OptionID)
		}
		v.Votes[c.User.ID] = ev.OptionID

		log.Debugf("<User %s> voted for <Option %d> in %s", c.User.ID, ev.OptionID, v.Type)
		res = r.voteFinishResult(v)
		if res != nil && res.OptionID == nevent.OptionIDFinish {
			err := r.setState(domain.Finished)
			if err != nil {
				return fmt.Errorf("error while setting lobby state: %w", err)
			}
		}

	default:
		log.Warnf("Unknown voting type: %T", vote)
		return fmt.Errorf("unknown voting type: %T", vote)
	}

	c.Broadcast(nevent.Voted{
		VoteID:   ev.VoteID,
		OptionID: ev.OptionID,
		User: nevent.UserJoined{
			ID:     c.User.ID,
			Name:   c.User.Name,
			Avatar: c.User.Avatar,
		},
	})

	if res != nil {
		c.Broadcast(res)
	}

	if r.state == domain.Finished {
		c.ForEach(r.emitFinish)
	}

	return nil
}

func (r *NRoom) voteMatchResult(vote nevent.MatchVote) *nevent.VoteResult {
	if len(vote.Votes) != len(r.usersMap) {
		return nil
	}
	likes := 0
	for _, vote := range vote.Votes {
		if vote == nevent.OptionIDLike {
			likes++
		}
	}
	if likes == len(r.usersMap) {
		return &nevent.VoteResult{
			VoteID:   vote.ID,
			Type:     nevent.VoteTypeMatch,
			OptionID: nevent.OptionIDLike,
		}
	}

	return &nevent.VoteResult{
		VoteID:   vote.ID,
		Type:     nevent.VoteTypeMatch,
		OptionID: nevent.OptionIDDislike,
	}
}

func (r *NRoom) voteFinishResult(vote nevent.FinishVote) *nevent.VoteResult {
	if len(vote.Votes) != len(r.usersMap) {
		return nil
	}
	likes := 0
	for _, vote := range vote.Votes {
		if vote == nevent.OptionIDFinish {
			likes++
		}
	}
	if likes == len(r.usersMap) {
		return &nevent.VoteResult{
			VoteID:   vote.ID,
			Type:     nevent.VoteTypeFinish,
			OptionID: nevent.OptionIDFinish,
		}
	}
	return nil
}

func (r *NRoom) emitFinish(c *state.Context[*NRoom]) {
	var matchResult *domain.Place
	if len(r.matches) > 0 {
		matchResult = r.matches[len(r.matches)-1].Place
	}

	c.Emit(nevent.Finish{
		Result:  matchResult,
		Matches: r.matches,
	})
}
