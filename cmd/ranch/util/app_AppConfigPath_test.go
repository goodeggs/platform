package util

import (
	"testing"
	"os"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AppConfigPathTestSuite struct {
	suite.Suite
}

func (suite *AppConfigPathTestSuite) SetupSuite() {
	removeTestFiles()
}

func (suite *AppConfigPathTestSuite) TestRun() {
	assert := assert.New(suite.T())

	// create .ranch.yaml file if it 
	os.Create(TEST_RANCHY)

	// test case: no config file passed, attempt to use .ranch.yaml of cwd
	file, err := _appConfigPath("")
	assert.Equal(TEST_RANCHY, file, "file returned should be '.ranch.yaml'")
	assert.Nil(err)

	// test case: no config file passed, scans cwd for .ranch.*.yaml, finds one,
	// so doesn't require you to specify -f
	os.Create(TEST_FILE_1)
	file, err = _appConfigPath("")
	assert.Equal(TEST_RANCHY, file, "file returned should be '.ranch.yaml'")
	assert.Nil(err)

	// test case: no config file passed, scans cwd for .ranch.*.yaml, finds two,
	// so then should error
	os.Create(TEST_FILE_2)

	file, err = _appConfigPath("")
	assert.Equal("", file)
	assert.Error(err, "should complain about multiple ranch configs")

	// test case: filename flat is provided
	file, err = _appConfigPath(TEST_FILE_1)
	assert.Equal(TEST_FILE_1, file, "returned file should be test file 1")
	assert.Nil(err)
}

func (suite *AppConfigPathTestSuite) TestTearDown() {
	removeTestFiles()
}

func TestAppConfigPathTestSuite(t *testing.T) {
	suite.Run(t, new(AppConfigPathTestSuite))
}
