package services

import (
	"context"
	"encoding/json"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
	"github.com/quangnc2k/do-an-gang/pkg/something"
	"log"
	"math"
)

var ErrBlacklistPrefix = "blacklist filter failed: "

func ProcessGeneral(ctx context.Context, data string) (marked bool, threat model.Threat) {
	var connLog model.ConnLog

	var credit float64
	var xtra interface{}

	err := json.Unmarshal([]byte(data), &connLog)
	if err != nil {
		log.Println(ErrFilePrefix, err)
		return
	}

	if connLog.LocalOrig {
		marked, credit, xtra, err = persistance.BlacklistEngine.Check(ctx, connLog.ID.ResponseHost)
		if err != nil {
			log.Println(ErrBlacklistPrefix, err)
			return
		}
	} else {
		marked, credit, xtra, err = persistance.BlacklistEngine.Check(ctx, connLog.ID.OriginalHost)
		if err != nil {
			log.Println(ErrBlacklistPrefix, err)
			return
		}
	}

	threat = model.Threat{
		AffectedHost: connLog.ID.OriginalHost,
		AttackerHost: connLog.ID.ResponseHost,
		ConnID:       connLog.UID,
		Phase:        model.CnC,
	}

	if credit > 0 {
		threat.Severity += int(math.Floor(credit * 0.25))
		threat.Confidence += credit
	}

	threat.Metadata = something.CombineAsMetadata(connLog.Metadata, xtra)

	return
}
