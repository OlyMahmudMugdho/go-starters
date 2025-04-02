package types

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

func (c *Context) Json(message interface{}) {
	c.Writer.Header().Add("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)

	err := json.NewEncoder(c.Writer).Encode(message)

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
