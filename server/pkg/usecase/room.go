package usecase

import (
	"context"
	"fmt"
	"slices"
	"sync"

	algo "dishdash.ru/pkg/algo"
	"dishdash.ru/pkg/domain"
	"dishdash.ru/pkg/usecase/event"
	"dishdash.ru/pkg/usecase/state"
	log "github.com/sirupsen/logrus"
)

type Room struct {
	lobby *domain.Lobby

	lock sync.RWMutex

	usersMap       map[string]*domain.User
	connectedUsers map[string]*domain.User

	userSwiped    map[string]int // count of swiped cards
	userCardsSeen map[string]int // count of cards sended to user
	cards         []*domain.Place
	swipes        []*domain.Swipe
	results       event.Results

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
		lobby:            lobby,
		lobbyUseCase:     lobbyUseCase,
		placeUseCase:     placeUseCase,
		swipeUseCase:     swipeUseCase,
		userUseCase:      userUseCase,
		placeRecommender: placeRecommender,

		usersMap:       make(map[string]*domain.User),
		connectedUsers: make(map[string]*domain.User),
		cards:          lobby.Places,
		userSwiped:     make(map[string]int),
		userCardsSeen:  make(map[string]int),
		swipes:         lobby.Swipes,
		log:            log.WithFields(log.Fields{"room": lobby.ID}),
	}

	r.log.Debugf("state: %s", r.lobby.State)

	err := r.load()
	if err != nil {
		return nil, fmt.Errorf("failed to load room: %w", err)
	}

	return r, nil
}

func (r *Room) ID() string {
	return r.lobby.ID
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

func (r *Room) load() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, user := range r.lobby.Users {
		r.usersMap[user.ID] = user
	}

	err := r.evalUserCards()
	if err != nil {
		return fmt.Errorf("failed to eval user cards: %w", err)
	}

	r.results, err = r.evalResults()
	if err != nil {
		return fmt.Errorf("failed to eval results: %w", err)
	}

	return nil
}

func (r *Room) evalUserCards() error {
	for _, swipe := range r.swipes {
		r.userSwiped[swipe.UserID]++
	}
	for id, swiped := range r.userSwiped {
		r.userCardsSeen[id] = swiped
	}
	return nil
}

func (r *Room) evalResults() (event.Results, error) {
	card2Likes := make(map[int64][]string)

	for _, swipe := range r.swipes {
		if swipe.Type == domain.LIKE {
			card2Likes[swipe.CardID] = append(card2Likes[swipe.CardID], swipe.UserID)
		}
	}

	cards := make(map[int64]*domain.Place)
	card2Position := make(map[int64]int)
	for i, c := range r.cards {
		cards[c.ID] = c
		card2Position[c.ID] = i
	}

	top := make([]event.TopPosition, 0)

	for cardID, likes := range card2Likes {
		card := cards[cardID]
		if card == nil {
			continue
		}
		top = append(top, event.TopPosition{
			Card: card,
			Likes: algo.Map(likes, func(id string) *domain.User {
				return r.usersMap[id]
			}),
		})
	}

	slices.SortFunc(top, func(a, b event.TopPosition) int {
		if len(b.Likes) != len(a.Likes) {
			return len(b.Likes) - len(a.Likes)
		}
		return card2Position[a.Card.ID] - card2Position[b.Card.ID]
	})

	return event.Results{Top: top}, nil
}

func (r *Room) setState(state domain.LobbyState) error {
	r.lobby.State = state
	err := r.lobbyUseCase.SetLobbyState(context.Background(), r.lobby.ID, state)
	if err != nil {
		return fmt.Errorf("failed to set lobby state: %w", err)
	}
	return nil
}

func (r *Room) OnJoin(c *state.Context[*Room]) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.connectedUsers[c.User.ID] = c.User
	c.Log.Debug("user joined lobby")
	switch r.lobby.State {
	case domain.InLobby:
		r.usersMap[c.User.ID] = c.User
		err := r.syncUsersWithBd(c.Ctx)
		if err != nil {
			return fmt.Errorf("failed to sync users with db: %w", err)
		}

	case domain.Swiping:
		_, has := r.usersMap[c.User.ID]
		if !has {
			r.usersMap[c.User.ID] = c.User
			r.userCardsSeen[c.User.ID] = 0
			r.userSwiped[c.User.ID] = 0
		}
		r.userCardsSeen[c.User.ID] = r.userSwiped[c.User.ID]
		c.Emit(event.StartSwipes{})
		r.emitCardsForUser(c, c.User.ID)
		c.Emit(r.results)
	}
	c.Log.Debug("user joined lobby")

	c.BroadcastToOthers(event.UserJoined{
		ID:     c.User.ID,
		Name:   c.User.Name,
		Avatar: c.User.Avatar,
	})
	for _, u := range r.users() {
		c.Emit(event.UserJoined{
			ID:     u.ID,
			Name:   u.Name,
			Avatar: u.Avatar,
		})
	}

	c.Emit(event.SettingsUpdate(r.lobby.Settings))

	return nil
}

