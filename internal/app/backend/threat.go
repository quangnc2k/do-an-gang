package backend

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/quangnc2k/do-an-gang/internal/persistance"
	"github.com/quangnc2k/do-an-gang/pkg/hxxp"
)

type bucket struct {
	asString string
	duration time.Duration
}

var buckets = []bucket{
	{
		asString: "15 min",
		duration: 15 * time.Minute,
	},
	{
		asString: "1 hour",
		duration: 1 * time.Hour,
	},
	{
		asString: "2 hour",
		duration: 2 * time.Hour,
	},
	{
		asString: "4 hour",
		duration: 4 * time.Hour,
	},
	{
		asString: "6 hour",
		duration: 6 * time.Hour,
	},
	{
		asString: "12 hour",
		duration: 12 * time.Hour,
	},
	{
		asString: "1 day",
		duration: 24 * time.Hour,
	},
	{
		asString: "3 day",
		duration: 3 * 24 * time.Hour,
	},
	{
		asString: "7 day",
		duration: 7 * 24 * time.Hour,
	},
	{
		asString: "14 day",
		duration: 14 * 24 * time.Hour,
	},
}

func threatsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pageString := r.URL.Query().Get("page")
	if pageString == "" {
		pageString = "1"
	}

	page, _ := strconv.Atoi(pageString)

	ppageString := r.URL.Query().Get("perPage")
	if ppageString == "" {
		ppageString = "1"
	}

	ppage, _ := strconv.Atoi(ppageString)

	orderByString := r.URL.Query().Get("orderBy")
	if orderByString != "" {
		orderBy := strings.Split(orderByString, " ")
		if orderBy[0] != "created_at" && orderBy[0] != "seen_at" && orderBy[0] != "confidence" && orderBy[0] != "severity" && orderBy[0] != "phase" {
			hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
			return
		}

		if len(orderBy) == 2 {
			if orderBy[1] != "asc" && orderBy[1] != "desc" {
				hxxp.RespondJson(w, http.StatusBadRequest, fmt.Sprintf("invalid query: %s", orderBy), nil)
				return
			}
		}

		if len(orderBy) > 2 {
			hxxp.RespondJson(w, http.StatusBadRequest, fmt.Sprintf("invalid query: %s", orderBy), nil)
			return
		}
	}

	search := r.URL.Query().Get("search")

	startString := r.URL.Query().Get("start")
	if len(startString) == 0 {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	from, err := time.Parse("2006-01-02T15:04:05.000Z", startString)
	if err != nil {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	endString := r.URL.Query().Get("end")
	if len(endString) == 0 {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	to, err := time.Parse("2006-01-02T15:04:05.000Z", endString)
	if err != nil {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	data, err := persistance.GetRepoContainer().ThreatRepository.Paginate(ctx, page, ppage, orderByString, search, from, to)
	if err != nil {
		hxxp.RespondJson(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	hxxp.RespondJson(w, 200, "Success", data)
}

func threatStatsPhase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	startString := r.URL.Query().Get("start")
	if len(startString) == 0 {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	from, err := time.Parse("2006-01-02T15:04:05.000Z", startString)
	if err != nil {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	endString := r.URL.Query().Get("end")
	if len(endString) == 0 {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	to, err := time.Parse("2006-01-02T15:04:05.000Z", endString)
	if err != nil {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	data, err := persistance.GetRepoContainer().ThreatRepository.StatsByPhase(ctx, from, to)
	if err != nil {
		hxxp.RespondJson(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	hxxp.RespondJson(w, 200, "Success", data)
}

func threatStatsSeverity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	startString := r.URL.Query().Get("start")
	if len(startString) == 0 {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	from, err := time.Parse("2006-01-02T15:04:05.000Z", startString)
	if err != nil {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	endString := r.URL.Query().Get("end")
	if len(endString) == 0 {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	to, err := time.Parse("2006-01-02T15:04:05.000Z", endString)
	if err != nil {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	data, err := persistance.GetRepoContainer().ThreatRepository.StatsBySeverity(ctx, from, to)
	if err != nil {
		hxxp.RespondJson(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	hxxp.RespondJson(w, 200, "Success", data)
}

func threatTopAffectedHost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	startString := r.URL.Query().Get("start")
	if len(startString) == 0 {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	from, err := time.Parse("2006-01-02T15:04:05.000Z", startString)
	if err != nil {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	endString := r.URL.Query().Get("end")
	if len(endString) == 0 {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	to, err := time.Parse("2006-01-02T15:04:05.000Z", endString)
	if err != nil {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	data, err := persistance.GetRepoContainer().ThreatRepository.TopHostAffected(ctx, from, to)
	if err != nil {
		hxxp.RespondJson(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	hxxp.RespondJson(w, 200, "Success", data)
}

func threatTopAttacker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	startString := r.URL.Query().Get("start")
	if len(startString) == 0 {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	from, err := time.Parse("2006-01-02T15:04:05.000Z", startString)
	if err != nil {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	endString := r.URL.Query().Get("end")
	if len(endString) == 0 {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	to, err := time.Parse("2006-01-02T15:04:05.000Z", endString)
	if err != nil {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	data, err := persistance.GetRepoContainer().ThreatRepository.TopAttacker(ctx, from, to)
	if err != nil {
		hxxp.RespondJson(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	hxxp.RespondJson(w, 200, "Success", data)
}

func threatHistogramAffected(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	startString := r.URL.Query().Get("start")
	if len(startString) == 0 {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	from, err := time.Parse("2006-01-02T15:04:05.000Z", startString)
	if err != nil {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	endString := r.URL.Query().Get("end")
	if len(endString) == 0 {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	to, err := time.Parse("2006-01-02T15:04:05.000Z", endString)
	if err != nil {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	if to.Before(from) {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
		return
	}

	diff := to.Sub(from)
	spaceString := ""
	space := diff / 10

	for _, b := range buckets {
		if b.duration < space {
			spaceString = b.asString
		}
	}

	if spaceString == "" {
		spaceString = "15 min"
	}

	data, err := persistance.GetRepoContainer().ThreatRepository.HistogramAffected(ctx, spaceString, from, to)
	if err != nil {
		hxxp.RespondJson(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	hxxp.RespondJson(w, 200, "Success", data)
}
