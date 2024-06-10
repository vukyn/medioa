package server

import (
	initStorage "medioa/internal/storage/init"

	"github.com/gin-gonic/gin"
)

func (s *Server) initHandler(group *gin.RouterGroup) {
	// Init storage

	storage := initStorage.NewInit(s.lib, s.cfg)
	storage.Handler.MapRoutes(group)

}
