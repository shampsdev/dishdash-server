package e2e

import (
	"context"
	"dishdash.ru/cmd/server/config"
	"dishdash.ru/e2e/pg_test"
	"dishdash.ru/e2e/server_test"
	"dishdash.ru/internal/repo/pg"
	"github.com/stretchr/testify/suite"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

type E2ETestSuite struct {
	suite.Suite
	testDB *pg_test.TestDatabase

	stopServer context.CancelFunc
}

func (suite *E2ETestSuite) SetupSuite() {
	config.Print()
	suite.testDB = pg_test.SetupTestDatabase()
	if err := pg.MigrateDB(); err != nil {
		suite.T().Fatal(err)
	}
	suite.stopServer = server_test.StartServer(suite.testDB.DbInstance)
}

func (suite *E2ETestSuite) TearDownSuite() {
	suite.stopServer()
	time.Sleep(3 * time.Second)
	suite.testDB.TearDown()
}

func (suite *E2ETestSuite) Test1() {
	time.Sleep(3 * time.Second)
	t := suite.T()
	assert.Equal(t, 1, 1)
}

func (suite *E2ETestSuite) Test2() {
	time.Sleep(3 * time.Second)
	t := suite.T()
	assert.Equal(t, 2, 2)
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
