package vio

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVCSHasUncommittedChanges(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, err)

	err = ioutil.WriteFile(path+"/bar", []byte(""), 0644)
	assert.Nil(t, err)

	createAndSeedTestRepo(t, path, []string{"bar"})

	has, err := HasUncommittedChanges(path)
	assert.Nil(t, err)
	assert.False(t, has)

	// unversioned files shouldn't be considered
	err = ioutil.WriteFile(path+"/untracked", []byte("changed"), 0644)
	assert.Nil(t, err)

	has, err = HasUncommittedChanges(path)
	assert.Nil(t, err)
	assert.False(t, has)

	// only changes to committed files should be considered
	err = ioutil.WriteFile(path+"/bar", []byte("changed"), 0644)
	assert.Nil(t, err)

	has, err = HasUncommittedChanges(path)
	assert.Nil(t, err)
	assert.True(t, has)
}

func TestGetCurrentCommitId(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, err)

	createAndSeedTestRepo(t, path, []string{})

	id, err := GetCurrentCommitId(path)
	assert.Nil(t, err)
	assert.Equal(t, len(id), 7)
	r, err := regexp.Compile(`[0123456789abcdef]+`)
	assert.Nil(t, err)
	assert.True(t, r.MatchString(id))
}

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

func TestGetVersionedFiles(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, err)

	createAndSeedTestRepo(t, path, []string{})

	files, err := GetVersionedFiles(path)
	assert.Nil(t, err)
	assert.Equal(t, len(files), 2)
}

func TestVersionToString(t *testing.T) {
	v := NewVersion("1234567890")
	assert.NotNil(t, v)
	assert.Equal(t, fmt.Sprintf("%v", v), "1234567890 "+v.timestamp.Format(time.RFC3339))

	ts_str := "2012-11-01T22:08:41Z"
	v = NewVersion("1234567890 " + ts_str)
	ts, err := time.Parse(time.RFC3339, ts_str)
	assert.Nil(t, err)
	assert.NotNil(t, ts)
	assert.NotNil(t, v)
	assert.Equal(t, fmt.Sprintf("%v", v), "1234567890 "+ts_str)
}

func TestContainsVersion(t *testing.T) {
	v1 := NewVersion("1234567890 2012-11-01T22:08:41Z")
	assert.NotNil(t, v1)
	v2 := NewVersion("5713943128 2015-02-03T23:59:12Z")
	assert.NotNil(t, v2)

	var vs []version
	assert.False(t, ContainsVersion(vs, v1))
	vs = append(vs, *NewVersion("1234567890 2012-11-01T22:08:41Z"))
	assert.Equal(t, len(vs), 1)
	assert.True(t, ContainsVersion(vs, v1))
	assert.False(t, ContainsVersion(vs, v2))
	vs = append(vs, *NewVersion("5713943128 2015-02-03T23:59:12Z"))
	assert.True(t, ContainsVersion(vs, v1))
	assert.True(t, ContainsVersion(vs, v2))
}
