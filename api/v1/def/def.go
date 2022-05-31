package def

const (
	ROUTE_V1                        = "/api/v1"                // prefix
	HEADER_TOKEN_KEY                = "Token"                  // token's key name in request header
	HEADER_SecWebSocketProtocol_KEY = "Sec-WebSocket-Protocol" // token's key name in websocket req header
)

type BaseResp struct {
	Error string `json:"error"`
}

const (
	StatusSuccess = iota
	StatusFail
)
