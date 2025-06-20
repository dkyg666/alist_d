package bootstrap

import (
	log "github.com/sirupsen/logrus"
	"github.com/vscodev/alist/v3/internal/search"
)

func InitIndex() {
	progress, err := search.Progress()
	if err != nil {
		log.Errorf("init index error: %+v", err)
		return
	}
	if !progress.IsDone {
		progress.IsDone = true
		search.WriteProgress(progress)
	}
}
