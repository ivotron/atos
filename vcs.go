package vio

import (
	"os"
	"os/exec"
	"strings"
)

func runCmd(path string, cmdstr string) (out string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	err = os.Chdir(path)
	if err != nil {
		return
	}
	if len(strings.TrimSpace(cmdstr)) == 0 {
		err = AnError{"Empty command"}
		return
	}
	cmd_items := strings.Split(cmdstr, " ")
	cmd := cmd_items[0]
	args := []string{}
	if len(cmd_items) > 1 {
		args = cmd_items[1:]
	}
	o, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return
	}
	out = string(o)
	err = os.Chdir(cwd)
	return
}

func HasUncommittedChanges(repoPath string) (has bool, err error) {
	out, err := runCmd(repoPath, "git status -s -uno")
	if err != nil {
		return
	}
	has = strings.TrimSpace(out) != ""
	return
}

func GetVersionedFiles(repoPath string) (versioned []string, err error) {
	out, err := runCmd(repoPath, "git ls-files")
	if err != nil {
		return
	}
	versioned = strings.Split(strings.TrimSpace(out), "\n")
	return
}

func GetCurrentCommitId(repoPath string) (id string, err error) {
	out, err := runCmd(repoPath, "git rev-parse --verify --short HEAD")
	if err != nil {
		return
	}
	id = strings.TrimSpace(out)
	return
}
