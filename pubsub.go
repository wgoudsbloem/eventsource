package main

import (
	"strconv"
	"time"
)

type id string
type notifyFunc func()
type topic map[id]notifyFunc

type pubsub struct {
	t topic
}

// subscribe the notification function and receive the reference id in case for future maintenance
func (ps *pubsub) subscribe(nf notifyFunc) (_id id) {
	_id = id(guuid())
	ps.t[_id] = nf
	return
}

// unsubsribe the id from the pubsub
func (ps *pubsub) unsubscribe(_id id) {
	delete(ps.t, _id)
}

// call when to notify all subscribers
func (ps *pubsub) notify() {
	for _, v := range ps.t {
		v()
	}
}

// default pubsub factory
func newPubsub() *pubsub {
	m := make(topic)
	return &pubsub{m}
}

// obviously very flawed guuid
func guuid() string {
	return strconv.Itoa(int(time.Now().Unix()))
}
