package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func WriteErrorResponse(c *gin.Context, err error) {
	//if errSt, ok := err.(core.StatusCodeCarrier); ok {
	//	c.JSON(errSt.StatusCode(), errSt)
	//	return
	//}

	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
