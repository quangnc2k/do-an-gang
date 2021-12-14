package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
	"github.com/quangnc2k/do-an-gang/pkg/something"
	"log"
	"strings"
)

var ErrNoticePrefix = "notice filter failed: "

func ProcessNotice(ctx context.Context, data string) (marked bool, threat model.Threat) {
	var noticeLog model.NoticeLog

	err := json.Unmarshal([]byte(data), &noticeLog)
	if err != nil {
		log.Println(ErrNoticePrefix, err)
		return
	}

	if strings.HasPrefix(noticeLog.Note, "CaptureLoss") {
		return
	}

	err = noticeLog.SetMetadata()
	if err != nil {
		log.Println(ErrNoticePrefix, err)
		return
	}

	noticeLog.SetExtraResource()
	if noticeLog.ExtraResource == nil {
		return
	}

	src := noticeLog.Source
	if src == "" && noticeLog.ID != nil {
		src = noticeLog.ID.OriginalHost
	}

	dest := noticeLog.Destination
	if dest == "" && noticeLog.ID != nil {
		dest = noticeLog.ID.ResponseHost
	}

	fmt.Println("asdasdasdasdasdasdds", src)

	marked, credit, xtra, err := persistance.IPEngine.Check(ctx, src)
	if err != nil {
		log.Println(ErrNoticePrefix, err)
		return
	}

	m := something.CombineAsMetadata(noticeLog.Metadata, xtra, noticeLog.ExtraResource)

	threat = model.Threat{
		SeenAt:       something.ToTime(noticeLog.TS),
		AffectedHost: dest,
		AttackerHost: src,
		Confidence:   something.ExtractFromJsonMap(m, "confidence").(float64),
		Severity:     something.ExtractFromJsonMap(m, "severity").(float64) / 2,
		Phase:        something.ExtractFromJsonMap(m, "phase").(string),
	}

	if marked {
		threat.Severity += credit / 2
		threat.Confidence = 1
	}

	threat.Metadata = m
	return
}
