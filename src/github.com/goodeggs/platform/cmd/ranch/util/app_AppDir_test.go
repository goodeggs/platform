package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AppDirTestSuite struct {
	suite.Suite
}

func (suite *AppDirTestSuite) SetupTest() {
	removeTestFiles()
}

func (suite *AppDirTestSuite) TearDownTest() {
	removeTestFiles()
}

func (suite *AppDirTestSuite) TestWorks() {
	assert := assert.New(suite.T())
	testCmd := mockCmd("")

	_, err := AppDir(&testCmd)
	assert.Nil(err)
}

func TestAppDirTestSuite(t *testing.T) {
	suite.Run(t, new(AppDirTestSuite))
}
