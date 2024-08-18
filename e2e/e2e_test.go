package e2e

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"dishdash.ru/e2e/tests"

	"dishdash.ru/cmd/server/config"
	"dishdash.ru/e2e/pg_test"
	"dishdash.ru/e2e/server_test"
	"dishdash.ru/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var updateGolden = flag.Bool("update-golden", false, "Update golden files")

type E2ETestSuite struct {
	suite.Suite
	testDB *pg_test.TestDatabase
	cases  usecase.Cases

	stopServer context.CancelFunc
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func (suite *E2ETestSuite) SetupSuite() {
	suite.testDB = pg_test.SetupTestDatabase()
	suite.cases = usecase.Setup(suite.testDB.DbInstance)
	suite.stopServer = server_test.StartServer(suite.cases)
	config.Print()
	tests.SIOHost = fmt.Sprintf("http://localhost:%d", config.C.Server.Port)
	tests.ApiHost = fmt.Sprintf("http://localhost:%d/api/v1", config.C.Server.Port)
}

func (suite *E2ETestSuite) TearDownSuite() {
	suite.stopServer()
	time.Sleep(3 * time.Second)
	suite.testDB.TearDown()
}

func (suite *E2ETestSuite) Test_GetAllTags() {
	t := suite.T()
	err := pg_test.ResetData(suite.cases)
	assert.NoError(t, err)
	tests.GetAllTags(t)
}

func (suite *E2ETestSuite) Test_UpdateUser() {
	t := suite.T()
	err := pg_test.ResetData(suite.cases)
	assert.NoError(t, err)
	tests.UpdateUser(t)
}

type sessionTest struct {
	Name       string
	GoldenFile string
	Run        func(t *testing.T) *tests.SocketIOSession
}

var sessionTests = []sessionTest{
	{
		Name:       "LobbyJoin",
		GoldenFile: "lobby_join",
		Run:        tests.LobbyJoin,
	},
	{
		Name:       "LobbySwipe",
		GoldenFile: "lobby_swipe",
		Run:        tests.LobbySwipe,
	},
}

func (suite *E2ETestSuite) Test_SessionTests() {
	for _, td := range sessionTests {
		suite.T().Run(td.Name, func(t *testing.T) {
			err := pg_test.ResetData(suite.cases)
			assert.NoError(t, err)
			s := td.Run(t)
			gp := goldenPath(td.GoldenFile)
			if *updateGolden {
				assert.NoError(t, s.Save(gp))
			} else {
				exp, err := tests.LoadSocketIOSession(gp)
				assert.NoError(t, err)
				tests.AssertSocketIOSession(t, exp, s)
				if t.Failed() {
					assert.NoError(t, s.Save(goldenPath("ERROR_"+td.GoldenFile)))
				}
			}
		})
	}
}

func goldenPath(name string) string {
	return filepath.Join("testdata", name+".golden.json")
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
