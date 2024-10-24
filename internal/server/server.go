package server

import (
	"sync"

	"github.com/derickit/go-rest-api/internal/db"
	"github.com/derickit/go-rest-api/internal/logger"
	"github.com/derickit/go-rest-api/internal/models"
)

var startOnce sync.Once

func StartOnce(svcEnv models.ServiceEnv, dbMgr db.MongoManager, lgr *logger.AppLogger) {
	startOnce.Do(func() {
		StartServer(svcEnv, dbMgr, lgr)
	})
}
