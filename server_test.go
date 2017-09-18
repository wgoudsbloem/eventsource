package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootHandler(t *testing.T) {
	s := server{}
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/html/index.html", nil)
	if err != nil {
		t.Error(err)
	}
	s.rootHandler(rec, req)
	b, err := ioutil.ReadFile("html/index.html")
	if err != nil {
		t.Error(err)
	}
	want := string(b)
	got, err := rec.Body.ReadString('|')
	if err != nil && err != io.EOF {
		t.Error(err)
	}
	if want != got {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want, got)
	}
}

func TestApiHandler(t *testing.T) {
	tdb := logDb{&bytes.Buffer{}, &pubsub{}}
	s := server{tdb, nil}
	rec := httptest.NewRecorder()
	body := struct {
		Amt int `josn:"amt"`
	}{}
	body.Amt = 555
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(body)
	req, err := http.NewRequest("POST", "/api", &b)
	if err != nil {
		t.Error(err)
	}
	s.apiHandler(rec, req)
	want := http.StatusCreated
	got := rec.Result().StatusCode
	if err != nil {
		t.Error(err)
	}
	if want != got {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want, got)
	}
}

type testReadWriteCloser struct {
	bytes.Buffer
}

func (*testReadWriteCloser) Close() error {
	return nil
}
