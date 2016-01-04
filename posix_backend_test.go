package vio

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"gopkg.in/ini.v1"

	"github.com/stretchr/testify/assert"
)

func getNewPosixBackend(t *testing.T, path string) (b Backend) {
	opts := ini.Empty()
	assert.NotNil(t, opts)

	opts.Section("").Key("repo_path").SetValue(path)
	opts.Section("").Key("snapshots_path").SetValue(path + "/.snapshots")
	opts.Section("").Key("backend_type").SetValue("posix")

	b, err := InstantiateBackend(opts)
	assert.NotNil(t, b)
	assert.Nil(t, err)

	return
}

func TestPosixBackendInitWithUninitializedRepo(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, os.Chdir(path))

	backend := getNewPosixBackend(t, path)

	assert.False(t, backend.IsInitialized())
	err = backend.Init()
	assert.Nil(t, err)

	assert.True(t, backend.IsInitialized())

	_, err = os.Stat(path + "/.snapshots/index")
	assert.Nil(t, err)
}

func TestPosixBackendInitWithInitializedRepo(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, os.Chdir(path))

	createAndSeedTestRepo(t, path, []string{})

	backend := getNewPosixBackend(t, path)

	assert.False(t, backend.IsInitialized())
	err = backend.Init()
	assert.Nil(t, err)

	assert.True(t, backend.IsInitialized())

	_, err = os.Stat(path + "/.snapshots/index")
	assert.Nil(t, err)
}

func TestPosixBackendCommit(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, os.Chdir(path))

	createAndSeedTestRepo(t, path, []string{})

	backend := getNewPosixBackend(t, path)

	err = backend.Init()
	assert.Nil(t, err)

	err = ioutil.WriteFile(path+"/bar", []byte("yeah"), 0644)
	assert.Nil(t, err)
	err = ioutil.WriteFile(path+"/toz", []byte("ok"), 0644)
	assert.Nil(t, err)

	// commit everything that is ignored or untracked
	v, err := backend.Commit(map[string]string{})
	assert.Nil(t, err)
	assert.NotNil(t, v)

	snapPath := path + "/.snapshots/" + v.revision
	_, err = os.Stat(snapPath)
	assert.Nil(t, err)

	unixTime := fmt.Sprintf("%d", v.timestamp.Unix())
	_, err = os.Stat(snapPath + "/" + unixTime + "/toz")
	assert.Nil(t, err)
	_, err = os.Stat(snapPath + "/" + unixTime + "/bar")
	assert.Nil(t, err)

	_, err = os.Stat(snapPath + "/" + unixTime + "/.git")
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

	backend := getNewPosixBackend(t, path)

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
	_, err = runCmd(path, "git add .vioignore")
	assert.Nil(t, err)
	_, err = runCmd(path, "git commit -m committing_vio_ignored_files")
	assert.Nil(t, err)

	v, err := backend.Commit(map[string]string{})
	assert.Nil(t, err)
	assert.NotNil(t, v)

	snapPath := path + "/.snapshots/" + v.revision
	_, err = os.Stat(snapPath)
	assert.Nil(t, err)

	unixTime := fmt.Sprintf("%d", v.timestamp.Unix())
	_, err = os.Stat(snapPath + "/" + unixTime + "/toz")
	assert.Nil(t, err)
	_, err = os.Stat(snapPath + "/" + unixTime + "/bar")
	assert.Nil(t, err)

	_, err = os.Stat(snapPath + "/" + unixTime + "/ignored_folder")
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(snapPath + "/" + unixTime + "/.git")
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(snapPath + "/" + unixTime + "/README")
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(snapPath + "/" + unixTime + "/.gitignore")
	assert.True(t, os.IsNotExist(err))
}

