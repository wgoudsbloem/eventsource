package main

import (
	"bytes"
	"strconv"
	"testing"
)

func TestSave(t *testing.T) {
	var bb bytes.Buffer
	tdb := logDb{&bb, &pubsub{}}
	want := 999
	tdb.Save(want)
	b := make([]byte, len(string(want)+"\n"))
	if _, err := bb.Read(b); err != nil {
		t.Error(err)
	}
	got, err := strconv.Atoi(string(b))
	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want, got)
	}
}

// Bytes buffer doesn't have a seeker
type bufferWithSeeker struct {
	bytes.Buffer
}

func (bws *bufferWithSeeker) Seek(i int64, j int) (k int64, err error) {
	return
}

func TestRead(t *testing.T) {
	var bb bufferWithSeeker
	want := 111
	bb.WriteString("1\n")
	bb.WriteString("10\n")
	bb.WriteString("100\n")
	tvl := viewLog{&bb}
	got, err := tvl.GetSum(totaling)
	if err != nil {
		t.Error(err)
	}
	if got != 111 {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want, got)
	}
}

func TestGetPubSub(t *testing.T) {
	got := NewDb().GetPubSub().subscribe(func() {})
	if got == "" {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", "any id", got)
	}
}

func TestNewView(t *testing.T) {
	var got View
	got = NewView()
	if got == nil {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", "not nil", got)
	}
}
