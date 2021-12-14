package services

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"sync"
)

const workerCap = 40

func ProcessData(ctx context.Context, consumers *sync.WaitGroup,
	inputChan chan RawPayload, outputChan chan<- model.Threat) {
	for i := 0; i < workerCap; i++ {
		// setup workers to handle payloads from inputs channel
		consumers.Add(1) // this is to notify the whole handling process is complete
		go func(ctx context.Context, wg *sync.WaitGroup, input chan RawPayload, outputChan chan<- model.Threat) {
			defer wg.Done()
			for payload := range inputChan {
				switch payload.Channel {
				case "file":
					marked, threat := ProcessFile(ctx, payload.Data)
					if marked {
						outputChan <- threat
					}

				case "notice":
					marked, threat := ProcessNotice(ctx, payload.Data)
					if marked {
						outputChan <- threat
					}

				default:
					marked, threat := ProcessGeneral(ctx, payload.Data)
					if marked {
						outputChan <- threat
					}
				}
			}
		}(ctx, consumers, inputChan, outputChan)
	}
}
