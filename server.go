package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

// Db is the interface for storage persistence
type Db interface {
	Save(i int) error
	GetPubSub() *pubsub
}

// View is the interface for storage view
type View interface {
	GetSum(f func(r io.Reader) (int, error)) (int, error)
}

func main() {
	db := NewDb()
	v := NewView()
	s := server{db, v}
	http.HandleFunc("/", s.rootHandler)
	http.HandleFunc("/api", s.apiHandler)
	http.Handle("/socket", websocket.Handler(s.socketHandler))
	http.ListenAndServe(":8888", nil)
}

type server struct {
	Db
	View
}

func (s *server) rootHandler(w http.ResponseWriter, r *http.Request) {
	f, err := os.OpenFile(r.URL.String()[1:], os.O_RDONLY, 0666)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	defer f.Close()
	io.Copy(w, f)
}

func (s *server) apiHandler(w http.ResponseWriter, r *http.Request) {
	var acc struct {
		Amt int `json:"amt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		fmt.Println(err)
	}
	if err := s.Save(acc.Amt); err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *server) socketHandler(ws *websocket.Conn) {
	bb, err := ioutil.ReadAll(ws)
	if err != nil {
		fmt.Println(err)
	}
	t := time.Now().Format("Mon Jan _2 15:04:05")
	fmt.Printf("websocket open\t%+v\t%+s\n", string(bb), t)
	f := func() (err error) {
		tot, err := s.GetSum(totaling)
		if err != nil {
			return
		}
		if err = websocket.JSON.Send(ws, tot); err != nil {
			return
		}
		return
	}
	if err := f(); err != nil {
		fmt.Println(err)
	}
	ps := s.GetPubSub()
	idd := ps.subscribe(func() {
		f()
	})
	var b bytes.Buffer
	if _, err := ws.Read(b.Bytes()); err == io.EOF {
		if err := ws.Close(); err != nil {
			fmt.Println(err)
		}
		ps.unsubscribe(idd)
		fmt.Println("websocket connection closed")
	}
}
