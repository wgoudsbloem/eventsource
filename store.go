package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

const filename = "db.log"

//var ps = newPubsub()

type logDb struct {
	io.Writer
	*pubsub
}

func (ldb logDb) Save(i int) (err error) {
	if _, err = ldb.Write([]byte(fmt.Sprintf("%d\n", i))); err != nil {
		return
	}
	ldb.notify()
	return
}

func (ldb logDb) GetPubSub() *pubsub {
	return ldb.pubsub
}

// NewDb is the default factory method to get a DB
func NewDb() Db {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	ps := newPubsub()
	return logDb{f, ps}
}

type viewLog struct {
	io.ReadSeeker
}

func (vl viewLog) GetSum(f func(r io.Reader) (int, error)) (int, error) {
	vl.Seek(0, os.SEEK_SET)
	return f(vl)
}

// NewView is the default factory method to get a View
func NewView() View {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return viewLog{f}
}

func totaling(in io.Reader) (o int, err error) {
	bio := bufio.NewReader(in)
	for true {
		var val string
		val, err = bio.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
				return
			}
		}
		var i int
		i, err = strconv.Atoi(val[:len(val)-1])
		if err != nil {
			return
		}
		o += i
	}
	return
}
