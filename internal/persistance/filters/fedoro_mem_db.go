package filters

import (
	"bufio"
	"context"
	"github.com/quangnc2k/do-an-gang/pkg/hxxp"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type FeodoTrackerEngine struct {
	Blacklist map[string]bool
	mu        *sync.RWMutex
}

func (e *FeodoTrackerEngine) StoreIntoMem() (err error) {
	client := hxxp.NewHTTPClient()

	url := "https://feodotracker.abuse.ch/downloads/ipblocklist_recommended.txt"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	e.mu.Lock()
	defer e.mu.Unlock()
	for k := range e.Blacklist {
		e.Blacklist[k] = false
	}

	for scanner.Scan() {
		if scanner.Text()[0] == '#' {
			continue
		}

		if net.ParseIP(scanner.Text()) == nil {
			continue
		}

		e.Blacklist[scanner.Text()] = true
	}
	return nil
}

func (e *FeodoTrackerEngine) ClearFromMem() (err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for k, v := range e.Blacklist {
		if !v {
			delete(e.Blacklist, k)
		}
	}
	return nil
}

func (e *FeodoTrackerEngine) Check(ctx context.Context, resource string) (marked bool, credit float64, extraResource interface{}, err error) {
	extra := make(map[string]interface{})

	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.Blacklist[resource] {
		marked = true
		credit = 5
		extra["description"] = "Detected by Feodo Tracker block list"
	}

	log.Println("Checked with Feodo:", resource, "found:", marked)

	extraResource = extra
	return
}

func InitFedoroEngine(ctx context.Context) *FeodoTrackerEngine {
	mu := sync.RWMutex{}
	var m = make(map[string]bool)
	blacklistEngine := FeodoTrackerEngine{
		Blacklist: m,
		mu:        &mu,
	}

	err := blacklistEngine.StoreIntoMem()
	if err != nil {
		log.Fatalln("Initiate Blacklist Engine", err)
	}

	ticker := time.NewTicker(1 * time.Hour)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err = blacklistEngine.StoreIntoMem()
				if err != nil {
					log.Println("Update Blacklist Engine", err)
				}

				_ = blacklistEngine.ClearFromMem()
			}
		}
	}(ctx)

	return &blacklistEngine
}
