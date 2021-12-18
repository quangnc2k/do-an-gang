package services

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
	"log"
	"sync"
	"time"
)

type RawPayload struct {
	Channel string
	Data    string
}

const popper = 1

func PopDataFromQueue(ctx context.Context, producers *sync.WaitGroup, outputChan chan<- RawPayload) {
	// setup workers to pop and send payloads to inputs channel
	for i := 0; i < popper; i++ {
		producers.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, outputChan chan<- RawPayload) {
			// this goroutine ends when the worker encounters an error from rdb
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					ch, data, err := persistance.GetQueue().Pop(ctx)
					if err != nil {
						if err != persistance.ErrOutOfItem {
							log.Println("worker encountered error", "error", err.Error())
							return
						}
						continue
					}
					if len(outputChan) > 500 {
						time.Sleep(5 * time.Second)
					}
					outputChan <- RawPayload{Channel: ch, Data: data}
				}
			}
		}(ctx, producers, outputChan)
	}
}
