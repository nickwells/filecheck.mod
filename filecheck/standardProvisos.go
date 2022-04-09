package filecheck

import (
	"github.com/nickwells/check.mod/v2/check"
)

// DirExists returns a Provisos that will check that the value refers
// to a directory, which must exist. This is a common Provisos value and this
// func is provided to simplify your code.
func DirExists() Provisos {
	return Provisos{
		Existence: MustExist,
		Checks:    []check.FileInfo{check.FileInfoIsDir},
	}
}

// FileExists returns a Provisos that will check that the value refers to a
// regular file, which must exist. This is a common Provisos value and this
// func is provided to simplify your code.
func FileExists() Provisos {
	return Provisos{
		Existence: MustExist,
		Checks:    []check.FileInfo{check.FileInfoIsRegular},
	}
}

// FileNonEmpty returns a Provisos that will check that the value refers to a
// regular file, which must exist and must not be empty. This is a common
// Provisos value and this func is provided to simplify your code.
func FileNonEmpty() Provisos {
	return Provisos{
		Existence: MustExist,
		Checks: []check.FileInfo{
			check.FileInfoIsRegular,
			check.FileInfoSize(check.ValGT[int64](0)),
		},
	}
}

// IsNew returns a Provisos that will check that the value refers to a
// non-existent file/directory. This is a common Provisos value and this func
// is provided to simplify your code.
func IsNew() Provisos {
	return Provisos{Existence: MustNotExist}
}
