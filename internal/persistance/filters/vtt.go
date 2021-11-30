package filters

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/quangnc2k/do-an-gang/internal/config"
	"github.com/quangnc2k/do-an-gang/pkg/hxxp"
	"net/http"
	"sync"
	"time"
)

type VTTEngine struct {
	apiKey string
	mu     *sync.Mutex
}

type VTTResponse struct {
	Data struct {
		LastAnalysisStats struct {
			Harmless   int `json:"harmless"`
			Suspicious int `json:"suspicious"`
			Malicious  int `json:"malicious"`
			Undetected int `json:"undetected"`
		} `json:"last_analysis_stats"`
	} `json:"data"`
}

func (e *VTTEngine) Check(ctx context.Context, resource string) (marked bool, credit float64, extraResource interface{}, err error) {
	client := hxxp.NewHTTPClient()

	url := fmt.Sprintf(" https://www.virustotal.com/api/v3/files/%s", resource)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	req.Header.Set("x-apikey", e.apiKey)

	e.mu.Lock()
	defer func() {
		time.Sleep(time.Duration(60 / config.Env.VTTMaxFilerPerMin) * 1000)
		e.mu.Unlock()
	}()

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	var respData VTTResponse

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return
	}

	if respData.Data.LastAnalysisStats.Malicious+respData.Data.LastAnalysisStats.Suspicious > 0 {
		marked = true
		credit = respData.getCredit()
	}

	return marked, credit, respData.Data.LastAnalysisStats, nil
}

func (r *VTTResponse) getCredit() float64 {
	return float64((r.Data.LastAnalysisStats.Malicious + r.Data.LastAnalysisStats.Suspicious) / (r.Data.LastAnalysisStats.Malicious + r.Data.LastAnalysisStats.Suspicious + r.Data.LastAnalysisStats.Harmless + r.Data.LastAnalysisStats.Undetected))
}

func InitVTTEngine(ctx context.Context) *VTTEngine{
	mu := new(sync.Mutex)

	return &VTTEngine{
		apiKey: config.Env.VTTApiKey,
		mu:     mu,
	}
}
