package filecheck_test

import (
	"os"
	"testing"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestStatusCheck(t *testing.T) {
	const noSuchFile = "testdata/nonesuch"
	const isAFile = "testdata/IsAFile"
	const isAFile600 = "testdata/IsAFile.PBits0600"
	const symlinkToAFile = "testdata/IsASymlinkToAFile"
	const symlinkToNothing = "testdata/IsASymlinkToNothing"

	err := os.Chmod(isAFile600, 0o600) // force the file mode
	if err != nil {
		t.Fatalf("Cannot set the file permission bits on %q: %s\n",
			isAFile600, err)
	}

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		fileName string
		p        filecheck.Provisos
	}{
		{
			ID:       testhelper.MkID("does not exist but should"),
			fileName: noSuchFile,
			p:        filecheck.Provisos{Existence: filecheck.MustExist},
			ExpErr: testhelper.MkExpErr(noSuchFile,
				"should exist but does not"),
		},
		{
			ID:       testhelper.MkID("does not exist"),
			fileName: noSuchFile,
			p:        filecheck.Provisos{Existence: filecheck.MustNotExist},
		},
		{
			ID:       testhelper.MkID("need not exist and does not"),
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
			ID:       testhelper.MkID("file - exists and should not"),
			fileName: isAFile,
			p:        filecheck.Provisos{Existence: filecheck.MustNotExist},
			ExpErr:   testhelper.MkExpErr(isAFile, "should not exist but does"),
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
				"should exist but does not"),
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
			fileName: isAFile600,
			p: filecheck.Provisos{
				Existence: filecheck.MustExist,
				Checks: []check.FileInfo{
					check.FileInfoPerm(check.FilePermEQ(0o600)),
				},
			},
		},
		{
			ID:       testhelper.MkID("file - perms do not equal 0664"),
			fileName: isAFile600,
			p: filecheck.Provisos{
				Existence: filecheck.MustExist,
				Checks: []check.FileInfo{
					check.FileInfoPerm(check.FilePermEQ(0o644)),
				},
			},
			ExpErr: testhelper.MkExpErr("the file permissions of",
				"IsAFile.PBits0600",
				"incorrect: the permissions (0600) should equal 0644"),
		},
	}

	for _, tc := range testCases {
		err := tc.p.StatusCheck(tc.fileName)
		testhelper.CheckExpErr(t, err, tc)
	}
}

func TestProvisosToString(t *testing.T) {
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
					check.FileInfoSize(check.ValEQ[int64](0)),
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
					check.FileInfoSize(check.ValEQ[int64](0)),
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
					check.FileInfoSize(check.ValEQ[int64](0)),
				},
			},
			expVal: "The filesystem object need not exist" +
				" but if it does it must satisfy further checks",
		},
	}

	for _, tc := range testCases {
		testhelper.DiffString(t, tc.IDStr(), "proviso string",
			tc.p.String(), tc.expVal)
	}
}