func TestPosixBackendCheckout(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, os.Chdir(path))

	createAndSeedTestRepo(t, path, []string{})

	fmt.Println("path: %s\n" + path)

	backend := getNewPosixBackend(t, path)

	err = backend.Init()
	assert.Nil(t, err)

	err = ioutil.WriteFile(path+"/bar", []byte("yeah"), 0644)
	assert.Nil(t, err)
	err = ioutil.WriteFile(path+"/toz", []byte("ok"), 0644)
	assert.Nil(t, err)

	// commit everything that is ignored or untracked
	v, err := backend.Commit(map[string]string{})
	assert.Nil(t, err)
	assert.NotNil(t, v)

	err = os.Remove(path + "/bar")
	assert.Nil(t, err)
	err = os.Remove(path + "/toz")
	assert.Nil(t, err)

	// v = NewVersion(fmt.Sprintf("%s#%d", v.revision, v.timestamp.Unix()))

	err = backend.Checkout(v)
	assert.Nil(t, err)

	_, err = os.Stat(path + "/bar")
	assert.Nil(t, err)
	_, err = os.Stat(path + "/toz")
	assert.Nil(t, err)
}

func TestPosixBackendCommitWithTags(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, os.Chdir(path))

	createAndSeedTestRepo(t, path, []string{})

	backend := getNewPosixBackend(t, path)

	err = backend.Init()
	assert.Nil(t, err)

	err = ioutil.WriteFile(path+"/bar", []byte("yeah"), 0644)
	assert.Nil(t, err)
	err = ioutil.WriteFile(path+"/toz", []byte("ok"), 0644)
	assert.Nil(t, err)

	// commit everything that is ignored or untracked
	v, err := backend.Commit(map[string]string{"foo": "bar", "hello": "goodbye"})
	assert.Nil(t, err)
	assert.NotNil(t, v)

	contents, err := ioutil.ReadFile(path + "/.snapshots/index")
	assert.Nil(t, err)
	assert.Equal(t, strings.TrimSpace(string(contents)), fmt.Sprintf("%v", v))
}

func TestPosixBackendAddVersionToIndex(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, os.Chdir(path))

	v1_str := "1234567890#1405544146"
	v2_str := "5713943128#2435869343"
	v1 := NewVersion(v1_str)
	assert.NotNil(t, v1)
	v2 := NewVersion(v2_str)
	assert.NotNil(t, v2)
	meta := map[string]string{"foo": "bar", "hello": "goodbye"}
	v3_str := "3943943128#5635869343"
	v3 := NewVersionWithMeta(v3_str, meta)
	assert.NotNil(t, v3)

	err = ioutil.WriteFile(path+"/index", []byte(""), 0644)
	assert.Nil(t, err)

	err = addVersionToIndex(v1, path+"/index")
	assert.Nil(t, err)
	actual_str, err := ioutil.ReadFile(path + "/index")
	assert.Nil(t, err)
	assert.Equal(t, string(actual_str), v1_str+",{}\n")

	err = addVersionToIndex(v2, path+"/index")
	assert.Nil(t, err)
	actual_str, err = ioutil.ReadFile(path + "/index")
	assert.Nil(t, err)
	assert.Equal(t, string(actual_str), v1_str+",{}\n"+v2_str+",{}\n")

	err = addVersionToIndex(v3, path+"/index")
	assert.Nil(t, err)
	actual_str, err = ioutil.ReadFile(path + "/index")
	assert.Nil(t, err)
	assert.Equal(t, string(actual_str),
		v1_str+",{}\n"+v2_str+",{}\n"+v3_str+","+"{\"foo\":\"bar\",\"hello\":\"goodbye\"}\n")
}

func TestPosixBackendGetVersions(t *testing.T) {
	path, err := ioutil.TempDir("", "testing")
	assert.Nil(t, os.Chdir(path))

	backend := getNewPosixBackend(t, path)

	err = os.Mkdir(path+"/.snapshots", 0755)
	assert.Nil(t, err)
	err = ioutil.WriteFile(path+"/.snapshots/index", []byte(""), 0644)
	assert.Nil(t, err)

	vs, err := backend.GetVersions()
	assert.Nil(t, err)
	assert.NotNil(t, vs)
	assert.Equal(t, len(vs), 0)

	v1 := NewVersion("1234567890#1405544146")
	assert.NotNil(t, v1)
	v2 := NewVersion("5713943128#2435869343")
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
