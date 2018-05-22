package util

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AppNameTestSuite struct {
	suite.Suite
}

func (suite *AppNameTestSuite) SetupTest() {
	removeTestFiles()
}

func (suite *AppNameTestSuite) TearDownTest() {
	removeTestFiles()
}

func (suite *AppNameTestSuite) TestUsesAppFlag() {
	assert := assert.New(suite.T())
	testCmd := mockCmd("-a superuberfantastic")

	name, err := AppName(&testCmd)
	assert.Equal("superuberfantastic", name)
	assert.Nil(err)
}

func (suite *AppNameTestSuite) TestErrorsIfDefaultRanchfileIsEmpty() {
	assert := assert.New(suite.T())
	testCmd := mockCmd("")

	os.Create(TEST_RANCHY)

	name, err := AppName(&testCmd)
	assert.Equal("", name)
	assert.Error(err)
}

func (suite *AppNameTestSuite) TestUsesDefaultRanchfile() {
	assert := assert.New(suite.T())
	testCmd := mockCmd("")

	err := ioutil.WriteFile(TEST_RANCHY, []byte("name: fantastic\n"), 0644)
	assert.Nil(err)

	name, err := AppName(&testCmd)
	assert.Equal("fantastic", name)
	assert.Nil(err)
}

func (suite *AppNameTestSuite) TestErrorsIfExtraRanchfilesPresent() {
	assert := assert.New(suite.T())
	testCmd := mockCmd("")

	err := ioutil.WriteFile(TEST_RANCHY, []byte("name: fantastic\n"), 0644)
	assert.Nil(err)

	os.Create(TEST_FILE_1)

	name, err := AppName(&testCmd)
	assert.Equal("", name)
	assert.Error(err)
}

func TestAppNameTestSuite(t *testing.T) {
	suite.Run(t, new(AppNameTestSuite))
}
