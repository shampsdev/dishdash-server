package sdk

import (
	"flag"
	"fmt"
	"testing"

	"dishdash.ru/cmd/server/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type SessionTest struct {
	GoldenFile string
	Run        func(t *testing.T) *SocketIOSession
}

var updateGolden = flag.Bool("update-golden", false, "Update golden files")

func RunSessionTest(t *testing.T, td SessionTest) {
	flag.Parse()
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetLevel(log.DebugLevel)

	config.Load("../e2e.env")
	err := SetupDB()
	assert.NoError(t, err)
	defer func() {
		err = CleanDB()
		if err != nil {
			fmt.Printf("error while cleaning: %s", err.Error())
		}
	}()

	var s *SocketIOSession
	defer func() {
		if t.Failed() && s != nil {
			_ = s.Save(goldenPath("ERROR_" + td.GoldenFile))
		}
	}()

	s = td.Run(t)

	gp := goldenPath(td.GoldenFile)
	if *updateGolden {
		assert.NoError(t, s.Save(gp))
	} else {
		exp, err := LoadSocketIOSession(gp)
		assert.NoError(t, err)
		AssertSocketIOSession(t, exp, s)
	}
}

func goldenPath(name string) string {
	return name + ".golden.json"
}
