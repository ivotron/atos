package vio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/ini.v1"

	"github.com/tgulacsi/go-locking"
)

type PosixBackend struct {
	snapshotsPath string
	repoPath      string
}

func NewPosixBackend(o *ini.File) (b Backend, err error) {
	if !o.Section("").HasKey("snapshots_path") {
		return nil, AnError{"Expecting key 'snapshots_path' in configuration."}
	}
	if !o.Section("").HasKey("repo_path") {
		return nil, AnError{"Expecting key 'repo_path' in configuration."}
	}
	return &PosixBackend{
		snapshotsPath: o.Section("").Key("snapshots_path").String(),
		repoPath:      o.Section("").Key("repo_path").String()}, nil
}

func (b PosixBackend) Init() (err error) {
	if err = os.Mkdir(b.snapshotsPath, 0755); err != nil {
		return
	}

	if _, err = os.Stat(b.snapshotsPath + "/index"); err == nil {
		return AnError{"Repository already initialized"}
	}

	if err = ioutil.WriteFile(b.snapshotsPath+"/index", []byte(""), 0644); err != nil {
		return
	}

	return
}

func (b PosixBackend) isRepoOK() (err error) {
	if !b.IsInitialized() {
		return AnError{"Uninitialized repository."}
	}

	hasUncommitted, err := HasUncommittedChanges(b.repoPath)

	if err != nil {
		return
	}
	if hasUncommitted {
		return AnError{"Uncommitted changes in repo."}
	}

	return
}

func (b PosixBackend) Open() error {
	return nil
}

func (b PosixBackend) IsInitialized() bool {
	_, err := os.Stat(b.snapshotsPath + "/index")
	return err == nil
}

func (b PosixBackend) GetStatus() (Status, error) {
	return Committed, nil
}

func (b PosixBackend) Checkout(v *version) (err error) {
	if err = b.isRepoOK(); err != nil {
		return
	}

	// acquire a lock on the index file
	flock, err := locking.NewFLock(b.snapshotsPath + "/index")
	if err != nil {
		return
	}
	if err = flock.Lock(); err != nil {
		return
	}
	defer flock.Unlock()

	idx, err := b.GetVersions()
	if err != nil {
		return
	}

	if !ContainsVersion(idx, v) {
		return AnError{
			fmt.Sprintf("Version %s#%d not in index", v.revision, v.timestamp.Unix())}
	}

	return checkoutSnapshot(b.repoPath, b.snapshotsPath, v)
}

func (b PosixBackend) Commit(meta map[string]string) (v *version, err error) {
	if err = b.isRepoOK(); err != nil {
		return
	}
	versionedFiles, err := GetVersionedFiles(b.repoPath)
	if err != nil {
		return
	}

	id, err := GetCurrentCommitId(b.repoPath)
	if err != nil {
		return
	}

	v = NewVersionWithMeta(id, meta)

	// acquire a lock on the index file
	flock, err := locking.NewFLock(b.snapshotsPath + "/index")
	if err != nil {
		return
	}
	if err = flock.Lock(); err != nil {
		return
	}
	defer flock.Unlock()

	idx, err := b.GetVersions()
	if err != nil {
		return
	}
	if ContainsVersion(idx, v) {
		return nil, AnError{"Version " + fmt.Sprintf("%v", v) + " already in index."}
	}

	if err = createSnapshot(b.repoPath, b.snapshotsPath, v, versionedFiles); err != nil {
		return
	}

	if err = addVersionToIndex(v, b.snapshotsPath+"/index"); err != nil {
		return
	}

	return
}

func checkoutSnapshot(repoPath string, snapsPath string, v *version) (err error) {

	if _, err = os.Stat(snapsPath + "/" + v.revision); err != nil {
		return
	}

	unixTime := fmt.Sprintf("%d", v.timestamp.Unix())
	if _, err = os.Stat(snapsPath + "/" + v.revision + "/" + unixTime); err != nil {
		return
	}

	srcPath := snapsPath + "/" + v.revision + "/" + unixTime + "/"

	var args []string
	args = append(args, "-a")

	// source
	args = append(args, srcPath)

	// destination
	args = append(args, repoPath)

	_, err = exec.Command("rsync", args...).CombinedOutput()

	return
}

func createSnapshot(repoPath string,
	snapsPath string, v *version, versionedFiles []string) (err error) {

	if err = os.MkdirAll(snapsPath+"/"+v.revision, 0755); err != nil {
		return
	}

	unixTime := fmt.Sprintf("%d", v.timestamp.Unix())
	if err = os.Mkdir(snapsPath+"/"+v.revision+"/"+unixTime, 0755); err != nil {
		return
	}
	destPath := snapsPath + "/" + v.revision + "/" + unixTime

	var args []string
	args = append(args, "-a")
	for _, vfile := range versionedFiles {
		args = append(args, "--exclude="+vfile)
	}

	args = append(args, "--exclude=.git/")

	if _, err := os.Stat(repoPath + "/.vioignore"); err == nil {
		args = append(args, "--filter=:-_/.vioignore")
	}

	// source
	args = append(args, repoPath+"/")

	// destination
	args = append(args, destPath)

	_, err = exec.Command("rsync", args...).CombinedOutput()

	return
}

func addVersionToIndex(v *version, filename string) (err error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}

	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%v\n", v))

	return
}

func (b PosixBackend) GetVersions() (versions []version, err error) {
	contents, err := ioutil.ReadFile(b.snapshotsPath + "/index")
	if err != nil {
		return
	}

	versions = []version{}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		i := strings.Index(line, ",")
		if i < 0 {
			return nil, AnError{"Malformed version in index: " + line}
		}
		v_str := line[:i]
		meta_str := line[i+1:]

		var meta map[string]string

		err = json.Unmarshal([]byte(meta_str), &meta)
		if err != nil {
			return
		}

		v := *NewVersionWithMeta(v_str, meta)
		versions = append(versions, v)
	}
	return
}

func (b PosixBackend) Diff(v1 *version, v2 *version, obj string) (string, error) {
	return "", AnError{"not yet"}
}
