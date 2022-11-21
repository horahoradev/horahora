package sockets

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	socketio "github.com/googollee/go-socket.io"
)

func New() *socketio.Server {
	s := socketio.NewServer(nil)
	s.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Errorf("connected:", s.ID())
		s.Join("bcast")
		return nil
	})

	return s
}

func Run(s *socketio.Server) {
	go s.Serve()
	defer s.Close()

	http.Handle("/socket.io/", s)

	log.Println("Serving at 0.0.0.0:8000...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}
