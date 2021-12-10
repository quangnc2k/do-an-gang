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
		Attributes struct {
			LastAnalysisStats struct {
				Harmless   int `json:"harmless"`
				Suspicious int `json:"suspicious"`
				Malicious  int `json:"malicious"`
				Undetected int `json:"undetected"`
			} `json:"last_analysis_stats"`
		} `json:"attributes"`
	} `json:"data"`
}

func (e *VTTEngine) Check(ctx context.Context, resource string) (marked bool, credit float64, extraResource interface{}, err error) {
	client := hxxp.NewHTTPClient()

	url := fmt.Sprintf("https://www.virustotal.com/api/v3/files/%s", resource)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	req.Header.Set("x-apikey", e.apiKey)

	e.mu.Lock()
	defer func() {
		time.Sleep(time.Duration(60/config.Env.VTTMaxFilerPerMin) * 1000000000)
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

	fmt.Println("Checked with vtt:", url, "got result", respData)

	if respData.Data.Attributes.LastAnalysisStats.Malicious+
		respData.Data.Attributes.LastAnalysisStats.Suspicious > 0 {
		marked = true
		credit = respData.getCredit()
	}

	return marked, credit, respData.Data.Attributes.LastAnalysisStats, nil
}

func (r *VTTResponse) getCredit() float64 {
	total := r.Data.Attributes.LastAnalysisStats.Malicious +
		r.Data.Attributes.LastAnalysisStats.Suspicious +
		r.Data.Attributes.LastAnalysisStats.Harmless +
		r.Data.Attributes.LastAnalysisStats.Undetected

	score := r.Data.Attributes.LastAnalysisStats.Malicious + r.Data.Attributes.LastAnalysisStats.Suspicious
	return float64(score) / float64(total)
}

func InitVTTEngine(ctx context.Context) *VTTEngine {
	mu := new(sync.Mutex)

	return &VTTEngine{
		apiKey: config.Env.VTTApiKey,
		mu:     mu,
	}
}
