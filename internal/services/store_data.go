package services

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
	"log"
	"sync"
	"time"
)

func StoreData(ctx context.Context, wg *sync.WaitGroup, threatChan chan model.Threat) {
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(2 * time.Second)
		var threatArr []model.Threat
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := persistance.GetRepoContainer().ThreatRepository.StoreThreatInBatch(ctx, threatArr)
				if err != nil {
					log.Println("store threat failed:", err)
				}

				threatArr = []model.Threat{}
			case threat := <-threatChan:
				threatArr = append(threatArr, threat)
			}

		}
	}()
}
