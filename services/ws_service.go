package services

import (
	"chatRoom/utilities"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/kataras/iris/v12/context"
	"log"
	"net/http"
)

type IWebsocketService interface {
	IsStart() bool
	Start(write context.ResponseWriter, request *http.Request) (*websocket.Conn, error)
	StartReceive(callBack func(IWebsocketService, []byte))
	WriteMessage(message string) error
	Close()
}

type WebsocketService struct {
	Configure websocket.Upgrader
	Conn      *websocket.Conn
}

func NewWebsocketService(cfg websocket.Upgrader) *WebsocketService {
	return &WebsocketService{
		Configure: cfg,
	}
}

func (w *WebsocketService) IsStart() bool {
	return w.Conn != nil
}

func (w *WebsocketService) Start(write context.ResponseWriter, request *http.Request) (*websocket.Conn, error) {
	ws, err := w.Configure.Upgrade(write, request, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return nil, err
	}
	w.Conn = ws

	return ws, nil
}

func (w *WebsocketService) StartReceive(callback func(IWebsocketService, []byte)) {
	defer w.Close()
	for {
		_, body, err := w.Conn.ReadMessage()
		utilities.FailOnError(err, "Fatal on websocket received")

		if err != nil {
			break
		}

		go callback(w, body)
	}
}

func (w *WebsocketService) WriteMessage(message string) error {
	if !w.IsStart() {
		return errors.New("Should Start server")
	}

	w.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	return nil
}

func (w *WebsocketService) Close() {
	w.Conn.Close()
}
