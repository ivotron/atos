package vio

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdInitPosix(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, os.Chdir(path))
	createAndSeedTestRepo(t, path, []string{})
	err = Init(path, path+"/.snapshots", "posix")
	assert.Nil(t, err)
	_, err = os.Stat(path + "/.vioconfig")
	assert.Nil(t, err)
}
func TestCmdCommitPosix(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, os.Chdir(path))
	createAndSeedTestRepo(t, path, []string{})
	err = Init(path, path+"/.snapshots", "posix")
	assert.Nil(t, err)
	err = ioutil.WriteFile(path+"/bar", []byte("yeah"), 0644)
	assert.Nil(t, err)
	err = ioutil.WriteFile(path+"/toz", []byte("ok"), 0644)
	assert.Nil(t, err)
	err = Commit()
	assert.Nil(t, err)
}
