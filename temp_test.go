// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on ioutil Go package, modified to add
// file extension option to File() and renamed to temp
// Lefteris Zafiris 2017

package temp

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestFile(t *testing.T) {
	f, err := File("/_not_exists_", "foo", "bar")
	if f != nil || err == nil {
		t.Errorf("File(`/_not_exists_`, `foo`, `bar`) = %v, %v", f, err)
	}

	dir := os.TempDir()
	f, err = File(dir, "temp_test", "tmp")
	if f == nil || err != nil {
		t.Errorf("File(dir, `temp_test`, `tmp`) = %v, %v", f, err)
	}
	if f != nil {
		f.Close()
		os.Remove(f.Name())
		re := regexp.MustCompile("^" + regexp.QuoteMeta(filepath.Join(dir, "temp_test")) + "[0-9]+\\.tmp$")
		if !re.MatchString(f.Name()) {
			t.Errorf("File(`"+dir+"`, `temp_test`, `tmp`) created bad name %s", f.Name())
		}
	}
	f, err = File(dir, "temp_test", "")
	if f == nil || err != nil {
		t.Errorf("File(dir, `temp_test`, ``) = %v, %v", f, err)
	}
	if f != nil {
		f.Close()
		os.Remove(f.Name())
		re := regexp.MustCompile("^" + regexp.QuoteMeta(filepath.Join(dir, "temp_test")) + "[0-9]+$")
		if !re.MatchString(f.Name()) {
			t.Errorf("File(`"+dir+"`, `temp_test`, ``) created bad name %s", f.Name())
		}
	}
	f, err = File(dir, "", "tmp")
	if f == nil || err != nil {
		t.Errorf("File(dir, ``, `tmp`) = %v, %v", f, err)
	}
	if f != nil {
		f.Close()
		os.Remove(f.Name())
		re := regexp.MustCompile("^" + regexp.QuoteMeta(dir) + "/[0-9]+\\.tmp$")
		if !re.MatchString(f.Name()) {
			t.Errorf("File(`"+dir+"`, ``, `tmp`) created bad name %s", f.Name())
		}
	}
}

func TestDir(t *testing.T) {
	name, err := Dir("/_not_exists_", "foo")
	if name != "" || err == nil {
		t.Errorf("Dir(`/_not_exists_`, `foo`) = %v, %v", name, err)
	}

	dir := os.TempDir()
	name, err = Dir(dir, "temp_test")
	if name == "" || err != nil {
		t.Errorf("Dir(dir, `temp_test`) = %v, %v", name, err)
	}
	if name != "" {
		os.Remove(name)
		re := regexp.MustCompile("^" + regexp.QuoteMeta(filepath.Join(dir, "temp_test")) + "[0-9]+$")
		if !re.MatchString(name) {
			t.Errorf("Dir(`"+dir+"`, `temp_test`) created bad name %s", name)
		}
	}
}

// test that we return a nice error message if the dir argument to TempDir doesn't
// exist (or that it's empty and os.TempDir doesn't exist)
func TestDir_BadDir(t *testing.T) {
	dir, err := Dir("", "TestTempDir_BadDir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	badDir := filepath.Join(dir, "not-exist")
	_, err = Dir(badDir, "foo")
	if pe, ok := err.(*os.PathError); !ok || !os.IsNotExist(err) || pe.Path != badDir {
		t.Errorf("Dir error = %#v; want PathError for path %q satisifying os.IsNotExist", err, badDir)
	}
}
