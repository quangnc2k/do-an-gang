package model

import (
	"encoding/json"
	"strings"
)

type NoticeLog struct {
	TS                float64     `json:"ts,omitempty"`
	UID               string      `json:"uid,omitempty"`
	ID                *ConnID     `json:"id,omitempty"`
	FUID              string      `json:"fuid,omitempty"`
	FileMimeType      string      `json:"file_mime_type,omitempty"`
	FileDesc          string      `json:"file_desc,omitempty"`
	Proto             string      `json:"proto,omitempty"`
	Note              string      `json:"note,omitempty"`
	Message           string      `json:"msg,omitempty"`
	Sub               string      `json:"sub,omitempty"`
	Source            string      `json:"src,omitempty"`
	Destination       string      `json:"dst,omitempty"`
	Port              Port        `json:"p,omitempty"`
	N                 int         `json:"n,omitempty"`
	PeerName          string      `json:"peer_name,omitempty"`
	PeerDescr         string      `json:"peer_descr,omitempty"`
	Actions           []string    `json:"actions,omitempty"`
	EmailBodySections []string    `json:"email_body_sections,omitempty"`
	EmailDelayTokens  []string    `json:"email_delay_tokens,omitempty"`
	Identifier        string      `json:"identifier,omitempty"`
	SuppressFor       string      `json:"suppress_for,omitempty"`
	RemoteLocation    interface{} `json:"remote_location,omitempty"`
	Dropped           bool        `json:"dropped,omitempty"`
	OriginalMAC       string      `json:"orig_mac"`
	OriginalHostName  string      `json:"orig_host_name"`

	ExtraResource *NoticeType            `json:"-"`
	Metadata      map[string]interface{} `json:"-"`
}

func (log *NoticeLog) SetMetadata() error {
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

func (log *NoticeLog) SetExtraResource() {
	for _, v := range noticeNoteMap {
		if strings.HasPrefix(log.Note, v.TypeHint) && (strings.Contains(log.Sub, v.Sub) || log.Sub == "") {
			log.ExtraResource = &v
			return
		}
	}
}