package vio

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/tgulacsi/go-locking"
)

type PosixBackend struct {
	snapshotsPath string
	repoPath      string
	configFile    string
}

func NewPosixBackend(o Options) (b *PosixBackend, err error) {
	return &PosixBackend{
		snapshotsPath: o.SnapshotsPath,
		repoPath:      o.RepoPath,
		configFile:    o.ConfigFile}, nil
}

func (b PosixBackend) Init() (err error) {
	if err = os.Mkdir(b.snapshotsPath, 0755); err != nil {
		return
	}

	// TODO: check if /index already exists

	if err = ioutil.WriteFile(b.snapshotsPath+"/index", []byte(""), 0644); err != nil {
		return
	}

	str := []byte(b.snapshotsPath + "\n")

	// TODO: use gcfg instead
	if err = ioutil.WriteFile(b.repoPath+"/"+b.configFile, str, 0644); err != nil {
		return
	}

	return
}
func (b PosixBackend) Open() error {
	return nil
}

func (b PosixBackend) IsInitialized() bool {
	_, err := os.Stat(b.snapshotsPath + "/index")
	return os.IsNotExist(err) == false
}
func (b PosixBackend) GetStatus() (Status, error) {
	return Committed, nil
}
func (b PosixBackend) Checkout(v *version) error {
	return AnError{"not yet"}
}

func (b PosixBackend) Commit() (v *version, err error) {
	has, err := HasUncommittedChanges(b.repoPath)

	if err != nil {
		return
	}
	if has {
		return nil, AnError{"Uncommitted changes in repo."}
	}

	versionedFiles, err := GetVersionedFiles(b.repoPath)
	if err != nil {
		return
	}

	id, err := GetCurrentCommitId(b.repoPath)
	if err != nil {
		return
	}

	v = NewVersion(id)

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
	args = append(args, "-az")
	for _, vfile := range versionedFiles {
		args = append(args, "--exclude="+vfile)
	}

	args = append(args, "--exclude=.git/")

	if _, err := os.Stat(repoPath + "/.vioignore"); os.IsNotExist(err) == false {
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
		v := *NewVersion(line)
		versions = append(versions, v)
	}
	return
}

func (b PosixBackend) Diff(v1 *version, v2 *version, obj string) (string, error) {
	return "", AnError{"not yet"}
}
