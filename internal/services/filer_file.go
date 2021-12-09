package services

import (
	"context"
	"encoding/json"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
	"github.com/quangnc2k/do-an-gang/pkg/something"
	"log"
)

var ErrFilePrefix = "file filter failed: "

func ProcessFile(ctx context.Context, data string) (marked bool, threat model.Threat) {
	var fileLog model.FileLog

	err := json.Unmarshal([]byte(data), &fileLog)
	if err != nil {
		log.Println(ErrFilePrefix, err)
		return
	}

	err = fileLog.SetMetadata()
	if err != nil {
		log.Println(ErrFilePrefix, err)
		return
	}

	marked, credit, xtra, err := persistance.FileEngine.Check(ctx, fileLog.MD5)
	if err != nil {
		log.Println(ErrFilePrefix, err)
		return
	}

	transmitter := fileLog.TXHosts[0]

	receiver := fileLog.RXHosts[0]

	threat = model.Threat{
		AffectedHost: receiver,
		AttackerHost: transmitter,
		Phase:        model.LM,
	}

	m := something.CombineAsMetadata(fileLog.Metadata, xtra)

	if credit > 0 {
		threat.Severity += credit * 0.25
		threat.Confidence += credit
	}
	threat.ConnID = something.ExtractFromJsonMap(m, "fuid").(string)
	threat.Metadata = m

	return
}
