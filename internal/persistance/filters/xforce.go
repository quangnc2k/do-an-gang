package filters

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/quangnc2k/do-an-gang/internal/config"
	"github.com/quangnc2k/do-an-gang/pkg/hxxp"
	"log"
	"net/http"
	"sync"
)

type XForceEngine struct {
	username string
	password string
	mu       *sync.Mutex
}

type XForceResponse struct {
	IP                   string            `json:"ip"`
	Cats                 map[string]int    `json:"cats"`
	Geo                  map[string]string `json:"geo"`
	CategoryDescriptions map[string]string `json:"categoryDescriptions"`
	Score                float32           `json:"score"`
}

func (e *XForceEngine) Check(ctx context.Context, resource string) (marked bool, credit float64, extraResource interface{}, err error) {
	client := hxxp.NewHTTPClient()

	url := fmt.Sprintf("https://exchange.xforce.ibmcloud.com/api/ipr/%s", resource)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(e.username, e.password)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	var respData XForceResponse

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return
	}

	if respData.Score > 5 {
		marked = true
		credit = float64(respData.Score)
	}

	log.Println("Checked with XForce:", resource, "got score:", respData.Score)

	return marked, credit, respData, nil
}

func InitXForceEngine(ctx context.Context) *XForceEngine {
	mu := new(sync.Mutex)

	return &XForceEngine{
		username: config.Env.XForceUsername,
		password: config.Env.XForcePassword,
		mu:       mu,
	}
}
