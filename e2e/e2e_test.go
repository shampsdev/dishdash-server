package e2e

import (
	"context"
	"fmt"
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

type E2ETestSuite struct {
	suite.Suite
	testDB *pg_test.TestDatabase
	cases  usecase.Cases

	stopServer context.CancelFunc
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

func (suite *E2ETestSuite) Test_JoinLobby() {
	t := suite.T()
	err := pg_test.ResetData(suite.cases)
	assert.NoError(t, err)
	tests.JoinLobby(t)
}

func (suite *E2ETestSuite) Test_SwipeLobby() {
	t := suite.T()
	err := pg_test.ResetData(suite.cases)
	assert.NoError(t, err)
	tests.SwipeLobby(t)
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}