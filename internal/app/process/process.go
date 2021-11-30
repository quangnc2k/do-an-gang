package process

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"github.com/quangnc2k/do-an-gang/internal/services"
	"log"
	"sync"
)

func Run(ctx context.Context) (err error) {
	var producers, consumer, storer sync.WaitGroup

	var processChan chan services.RawPayload
	var threatChan chan model.Threat

	log.Println("Running....")
	services.PopDataFromQueue(ctx, &producers, processChan)
	services.ProcessData(ctx, &consumer, processChan, threatChan)
	services.StoreData(ctx, &storer, threatChan)

	producers.Wait()
	close(processChan)

	consumer.Wait()
	close(threatChan)

	storer.Wait()

	return nil
}
