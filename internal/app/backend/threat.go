package backend

import (
	"github.com/quangnc2k/do-an-gang/internal/persistance"
	"github.com/quangnc2k/do-an-gang/pkg/hxxp"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func threatsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pageString := r.URL.Query().Get("page")
	if pageString == "" {
		pageString = "1"
	}

	page, _ := strconv.Atoi(pageString)

	ppageString := r.URL.Query().Get("perPage")
	if pageString == "" {
		pageString = "1"
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
				hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
				return
			}
		}

		if len(orderBy) >= 2 {
			hxxp.RespondJson(w, http.StatusBadRequest, "invalid query", nil)
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