func (r *Room) emitCardsForUser(c *state.Context[*Room], id string) {
	swiped := r.userSwiped[id]
	seen := r.userCardsSeen[id]

	cards := make([]*domain.Place, 0)
	for seen < swiped+3 && seen < len(r.cards) {
		cards = append(cards, r.cards[seen])
		seen++
	}
	r.userCardsSeen[id] = seen
	c.Emit(event.Cards{
		Cards: cards,
	})
}

func (r *Room) OnLeave(c *state.Context[*Room]) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, has := r.usersMap[c.User.ID]
	if !has {
		return nil
	}
	delete(r.connectedUsers, c.User.ID)

	c.BroadcastToOthers(event.UserLeft{
		ID: c.User.ID,
	})

	return nil
}

func (r *Room) Active() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return len(r.connectedUsers) > 0
}

func (r *Room) Empty() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return len(r.usersMap) == 0
}

func (r *Room) OnLeaveLobby(c *state.Context[*Room], _ event.LeaveLobby) error {
	err := c.Close()
	if err != nil {
		return fmt.Errorf("error while closing connection: %w", err)
	}
	return nil
}

func (r *Room) syncUsersWithBd(ctx context.Context) error {
	var err error
	r.lobby.Users, err = r.lobbyUseCase.SetLobbyUsers(ctx,
		r.lobby.ID,
		algo.Map(r.users(), func(u *domain.User) string {
			return u.ID
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to set lobby state: %w", err)
	}
	return nil
}

func (r *Room) OnSettingsUpdate(c *state.Context[*Room], ev event.SettingsUpdate) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	err := r.updateLobbySettings(
		c.Ctx,
		domain.LobbySettings(ev),
	)
	if err != nil {
		return fmt.Errorf("error while updating lobby settings: %w", err)
	}

	c.Broadcast(ev)

	return nil
}

func (r *Room) updateLobbySettings(
	ctx context.Context,
	settings domain.LobbySettings,
) error {
	err := r.lobbyUseCase.SetLobbySettings(ctx, r.lobby.ID, settings)
	if err != nil {
		return err
	}
	r.lobby.Settings = settings
	return nil
}

func (r *Room) OnStartSwipes(c *state.Context[*Room], ev event.StartSwipes) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	c.Log.Debug("Starting swipes")

	var err error
	r.cards, err = r.placeRecommender.RecommendPlaces(c.Ctx,
		r.lobby.Settings,
	)
	if err != nil {
		c.Log.WithError(err).Error("Action 'GetPlacesForLobby' failed")
		return err
	}
	c.Log.Debugf("Get %d places from recommender", len(r.cards))

	err = r.lobbyUseCase.AttachOrderedPlacesToLobby(c.Ctx,
		algo.Map(r.cards, func(p *domain.Place) int64 { return p.ID }),
		r.lobby.ID,
	)
	if err != nil {
		return fmt.Errorf("error while attaching places to lobby: %w", err)
	}
	c.Log.Debugf("Attached %d places to lobby", len(r.cards))

	err = r.updateLobbySettings(c.Ctx, r.lobby.Settings)
	if err != nil {
		return fmt.Errorf("error while updating lobby settings: %w", err)
	}

	log.Info("Swipes successfully started")
	err = r.setState(domain.Swiping)
	if err != nil {
		return fmt.Errorf("error while setting state: %w", err)
	}

	c.Broadcast(ev)
	c.ForEach(func(cc *state.Context[*Room]) {
		r.emitCardsForUser(cc, cc.User.ID)
	})

	return nil
}

func (r *Room) OnSwipe(c *state.Context[*Room], ev event.Swipe) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	card := r.cards[r.userSwiped[c.User.ID]]

	swipe := &domain.Swipe{
		LobbyID: r.lobby.ID,
		CardID:  card.ID,
		UserID:  c.User.ID,
		Type:    ev.SwipeType,
	}

	r.swipes = append(r.swipes, swipe)
	err := r.swipeUseCase.SaveSwipe(context.Background(), swipe)
	if err != nil {
		return fmt.Errorf("error while saving swipe: %w", err)
	}

	r.userSwiped[c.User.ID]++

	if ev.SwipeType == domain.LIKE {
		likes := algo.Count(r.swipes, func(swipe *domain.Swipe) bool {
			return swipe.CardID == card.ID && swipe.Type == domain.LIKE
		})

		if likes == len(r.usersMap) {
			c.Broadcast(event.Match{Card: card})
		}

		r.results, err = r.evalResults()
		if err != nil {
			return fmt.Errorf("error while evaluating results: %w", err)
		}

		c.Broadcast(r.results)
	}

	r.emitCardsForUser(c, c.User.ID)

	return nil
}
