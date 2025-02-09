package framework

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"dishdash.ru/cmd/server/config"
	"dishdash.ru/e2e/framework/session"
	"dishdash.ru/pkg/domain"
	socketio "github.com/googollee/go-socket.io"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Framework struct {
	Cfg config.Config
	DB  *pgxpool.Pool
	Log *logrus.Logger

	HttpCli *http.Client
	ApiHost string
	SIOHost string

	Session *session.Session
}

func MustInit() *Framework {
	fw := &Framework{}

	fw.Session = session.New()

	fw.Log = logrus.New()
	fw.Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	fw.Log.SetLevel(logrus.DebugLevel)

	if !isE2ETesting() {
		return fw
	}

	config.Load("../e2e.env")
	fw.Cfg = config.C
	fw.Log.Info("Loaded config")

	var err error
	fw.DB, err = pgxpool.NewWithConfig(context.Background(), fw.Cfg.PGXConfig())
	if err != nil {
		panic(err)
	}
	fw.Log.Info("Connected to database")

	fw.HttpCli = &http.Client{Timeout: 10 * time.Second}
	fw.ApiHost = "http://localhost:8001/api/v1"
	fw.SIOHost = "http://localhost:8001"

	return fw
}

func (fw *Framework) MustNewClient(user *domain.User) *Client {
	c, err := fw.NewClient(user)
	if err != nil {
		panic(err)
	}
	return c
}

func (fw *Framework) NewClient(user *domain.User) (*Client, error) {
	cli, err := socketio.NewClient(fw.SIOHost, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	user, err = fw.postUserWithID(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	c := &Client{
		fw:   fw,
		User: user,
		cli:  cli,
	}
	c.Log = fw.Log.WithFields(logrus.Fields{
		"user": user.ID,
	})
	c.setup(allEvents)

	err = c.cli.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect client: %w", err)
	}

	return c, nil
}

func (fw *Framework) RecordEvents(events ...string) {
	fw.Session.SetRecordEvents(events...)
}

// Step
// 1. Adds a new step to the session
// 2. Runs the function
// 3. Locks until the number of responses to be recorded is reached
//
// It is important that all incoming events are counted and not just recorded ones.
func (fw *Framework) Step(name string, f func(), waitNResponses uint32) {
	fw.Log.Infof("Step: %s", name)
	fw.Session.NewStep(name)
	f()
	fw.Session.WaitNResponses(waitNResponses)
}
