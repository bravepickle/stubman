package main

import (
	"bufio"
	//	"net/http"
	//	"net/http/httptest"
	"os"
	//	"regexp"
	"strings"
	"testing"
)

type helperDbStruct struct {
	db   *Db
	test *testing.T
}

func (h *helperDbStruct) resetDbSchema() {
	h.test.Log(`Reset DB schema:`, h.db.DbName)

	err := h.db.Reset()
	if err != nil {
		h.test.Fatal(err)
	}
}

func (h *helperDbStruct) loadFixtures(stubs []Stub) {
	for _, stub := range stubs {
		_, err := repo.Insert(&stub)
		if err != nil {
			h.test.Fatal(err)
		}
	}
}

func NewHelperDb(t *testing.T) *helperDbStruct {
	return &helperDbStruct{DefaultDb, t}
}

func assertFileContains(t *testing.T, txt string, path string) {
	fh, err := os.Open(path)
	if err != nil {
		t.Fatal(`Failed to open requests log file:`, err)
	}

	scanner := bufio.NewScanner(fh)
	found := false
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), txt) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf(`Failed to find in file "%s" string: %s`, path, txt)
	}
}
