package common

import (
	"github.com/gin-gonic/gin"
	"github.com/viettranx/service-context/core"
	"net/http"
)

type CanWithDebug interface {
	WithDebug(debug string) *core.DefaultError
}

func WriteErrorResponse(c *gin.Context, err error) {
	if errSt, ok := err.(core.StatusCodeCarrier); ok {
		if !gin.IsDebugging() {
			errWithNoDebug := errSt.(CanWithDebug).WithDebug("")
			c.JSON(errSt.StatusCode(), errWithNoDebug)
			return
		}

		c.JSON(errSt.StatusCode(), errSt)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
