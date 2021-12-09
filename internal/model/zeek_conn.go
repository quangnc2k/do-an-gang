package model

import "encoding/json"

type ConnLog struct {
	TS                 float64  `json:"ts,omitempty"`
	UID                string   `json:"uid,omitempty"`
	ID                 *ConnID  `json:"id,omitempty"`
	Proto              string   `json:"proto,omitempty"`
	Service            string   `json:"service,omitempty"`
	Duration           string   `json:"duration,omitempty"`
	OrigBytes          int      `json:"orig_bytes,omitempty"`
	RespBytes          int      `json:"resp_bytes,omitempty"`
	ConnState          string   `json:"conn_state,omitempty"`
	LocalOrig          bool     `json:"local_orig,omitempty"`
	LocalResp          bool     `json:"local_resp,omitempty"`
	MissedBytes        int      `json:"missed_bytes,omitempty"`
	History            string   `json:"history,omitempty"`
	OrigPkts           int      `json:"orig_pkts,omitempty"`
	OrigIpBytes        int      `json:"orig_ip_bytes,omitempty"`
	RespPkts           int      `json:"resp_pkts,omitempty"`
	RespIpBytes        int      `json:"resp_ip_bytes,omitempty"`
	TunnelParents      []string `json:"tunnel_parents,omitempty"`
	OrigL2Addr         string   `json:"orig_l2_addr,omitempty"`
	RespL2Addr         string   `json:"resp_l2_addr,omitempty"`
	SpeculativeService string   `json:"speculative_service,omitempty"`

	Metadata map[string]interface{} `json:"-"`
}

func (log *ConnLog) SetMetadata() error {
	var m map[string]interface{}
	jsonized, err := json.Marshal(log)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonized, &m)
	if err != nil {
		return err
	}

	log.Metadata = m

	return nil
}
