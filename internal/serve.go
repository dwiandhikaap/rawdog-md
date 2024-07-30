package internal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// Serve built files
func Serve(watcherServer *WatcherServer, dir string, port int) {
	//mux := http.NewServeMux()

	http.Handle("/", http.FileServer(http.Dir(dir)))
	http.HandleFunc("/__rawdog-md/watcher", watcherHandler(watcherServer))

	portStr := fmt.Sprintf("127.0.0.1:%d", port)
	fmt.Println("Serving on http://" + portStr)
	http.ListenAndServe(portStr, nil)

}

type WatcherServer struct {
	clients []*websocket.Conn
}

func NewWatcherServer() *WatcherServer {
	return &WatcherServer{}
}

func (w *WatcherServer) AddClient(conn *websocket.Conn) {
	w.clients = append(w.clients, conn)
}

func (w *WatcherServer) RemoveClient(conn *websocket.Conn) {
	for i, c := range w.clients {
		if c == conn {
			w.clients = append(w.clients[:i], w.clients[i+1:]...)
			break
		}
	}
}

func (w *WatcherServer) Broadcast(message string) {
	msgBytes := []byte(message)
	for _, c := range w.clients {
		err := c.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func watcherHandler(watcherServer *WatcherServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return strings.HasPrefix(r.Host, "localhost")
			},
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		watcherServer.AddClient(conn)

		defer watcherServer.RemoveClient(conn)
		defer conn.Close()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}
}
