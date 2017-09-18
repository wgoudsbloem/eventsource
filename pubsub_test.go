package main

import (
	"reflect"
	"testing"
)

func TestSubscribe(t *testing.T) {
	mp := make(topic)
	ps := pubsub{mp}
	testNotifyFunc := func() {}
	id := ps.subscribe(testNotifyFunc)
	if string(id) == "" {
		t.Errorf("id can not be empty %+v", id)
	}
	if mp[id] == nil {
		t.Error("map topic can not be empty")
	}
	fn, ok := mp[id]
	if ok {
		if reflect.ValueOf(fn).Pointer() != reflect.ValueOf(testNotifyFunc).Pointer() {
			t.Error("the given functions are not equal")
		}
	}
}

func TestUnSubscribe(t *testing.T) {
	mp := make(topic)
	ps := pubsub{mp}
	testNotifyFunc := func() {}
	id := ps.subscribe(testNotifyFunc)
	ps.unsubscribe(id)
	_, ok := mp[id]
	if ok {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", false, ok)
	}
}

func TestNotify(t *testing.T) {
	ps := newPubsub()
	want := "testval"
	var got string
	tnf := func() {
		got = want
	}
	ps.subscribe(tnf)
	ps.notify()
	if got != want {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", want, got)
	}
}
