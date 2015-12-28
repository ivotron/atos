package vio

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"gopkg.in/ini.v1"

	"github.com/stretchr/testify/assert"
)

func createAndSeedTestRepo(t *testing.T, repoPath string, filesToAdd []string) {
	_, err := runCmd(repoPath, "git init")
	assert.Nil(t, err)

	err = ioutil.WriteFile(repoPath+"/README", []byte("foo\n"), 0644)
	assert.Nil(t, err)

	err = ioutil.WriteFile(repoPath+"/.gitignore", []byte(".snapshots\n"), 0644)
	assert.Nil(t, err)

	_, err = runCmd(repoPath, "git add *")
	assert.Nil(t, err)

	for _, file := range filesToAdd {
		_, err = runCmd(repoPath, "git add "+file)
		assert.Nil(t, err)
	}

	_, err = runCmd(repoPath, "git commit -m yeah")
	assert.Nil(t, err)

	return
}

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

func TestVersionToString(t *testing.T) {
	v := NewVersion("1234567890")
	assert.NotNil(t, v)
	assert.Equal(t, fmt.Sprintf("%v", v), "1234567890#"+fmt.Sprint(v.timestamp.Unix()))

	ts_str := "1405544146"
	v = NewVersion("1234567890#" + ts_str)
	i, err := strconv.ParseInt(ts_str, 10, 64)
	assert.Nil(t, err)
	ts := time.Unix(i, 0)
	assert.Nil(t, err)
	assert.NotNil(t, ts)
	assert.NotNil(t, v)
	assert.Equal(t, fmt.Sprintf("%v", v), "1234567890#"+ts_str)
}

func TestContainsVersion(t *testing.T) {
	v1 := NewVersion("1234567890#1405544146")
	assert.NotNil(t, v1)
	v2 := NewVersion("5713943128#2435869343")
	assert.NotNil(t, v2)

	var vs []version
	assert.False(t, ContainsVersion(vs, v1))
	vs = append(vs, *NewVersion("1234567890#1405544146"))
	assert.Equal(t, len(vs), 1)
	assert.True(t, ContainsVersion(vs, v1))
	assert.False(t, ContainsVersion(vs, v2))
	vs = append(vs, *NewVersion("5713943128#2435869343"))
	assert.True(t, ContainsVersion(vs, v1))
	assert.True(t, ContainsVersion(vs, v2))
}
