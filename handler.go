package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	// TODO: implement Handler
}

func NewHandler() *Handler {
	// TODO: implement Handler
	return &Handler{}
}

// HandleKurwa ...
func (h *Handler) HandleKurwa(ctx *gin.Context) {
	req := Request{}
	if err := ctx.BindJSON(req); err != nil {
		return
	}
	ctx.JSONP(http.StatusOK, req)
}

// HandlePshiek ...
func (h *Handler) HandlePshiek(ctx *gin.Context) {
	// TODO: implement stub
}

// HandlePierdole ...
func (h *Handler) HandlePierdole(ctx *gin.Context) {
	// TODO: implement stub
}

func Serve() error {
	h := NewHandler()
	g := gin.New()

	g.GET("/kurwa", h.HandleKurwa)
	g.POST("/pshiek", h.HandlePshiek)
	g.PUT("/pierdole", h.HandlePierdole)

	return g.Run(":8080")
}

func main() {
	Serve()
}

type Request struct {
	Entity string
}
