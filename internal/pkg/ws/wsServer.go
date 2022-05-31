package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/githubzjm/tuo/api/v1/def"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func WsHandler(c *gin.Context) {

	// use key Sec-WebSocket-Protocol of header to deliver token
	// if this key is used in websocket req, the resp of server must contains it
	token := c.Request.Header.Get(def.HEADER_SecWebSocketProtocol_KEY)

	var upgrader = websocket.Upgrader{
		// solve cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{token}, // https://www.jianshu.com/p/7b1deb1e0a07
	}

	// not sure is useful
	// //判断请求是否为websocket升级请求。
	// if websocket.IsWebSocketUpgrade(c.Request) {
	// 	conn, err := upgrader.Upgrade(w, r, w.Header())
	// } else {
	// 	//处理普通请求
	// 	c := newContext(w, r)
	// 	e.router.handle(c)
	// }

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer ws.Close()

	for {
		// Read message from client
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		// Print the message to the console
		log.Printf("recv: %s", msg)

		// Write message back to client
		rst := append([]byte("hello "), msg...)
		if err = ws.WriteMessage(msgType, rst); err != nil {
			log.Println("write:", err)
			return
		}
	}
}
