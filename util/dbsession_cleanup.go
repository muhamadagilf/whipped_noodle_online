// Package util
package util

import (
	"context"
	"log"
	"time"

	"github.com/muhamadagilf/whipped_noodle_online/internal/database"
)

func DBSessionCleanUp(q *database.Queries) {
	log.Println("[BG_WORKER_DEBUG]# DB sessions cleanup running")
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		<-ticker.C
		log.Println("[BG_WORKER_DEBUG]# cleanup schedule")
		err := q.DeleteExpiredSession(context.Background())
		if err != nil {
			log.Println("[BG_WORKER_ERROR]# ", err)
		}
	}
}
