package base

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/githubzjm/tuo/api/v1/def"
	"github.com/githubzjm/tuo/internal/pkg/jwt"
	"github.com/githubzjm/tuo/internal/pkg/ws"
	"github.com/githubzjm/tuo/routers"
)

func noRouteHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, def.BaseResp{
		Error: "404, page not exists!",
	})
}

// func pingHandler(c *gin.Context) {
// 	c.JSON(200, common.Resp{
// 		Status: common.StatusSuccess,
// 		Msg:    "pong",
// 	})
// }

// func TestAuthHandler(c *gin.Context) {
// 	// claim is put into gin.Context in JWTAuth()
// 	claims := c.MustGet("claims").(*jwt.CustomClaims)

// 	c.JSON(http.StatusOK, common.Resp{
// 		Status: common.StatusSuccess,
// 		Msg:    "token is valid",
// 		Data:   claims,
// 	})

// }

func Routers(e *gin.Engine) {
	e.NoRoute(noRouteHandler)

	// e.LoadHTMLGlob("web/**/*") // related to the bin's work path
	e.LoadHTMLFiles("web/templates/index.tmpl", "web/test/ws.html")
	e.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})
	e.GET("/test/ws", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test/ws.html", gin.H{
			"title": "Test WebSocket",
		})
	})

	test := e.Group("/ws")
	test.Use(jwt.WSJWTAuth())
	test.GET("", ws.WsHandler)

	// test APIs

	// v1 := e.Group("/api/v1")
	// {
	// 	v1.GET("/ping", pingHandler)

	// 	test := v1.Group("/api/v1/test")
	// 	test.Use(jwt.JWTAuth())
	// 	{
	// 		test.GET("/auth", TestAuthHandler)
	// 	}
	// }

}

func init() {
	routers.Include(Routers)
}
