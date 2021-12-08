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

var ErrFilePrefix = "file filter failed: "

func ProcessFile(ctx context.Context, data string) (marked bool, threat model.Threat) {
	var fileLog model.FileLog

	err := json.Unmarshal([]byte(data), &fileLog)
	if err != nil {
		log.Println(ErrFilePrefix, err)
		return
	}

	marked, credit, xtra, err := persistance.FileEngine.Check(ctx, fileLog.MD5)
	if err != nil {
		log.Println(ErrFilePrefix, err)
		return
	}

	src := fileLog.TXHosts[0]

	dest := fileLog.RXHosts[0]

	threat = model.Threat{
		AffectedHost: src,
		AttackerHost: dest,
		Phase:        model.LM,
	}

	m := something.CombineAsMetadata(fileLog.Metadata, xtra)

	if credit > 0 {
		threat.Severity += int(math.Floor(credit * 0.25))
		threat.Confidence += credit
	}
	threat.ConnID = something.ExtractFromJsonMap(m, "fuid").(string)
	threat.Metadata = m

	return
}
