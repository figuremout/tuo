package charts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1_def "github.com/githubzjm/tuo/api/v1/def"
	"github.com/githubzjm/tuo/routers"
)

func noRouteHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, v1_def.BaseResp{
		Error: "404, page not exists!",
	})
}

func Routers(e *gin.Engine) {
	e.NoRoute(noRouteHandler)

	e.LoadHTMLGlob("web/charts/**/*")
	// e.GET(fmt.Sprintf(def.ROUTE_CPU+def.SUBROUTE_PERCENT, ":"+clusterIDParam, ":"+nodeIDParam), func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "percent.html", gin.H{
	// 		"title": "CPU Percent",
	// 		"url":   "ws://localhost:8080/api/v1/clusters/1/nodes/4/cpu/percent", // TODO should not hard code
	// 	})
	// })

}

func init() {
	routers.Include(Routers)
}
