package vio

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPosixBackendInitWithUninitializedRepo(t *testing.T) {
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

	expected_str := []byte(path + "/.snapshots\n")
	actual_str, err := ioutil.ReadFile(path + "/" + opts.ConfigFile)
	assert.Equal(t, actual_str, expected_str)
}

func TestPosixBackendInitWithInitializedRepo(t *testing.T) {
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

	assert.False(t, backend.IsInitialized())
	err = backend.Init()
	assert.Nil(t, err)

	assert.True(t, backend.IsInitialized())

	_, err = os.Stat(path + "/.snapshots/index")
	assert.False(t, os.IsNotExist(err))

	_, err = os.Stat(path + "/" + opts.ConfigFile)
	assert.False(t, os.IsNotExist(err))

	expected_str := []byte(path + "/.snapshots\n")
	actual_str, err := ioutil.ReadFile(path + "/" + opts.ConfigFile)
	assert.Equal(t, actual_str, expected_str)
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

	err = backend.Init()
	assert.Nil(t, err)

	err = ioutil.WriteFile(path+"/bar", []byte("yeah"), 0644)
	assert.Nil(t, err)
	err = ioutil.WriteFile(path+"/toz", []byte("ok"), 0644)
	assert.Nil(t, err)
	_, err = runCmd(path, "git add .vioconfig")
	assert.Nil(t, err)
	_, err = runCmd(path, "git commit -m committing_vioconfig_file")
	assert.Nil(t, err)

	// commit everything that is ignored or untracked
	v, err := backend.Commit()
	assert.Nil(t, err)
	assert.NotNil(t, v)

	snapPath := path + "/.snapshots/" + v.revision
	_, err = os.Stat(snapPath)
	assert.False(t, os.IsNotExist(err))

	unixTime := fmt.Sprintf("%d", v.timestamp.Unix())
	_, err = os.Stat(snapPath + "/" + unixTime + "/toz")
	assert.False(t, os.IsNotExist(err))
	_, err = os.Stat(snapPath + "/" + unixTime + "/bar")
	assert.False(t, os.IsNotExist(err))

	_, err = os.Stat(snapPath + "/" + unixTime + "/.git")
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(snapPath + "/" + unixTime + "/.vioconfig")
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(snapPath + "/" + unixTime + "/README")
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(snapPath + "/" + unixTime + "/.gitignore")
	assert.True(t, os.IsNotExist(err))
}

func TestPosixBackendCommitWithIgnore(t *testing.T) {
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

	err = backend.Init()
	assert.Nil(t, err)

	err = ioutil.WriteFile(path+"/bar", []byte("yeah"), 0644)
	assert.Nil(t, err)
	err = ioutil.WriteFile(path+"/toz", []byte("ok"), 0644)
	assert.Nil(t, err)
	err = ioutil.WriteFile(path+"/.vioignore", []byte("ignored_folder\n.snapshots\n"), 0644)
	assert.Nil(t, err)
	err = os.Mkdir(path+"/ignored_folder", 0755)
	assert.Nil(t, err)
	err = ioutil.WriteFile(path+"/ignored_folder/foo", []byte("ignore this\n"), 0644)
	assert.Nil(t, err)
	_, err = runCmd(path, "git add .vioconfig .vioignore")
	assert.Nil(t, err)
	_, err = runCmd(path, "git commit -m committing_vio_files")
	assert.Nil(t, err)

	v, err := backend.Commit()
	assert.Nil(t, err)
	assert.NotNil(t, v)

	snapPath := path + "/.snapshots/" + v.revision
	_, err = os.Stat(snapPath)
	assert.False(t, os.IsNotExist(err))

	unixTime := fmt.Sprintf("%d", v.timestamp.Unix())
	_, err = os.Stat(snapPath + "/" + unixTime + "/toz")
	assert.False(t, os.IsNotExist(err))
	_, err = os.Stat(snapPath + "/" + unixTime + "/bar")
	assert.False(t, os.IsNotExist(err))

	_, err = os.Stat(snapPath + "/" + unixTime + "/ignored_folder")
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(snapPath + "/" + unixTime + "/.git")
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(snapPath + "/" + unixTime + "/.vioconfig")
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(snapPath + "/" + unixTime + "/README")
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(snapPath + "/" + unixTime + "/.gitignore")
	assert.True(t, os.IsNotExist(err))
}

func TestPosixBackendAddVersionToIndex(t *testing.T) {
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

func TestPosixBackendGetVersions(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, os.Chdir(path))

	opts := NewOptions()
	assert.NotNil(t, opts)

	opts.RepoPath = path
	opts.SnapshotsPath = path + "/.snapshots"

	backend, err := NewPosixBackend(opts)

	err = os.Mkdir(path+"/.snapshots", 0755)
	assert.Nil(t, err)
	err = ioutil.WriteFile(path+"/.snapshots/index", []byte(""), 0644)
	assert.Nil(t, err)

	vs, err := backend.GetVersions()
	assert.Nil(t, err)
	assert.NotNil(t, vs)
	assert.Equal(t, len(vs), 0)

	v1 := NewVersion("1234567890 2012-11-01T22:08:41Z")
	assert.NotNil(t, v1)
	v2 := NewVersion("5713943128 2015-02-03T23:59:12Z")
	assert.NotNil(t, v2)

	err = addVersionToIndex(v1, path+"/.snapshots/index")
	assert.Nil(t, err)
	vs, err = backend.GetVersions()
	assert.Nil(t, err)
	assert.NotNil(t, vs)
	assert.Equal(t, len(vs), 1)
	assert.Equal(t, vs[0].revision, v1.revision)
	assert.Equal(t, vs[0].timestamp, v1.timestamp)

	err = addVersionToIndex(v2, path+"/.snapshots/index")
	assert.Nil(t, err)

	assert.Nil(t, err)
	assert.NotNil(t, backend)

	vs, err = backend.GetVersions()
	assert.Nil(t, err)
	assert.NotNil(t, vs)
	assert.Equal(t, len(vs), 2)
	assert.Equal(t, vs[0].revision, v1.revision)
	assert.Equal(t, vs[0].timestamp, v1.timestamp)
	assert.Equal(t, vs[1].revision, v2.revision)
	assert.Equal(t, vs[1].timestamp, v2.timestamp)
}
