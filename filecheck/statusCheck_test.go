package filecheck_test

import (
	"os"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestStatusCheck(t *testing.T) {
	_ = os.Chmod("testdata/IsAFile.PBits0600", 0600) // force the file mode
	const noSuchFile = "testdata/nonesuch"
	const isAFile = "testdata/IsAFile"
	const symlinkToAFile = "testdata/IsASymlinkToAFile"
	const symlinkToNothing = "testdata/IsASymlinkToNothing"

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		fileName string
		p        filecheck.Provisos
	}{
		{
			ID:       testhelper.MkID("doesn't exist but should"),
			fileName: noSuchFile,
			p:        filecheck.Provisos{Existence: filecheck.MustExist},
			ExpErr: testhelper.MkExpErr(noSuchFile,
				"does not exist but should"),
		},
		{
			ID:       testhelper.MkID("doesn't exist"),
			fileName: noSuchFile,
			p:        filecheck.Provisos{Existence: filecheck.MustNotExist},
		},
		{
			ID:       testhelper.MkID("need not exist and doesn't"),
			fileName: noSuchFile,
			p:        filecheck.Provisos{Existence: filecheck.Optional},
		},
		{
			ID:       testhelper.MkID("need not exist and does"),
			fileName: isAFile,
			p:        filecheck.Provisos{Existence: filecheck.Optional},
		},
		{
			ID:       testhelper.MkID("file - exists and should"),
			fileName: isAFile,
			p:        filecheck.Provisos{Existence: filecheck.MustExist},
		},
		{
			ID:       testhelper.MkID("file - exists and shouldn't"),
			fileName: isAFile,
			p:        filecheck.Provisos{Existence: filecheck.MustNotExist},
			ExpErr:   testhelper.MkExpErr(isAFile, "exists but shouldn't"),
		},
		{
			ID:       testhelper.MkID("file (symlink) - exists and should"),
			fileName: symlinkToAFile,
			p:        filecheck.Provisos{Existence: filecheck.MustExist},
		},
		{
			ID:       testhelper.MkID("file (symlink) - exists - link to nothing"),
			fileName: symlinkToNothing,
			p:        filecheck.Provisos{Existence: filecheck.MustExist},
			ExpErr: testhelper.MkExpErr(symlinkToNothing,
				"does not exist but should"),
		},
		{
			ID: testhelper.MkID(
				"file (symlink) - exists - link to nothing - dont follow"),
			fileName: "testdata/IsASymlinkToNothing",
			p: filecheck.Provisos{
				Existence:          filecheck.MustExist,
				DontFollowSymlinks: true,
			},
		},
		{
			ID:       testhelper.MkID("file - perms equal 0600"),
			fileName: "testdata/IsAFile.PBits0600",
			p: filecheck.Provisos{
				Existence: filecheck.MustExist,
				Checks: []check.FileInfo{
					check.FileInfoPerm(check.FilePermEQ(0600)),
				},
			},
		},
		{
			ID:       testhelper.MkID("file - perms don't equal 0664"),
			fileName: "testdata/IsAFile.PBits0600",
			p: filecheck.Provisos{
				Existence: filecheck.MustExist,
				Checks: []check.FileInfo{
					check.FileInfoPerm(check.FilePermEQ(0644)),
				},
			},
			ExpErr: testhelper.MkExpErr("the check on the permissions of",
				"IsAFile.PBits0600",
				"failed: the permissions (0600) should be equal to 0644"),
		},
	}

	for _, tc := range testCases {
		err := tc.p.StatusCheck(tc.fileName)
		testhelper.CheckExpErr(t, err, tc)
	}

}

func TestESToString(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		p      filecheck.Provisos
		expVal string
	}{
		{
			ID:     testhelper.MkID("must not exist, no checks"),
			p:      filecheck.Provisos{Existence: filecheck.MustNotExist},
			expVal: "The filesystem object must not exist",
		},
		{
			ID: testhelper.MkID("must not exist, with (redundant) checks"),
			p: filecheck.Provisos{
				Existence: filecheck.MustNotExist,
				Checks: []check.FileInfo{
					check.FileInfoSize(check.Int64EQ(0)),
				},
			},
			expVal: "The filesystem object must not exist",
		},
		{
			ID:     testhelper.MkID("must exist, no checks"),
			p:      filecheck.Provisos{Existence: filecheck.MustExist},
			expVal: "The filesystem object must exist",
		},
		{
			ID: testhelper.MkID("must exist, with checks"),
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
			ID:     testhelper.MkID("need not exist, no checks"),
			p:      filecheck.Provisos{Existence: filecheck.Optional},
			expVal: "The filesystem object need not exist",
		},
		{
			ID: testhelper.MkID("need not exist, with checks"),
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

	for _, tc := range testCases {
		val := tc.p.String()
		if val != tc.expVal {
			t.Log(tc.IDStr())
			t.Logf("\t: Expected: %s\n", tc.expVal)
			t.Logf("\t:      Got: %s\n", val)
			t.Errorf("\t: bad string representation of the Provisos\n")
		}
	}
}
