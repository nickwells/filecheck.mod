package filecheck

import (
	"fmt"
	"os"

	"github.com/nickwells/check.mod/check"
)

// Exists records whether the file-system object should exist or
// not. In each case the check is only valid at the time the check is
// made and so any code using this should be aware of this
type Exists uint

const (
	// Optional indicates that no existence check should be made
	Optional Exists = iota
	// MustExist indicates that the object must exist
	MustExist
	// MustNotExist indicates that the object must not exist
	MustNotExist
)

// Provisos records the expectations of the file-system object. You can
// specify whether the object should or shouldn't exist and whether, if it's
// a symlink, you should follow the link. You can also perform a number of
// checks on the status of the file system object.
type Provisos struct {
	Existence          Exists
	Checks             []check.FileInfo
	DontFollowSymlinks bool
}

// StatusCheck checks that the file system object called 'name' satisfies
// the constraints. It returns a non-nil error if the constraint is not
// met. Note that if the file does not exist and it is not expected to
// exist then no further checks are performed (this may be obvious to you)
func (p Provisos) StatusCheck(name string) error {
	var info os.FileInfo
	var err error
	if p.DontFollowSymlinks {
		info, err = os.Lstat(name)
	} else {
		info, err = os.Stat(name)
	}

	if os.IsNotExist(err) {
		if p.Existence == MustExist {
			return fmt.Errorf("path: %q should exist but doesn't", name)
		}
		return nil
	}

	if p.Existence == MustNotExist {
		return fmt.Errorf("path: %q shouldn't exist but does", name)
	}

	if err != nil {
		return fmt.Errorf("path: %q error: %s", name, err.Error())
	}

	for _, c := range p.Checks {
		if err := c(info); err != nil {
			return err
		}
	}
	return nil
}

// String returns a string describing the Provisos
func (p Provisos) String() string {
	rval := ""
	prefix := ""

	switch p.Existence {
	case MustNotExist:
		return "The filesystem object must not exist"
	case MustExist:
		rval = "The filesystem object must exist"
		prefix = " and"
	case Optional:
		rval = "The filesystem object need not exist"
		prefix = " but if it does it"
	}

	if len(p.Checks) > 0 {
		rval += prefix + " must satisfy further checks"
	}

	return rval
}
