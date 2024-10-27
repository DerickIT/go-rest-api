package handlers

import (
	"net/http"

	"github.com/derickit/go-rest-api/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ServiceStatus string

const (
	UP   ServiceStatus = "ok"
	DOWN ServiceStatus = "down"
)

type StatusResponse struct {
	Status      ServiceStatus
	ServiceName string
	UpTime      string
	Environment string
	Version     string
}

type StatusController struct {
	dbMgr db.MongoManager
}

func NewStatusController(dbMgr db.MongoManager) *StatusController {
	return &StatusController{dbMgr: dbMgr}
}

func (s *StatusController) CheckStatus(c *gin.Context) {
	var stat ServiceStatus
	var code int
	if err := s.dbMgr.Ping(); err != nil {
		stat = UP
		code = http.StatusOK
	} else {
		log.Error().Msg("unable to connect to DB")
		stat = DOWN
		code = http.StatusFailedDependency
	}
	c.JSON(code, stat)
}
