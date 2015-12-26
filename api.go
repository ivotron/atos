package vio

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Status int

const (
	Committed Status = iota
	Staged
)

type version struct {
	revision  string
	timestamp time.Time
}

func NewVersion(revision string) *version {
	fields := strings.Fields(revision)

	if len(fields) == 2 {
		ts, err := time.Parse(time.RFC3339, fields[1])
		if err != nil {
			panic(err)
		}
		return &version{revision: fields[0], timestamp: ts}
	} else {
		return &version{revision: fields[0], timestamp: time.Now()}
	}
}

func ContainsVersion(vs []version, v *version) bool {
	for _, v_in := range vs {
		if reflect.DeepEqual(*v, v_in) {
			return true
		}
	}
	return false
}

func (v *version) String() string {
	return fmt.Sprintf("%s %s", v.revision, v.timestamp.Format(time.RFC3339))
}

type Backend interface {
	// inits backend in current directory
	Init() error

	// opens the backend
	Open() error

	// whether the backend has been initialized
	IsInitialized() bool

	// returns status of backend
	GetStatus() (Status, error)

	// checks out a commit
	Checkout(v *version) error

	// commits a version.
	Commit() (*version, error)

	// retrieves the string representation of the diff for an object
	Diff(v1 *version, v2 *version, obj string) (string, error)

	// return the list
	GetVersions() (versions []version, err error)
}

type AnError struct {
	Msg string
}

func (e AnError) Error() string {
	return e.Msg
}
