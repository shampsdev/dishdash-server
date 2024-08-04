package e2e

import (
	"context"
	"dishdash.ru/e2e/tests"
	"fmt"
	"testing"
	"time"

	"dishdash.ru/cmd/server/config"
	"dishdash.ru/e2e/pg_test"
	"dishdash.ru/e2e/server_test"
	"dishdash.ru/internal/usecase"

	"github.com/stretchr/testify/suite"
	"gotest.tools/v3/assert"
)

var host string

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
	host = fmt.Sprintf("http://localhost:%d/api/v1", config.C.Server.Port)
}

func (suite *E2ETestSuite) TearDownSuite() {
	suite.stopServer()
	time.Sleep(3 * time.Second)
	suite.testDB.TearDown()
}

func (suite *E2ETestSuite) Test_GetAllTags() {
	t := suite.T()
	err := pg_test.ResetData(suite.cases)
	assert.NilError(t, err)
	tests.GetAllTags(t, host)
}

func (suite *E2ETestSuite) Test_UpdateUser() {
	t := suite.T()
	err := pg_test.ResetData(suite.cases)
	assert.NilError(t, err)
	tests.UpdateUser(t, host)
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
