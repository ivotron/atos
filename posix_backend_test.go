package vio

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPosixBackendInit(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, os.Chdir(path))

	opts := NewOptions()
	assert.NotNil(t, opts)

	opts.RepoPath = path
	opts.SnapshotsPath = path + "/.snapshots"

	backend, err := NewPosixBackend(opts)
	assert.NotNil(t, backend)
	assert.Nil(t, err)

	assert.False(t, backend.IsInitialized())
	err = backend.Init()
	assert.Nil(t, err)

	assert.True(t, backend.IsInitialized())

	_, err = os.Stat(path + "/.snapshots/index")
	assert.False(t, os.IsNotExist(err))

	_, err = os.Stat(path + "/" + opts.ConfigFile)
	assert.False(t, os.IsNotExist(err))

	expected_str := []byte(path + "/.snapshots\n" + path)
	actual_str, err := ioutil.ReadFile(path + "/" + opts.ConfigFile)
	assert.Equal(t, actual_str, expected_str)
}

func TestPosixBackendGetVersions(t *testing.T) {

	//vers, err := backend.GetVersions()
	//assert.Equal(t, vers, []version{})
	//assert.Nil(t, err)

}

func TestPosixBackendCommit(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, os.Chdir(path))

	createAndSeedTestRepo(t, path, []string{})

	opts := NewOptions()
	assert.NotNil(t, opts)

	opts.RepoPath = path
	opts.SnapshotsPath = path + "/.snapshots"

	backend, err := NewPosixBackend(opts)
	assert.NotNil(t, backend)
	assert.Nil(t, err)

	err = ioutil.WriteFile(path+"/bar", []byte(""), 0644)
	assert.Nil(t, err)
	err = ioutil.WriteFile(path+"/toz", []byte(""), 0644)
	assert.Nil(t, err)

	// commit everything that is ignored or untracked
	//v, err := backend.Commit()
	//assert.Nil(t, err)
	//assert.NotNil(t, v)
}

func TestAddVersionToIndex(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, os.Chdir(path))

	v1_str := "1234567890 2012-11-01T22:08:41Z\n"
	v2_str := "5713943128 2015-02-03T23:59:12Z\n"
	v1 := NewVersion(v1_str)
	assert.NotNil(t, v1)
	v2 := NewVersion(v2_str)
	assert.NotNil(t, v2)

	err = ioutil.WriteFile(path+"/index", []byte(""), 0644)
	assert.Nil(t, err)

	err = addVersionToIndex(v1, path+"/index")
	assert.Nil(t, err)
	actual_str, err := ioutil.ReadFile(path + "/index")
	assert.Nil(t, err)
	assert.Equal(t, string(actual_str), v1_str)

	err = addVersionToIndex(v2, path+"/index")
	assert.Nil(t, err)
	actual_str, err = ioutil.ReadFile(path + "/index")
	assert.Nil(t, err)
	assert.Equal(t, string(actual_str), v1_str+v2_str)
}

func TestGetVersions(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, os.Chdir(path))

	err = os.Mkdir(path+"/.snapshots", 0755)
	assert.Nil(t, err)
	err = ioutil.WriteFile(path+"/.snapshots/index", []byte(""), 0644)
	assert.Nil(t, err)

	v1 := NewVersion("1234567890 2012-11-01T22:08:41Z")
	assert.NotNil(t, v1)
	v2 := NewVersion("5713943128 2015-02-03T23:59:12Z")
	assert.NotNil(t, v2)

	err = addVersionToIndex(v1, path+"/.snapshots/index")
	assert.Nil(t, err)
	err = addVersionToIndex(v2, path+"/.snapshots/index")
	assert.Nil(t, err)

	opts := NewOptions()
	assert.NotNil(t, opts)

	opts.RepoPath = path
	opts.SnapshotsPath = path + "/.snapshots"

	backend, err := NewPosixBackend(opts)
	assert.Nil(t, err)
	assert.NotNil(t, backend)

	//vs, err := backend.GetVersions()
	//assert.Nil(t, err)
	//assert.NotNil(t, vs)
	//assert.Equal(t, len(vs), 2)
	//assert.Equal(t, vs[0].revision, v1.revision)
	//assert.Equal(t, vs[0].timestamp, v1.timestamp)
	//assert.Equal(t, vs[1].revision, v2.revision)
	//assert.Equal(t, vs[1].timestamp, v2.timestamp)
}
