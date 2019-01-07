package filecheck_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestStatusCheck(t *testing.T) {
	os.Chmod("testdata/IsAFile.PBits0600", 0600) // force the file mode

	testCases := []struct {
		name           string
		fileName       string
		p              filecheck.Provisos
		errExpected    bool
		errMustContain []string
	}{
		{
			name:     "doesn't exist",
			fileName: "testdata/nonesuch",
			p:        filecheck.Provisos{Existence: filecheck.MustNotExist},
		},
		{
			name:     "need not exist and doesn't",
			fileName: "testdata/nonesuch",
			p:        filecheck.Provisos{Existence: filecheck.Optional},
		},
		{
			name:     "need not exist and does",
			fileName: "testdata/IsAFile",
			p:        filecheck.Provisos{Existence: filecheck.Optional},
		},
		{
			name:     "file - exists and should",
			fileName: "testdata/IsAFile",
			p:        filecheck.Provisos{Existence: filecheck.MustExist},
		},
		{
			name:        "file - exists and shouldn't",
			fileName:    "testdata/IsAFile",
			p:           filecheck.Provisos{Existence: filecheck.MustNotExist},
			errExpected: true,
		},
		{
			name:     "file (symlink) - exists and should",
			fileName: "testdata/IsASymlinkToAFile",
			p:        filecheck.Provisos{Existence: filecheck.MustExist},
		},
		{
			name:        "file (symlink) - exists - link to nothing",
			fileName:    "testdata/IsASymlinkToNothing",
			p:           filecheck.Provisos{Existence: filecheck.MustExist},
			errExpected: true,
		},
		{
			name:     "file (symlink) - exists - link to nothing - dont follow",
			fileName: "testdata/IsASymlinkToNothing",
			p: filecheck.Provisos{
				Existence:          filecheck.MustExist,
				DontFollowSymlinks: true,
			},
		},
		{
			name:     "file - perms equal 0600",
			fileName: "testdata/IsAFile.PBits0600",
			p: filecheck.Provisos{
				Existence: filecheck.MustExist,
				Checks: []check.FileInfo{
					check.FileInfoPerm(check.FilePermEQ(0600)),
				},
			},
		},
		{
			name:     "file - perms don't equal 0664",
			fileName: "testdata/IsAFile.PBits0600",
			p: filecheck.Provisos{
				Existence: filecheck.MustExist,
				Checks: []check.FileInfo{
					check.FileInfoPerm(check.FilePermEQ(0644)),
				},
			},
			errExpected: true,
		},
	}

	for i, tc := range testCases {
		testID := fmt.Sprintf("test %d: %s", i, tc.name)
		err := tc.p.StatusCheck(tc.fileName)
		testhelper.CheckError(t, testID, err, tc.errExpected, []string{})
	}

}

func TestESToString(t *testing.T) {
	testCases := []struct {
		name   string
		p      filecheck.Provisos
		expVal string
	}{
		{
			name:   "must not exist, no checks",
			p:      filecheck.Provisos{Existence: filecheck.MustNotExist},
			expVal: "The filesystem object must not exist",
		},
		{
			name: "must not exist, with (redundant) checks",
			p: filecheck.Provisos{
				Existence: filecheck.MustNotExist,
				Checks: []check.FileInfo{
					check.FileInfoSize(check.Int64EQ(0)),
				},
			},
			expVal: "The filesystem object must not exist",
		},
		{
			name:   "must exist, no checks",
			p:      filecheck.Provisos{Existence: filecheck.MustExist},
			expVal: "The filesystem object must exist",
		},
		{
			name: "must exist, with checks",
			p: filecheck.Provisos{
				Existence: filecheck.MustExist,
				Checks: []check.FileInfo{
					check.FileInfoSize(check.Int64EQ(0)),
				},
			},
			expVal: "The filesystem object must exist" +
				" and must satisfy further checks",
		},
		{
			name:   "need not exist, no checks",
			p:      filecheck.Provisos{Existence: filecheck.Optional},
			expVal: "The filesystem object need not exist",
		},
		{
			name: "need not exist, with checks",
			p: filecheck.Provisos{
				Existence: filecheck.Optional,
				Checks: []check.FileInfo{
					check.FileInfoSize(check.Int64EQ(0)),
				},
			},
			expVal: "The filesystem object need not exist" +
				" but if it does it must satisfy further checks",
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s :", i, tc.name)
		val := tc.p.String()
		if val != tc.expVal {
			t.Log(tcID)
			t.Logf("\t: Expected: %s\n", tc.expVal)
			t.Logf("\t:      Got: %s\n", val)
			t.Errorf("\t: bad string representation of the Provisos\n")
		}
	}
}
