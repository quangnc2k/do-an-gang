package model

import "encoding/json"

type FileLog struct {
	TS               float64                `json:"ts,omitempty"`
	FUID             string                 `json:"fuid,omitempty"`
	TXHosts          []string               `json:"tx_hosts,omitempty"`
	RXHosts          []string               `json:"rx_hosts,omitempty"`
	ConnUIDs         []string               `json:"conn_uids,omitempty"`
	Source           string                 `json:"source,omitempty"`
	Depth            int                    `json:"depth,omitempty"`
	Analyzers        []string               `json:"analyzers,omitempty"`
	MIMEType         string                 `json:"mime_type,omitempty"`
	Filename         string                 `json:"filename,omitempty"`
	Duration         string                 `json:"duration,omitempty"`
	LocalOrig        bool                   `json:"local_orig,omitempty"`
	IsOrig           bool                   `json:"is_orig,omitempty"`
	SeenBytes        int                    `json:"seen_bytes,omitempty"`
	TotalBytes       int                    `json:"total_bytes,omitempty"`
	MissingBytes     int                    `json:"missing_bytes,omitempty"`
	OverflowBytes    int                    `json:"overflow_bytes,omitempty"`
	Timedout         bool                   `json:"timedout,omitempty"`
	ParentFUID       string                 `json:"parent_fuid,omitempty"`
	MD5              string                 `json:"md5,omitempty"`
	SHA1             string                 `json:"sha1,omitempty"`
	X509             map[string]interface{} `json:"x509,omitempty"`
	Extracted        string                 `json:"extracted,omitempty"`
	ExtractedCutoff  bool                   `json:"extracted_cutoff,omitempty"`
	ExtractedSize    int                    `json:"extracted_size,omitempty"`
	Entropy          float64                `json:"entropy,omitempty"`
	Describe         string                 `json:"describe,omitempty"`
	OriginalMAC      string                 `json:"orig_mac,omitempty"`
	OriginalHostName string                 `json:"orig_host_name"`
	ResponseMAC      string                 `json:"resp_mac,omitempty"`

	Metadata map[string]interface{} `json:"-"`
}

func (log *FileLog) SetMetadata() error {
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
