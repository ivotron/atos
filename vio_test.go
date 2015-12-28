package vio

import (
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/ini.v1"

	"github.com/stretchr/testify/assert"
)

func testCmdInitPosix(t *testing.T, snapsPath string) {
	err := Init(snapsPath, "posix")
	assert.Nil(t, err)
	_, err = os.Stat(".vioconfig")
	assert.Nil(t, err)
	cfg, err := ini.Load(".vioconfig")
	assert.Nil(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, cfg.Section("").Key("snapshots_path").String(), snapsPath)
}

func TestCmdInitPosixRelative(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, err)
	assert.Nil(t, os.Chdir(path))
	createAndSeedTestRepo(t, path, []string{})
	testCmdInitPosix(t, ".snapshots")
}
func TestCmdInitPosixAbsolute(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, err)
	assert.Nil(t, os.Chdir(path))
	createAndSeedTestRepo(t, path, []string{})
	testCmdInitPosix(t, path+"/.snapshots")
}
func TestCmdCommitPosix(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, os.Chdir(path))
	createAndSeedTestRepo(t, path, []string{})
	err = Init(".snapshots", "posix")
	assert.Nil(t, err)
	err = ioutil.WriteFile(path+"/bar", []byte("yeah"), 0644)
	assert.Nil(t, err)
	err = ioutil.WriteFile(path+"/toz", []byte("ok"), 0644)
	assert.Nil(t, err)
	err = Commit()
	assert.Nil(t, err)
}
