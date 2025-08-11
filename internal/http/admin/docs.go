package adminhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/tanna.dev/openapi-doc-http-handler/elements"
)

var docsHandler http.HandlerFunc

func (s *ServerImpl) GetV1AdminDocs(c *gin.Context) {
	docsHandler.ServeHTTP(c.Writer, c.Request)
}

func init() {
	swagger, err := GetSwagger()
	docsHandler, err = elements.NewHandler(swagger, err)
	if err != nil {
		panic(err)
	}
}
