package util

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AppConfigPathTestSuite struct {
	suite.Suite
}

func (suite *AppConfigPathTestSuite) SetupTest() {
	removeTestFiles()
}

func (suite *AppConfigPathTestSuite) TearDownTest() {
	removeTestFiles()
}

func (suite *AppConfigPathTestSuite) TestReturnsEmptyWithoutRanchfiles() {
	assert := assert.New(suite.T())
	cmd := mockCmd("")

	file, err := AppConfigPath(&cmd)
	assert.Equal("", file)
	assert.Nil(err)
}

func (suite *AppConfigPathTestSuite) TestErrorsIfFileFlagDoesNotExist() {
	assert := assert.New(suite.T())
	cmd := mockCmd(fmt.Sprintf("-f %s", TEST_FILE_1))

	file, err := AppConfigPath(&cmd)
	assert.Equal("", file)
	assert.Error(err)
}

func (suite *AppConfigPathTestSuite) TestRespectsFileFlagIfExists() {
	assert := assert.New(suite.T())
	cmd := mockCmd(fmt.Sprintf("-f %s", TEST_FILE_1))

	os.Create(TEST_FILE_1)

	file, err := AppConfigPath(&cmd)
	assert.Equal(TEST_FILE_1, file)
	assert.Nil(err)
}

func (suite *AppConfigPathTestSuite) TestReturnsDefaultRanchfileIfPresent() {
	assert := assert.New(suite.T())
	cmd := mockCmd("")

	os.Create(TEST_RANCHY)

	file, err := AppConfigPath(&cmd)
	assert.Equal(TEST_RANCHY, file)
	assert.Nil(err)
}

func (suite *AppConfigPathTestSuite) TestErrorsIfExtraRanchfilesPresent() {
	assert := assert.New(suite.T())
	cmd := mockCmd("")

	os.Create(TEST_RANCHY)
	os.Create(TEST_FILE_1)

	file, err := AppConfigPath(&cmd)
	assert.Equal("", file)
	assert.Error(err)
}

func TestAppConfigPathTestSuite(t *testing.T) {
	suite.Run(t, new(AppConfigPathTestSuite))
}
