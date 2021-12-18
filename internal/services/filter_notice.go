package services

import (
	"context"
	"encoding/json"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
	"github.com/quangnc2k/do-an-gang/pkg/something"
	"log"
	"math"
	"net"
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
	} else {
		marked = true
	}

	marked2, credit, xtra, aff, sus, err := doubleCheckNoticeLog(ctx, noticeLog)
	if err != nil {
		log.Println(ErrNoticePrefix, err)
		return
	}

	m := something.CombineAsMetadata(noticeLog.Metadata, xtra, noticeLog.ExtraResource)

	threat = model.Threat{
		SeenAt:       something.ToTime(noticeLog.TS),
		AffectedHost: aff,
		AttackerHost: sus,
		Severity:     something.ExtractFromJsonMap(m, "severity").(float64),
		Phase:        something.ExtractFromJsonMap(m, "phase").(string),
	}

	if marked2 {
		threat.Severity += math.Floor(credit / 3)
	}

	threat.Metadata = m
	return
}

func doubleCheckNoticeLog(ctx context.Context, l model.NoticeLog) (marked bool, credit float64, xtra interface{},
	aff, sus string, err error) {
	if l.ExtraResource.IsHostVictim {
		aff = l.Source
		// Check if notice has destination logged and whether it is dirty
		if l.Destination != "" {
			marked, credit, xtra, err = persistance.IPEngine.Check(ctx, l.Destination)
			if err != nil {
				log.Println(ErrNoticePrefix, err)
				return
			}

			if marked {
				return marked, credit, xtra, l.Source, l.Destination, nil
			}
		}

		// Check dirt in list of suspected addresses
		for _, addr := range l.SuspectedAddr {
			marked, credit, xtra, err = persistance.IPEngine.Check(ctx, addr)
			if err != nil {
				log.Println(ErrNoticePrefix, err)
				return
			}

			if marked {
				return marked, credit, xtra, l.Source, addr, nil
			}
		}

		// Check dirt in list of suspected hostname
		for _, host := range l.SuspectedHosts {
			addrs, err1 := net.LookupIP(host)
			if err1 != nil {
				log.Println(ErrNoticePrefix, err1)
				return
			}

			for _, addr := range addrs {
				marked, credit, xtra, err = persistance.IPEngine.Check(ctx, addr.String())
				if err != nil {
					log.Println(ErrNoticePrefix, err)
					return
				}

				if marked {
					return marked, credit, xtra, l.Source, addr.String(), nil
				}
			}
		}
	} else {
		sus = l.Source
		marked, credit, xtra, err = persistance.IPEngine.Check(ctx, l.Source)
		if err != nil {
			log.Println(ErrNoticePrefix, err)
			return
		}

		if marked {
			return marked, credit, xtra, "", l.Source, nil
		}
	}
	return
}
