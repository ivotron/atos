package vio

import (
	"io/ioutil"
	"regexp"
	"testing"

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

func TestGetVersionedFiles(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, err)

	createAndSeedTestRepo(t, path, []string{})

	files, err := GetVersionedFiles(path)
	assert.Nil(t, err)
	assert.Equal(t, len(files), 2)
}
