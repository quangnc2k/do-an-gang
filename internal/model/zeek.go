package model

type ConnID struct {
	OriginalHost string `json:"orig_h,omitempty"`
	OriginalPort Port   `json:"orig_p,omitempty"`
	ResponseHost string `json:"resp_h,omitempty"`
	ResponsePort Port   `json:"resp_p,omitempty"`
}

type HTTPEntity struct {
	Filename string `json:"filename,omitempty"`
}

type Port struct {
	Port  int64  `json:"port,omitempty"`
	Proto string `json:"proto,omitempty"`
}
