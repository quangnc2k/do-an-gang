package model

const (
	GA  = "Gain Access"
	EF  = "Establish Foothold"
	DA  = "Deepen Access"
	LM  = "Lateral Movement"
	LLR = "Look, Learn, Remain"
	UK  = "Undefined"
)

type NoticeType struct {
	Resource   string  `json:"resource"`
	Origin     string  `json:"origin"`
	TypeHint   string  `json:"type_hint"`
	Sub        string  `json:"sub"`
	Severity   float64 `json:"severity"`
	Confidence float64 `json:"confidence"`
	Phase      string  `json:"phase"`

	IsHostVictim bool `json:"-"`
}

var noticeNoteMap = map[string]NoticeType{
	"Heartbleed": {
		Origin:   "SSL Heartbleed Detection",
		TypeHint: "Heartbleed::",
		Severity: 7.5,
		Phase:    GA,
	},

	"Beacon Detection": {
		Origin:   "Botnet Beacon Detection",
		TypeHint: "BOTNET::",
		Severity: 5.5,
		Phase:    EF,
	},

	"SQL Injection Attacker": {
		Origin:   "SQL Injection Detection",
		TypeHint: "HTTP::SQL_Injection_Attacker",
		Severity: 9.5,
		Phase:    GA,
	},

	"SQL Injection Victim": {
		Origin:       "SQL Injection Detection",
		TypeHint:     "HTTP::SQL_Injection_Victim",
		Severity:     9.5,
		Phase:        GA,
		IsHostVictim: true,
	},

	"SSH Brute": {
		Origin:   "SSH Bruteforce Detection",
		TypeHint: "SSH::Password_Guessing",
		Severity: 8.4,
		Phase:    DA,
	},

	"NCSA Scan": {
		Origin:   "Scan Detection",
		TypeHint: "Scan::",
		Severity: 4.2,
		Phase:    EF,
	},

	"BZAR Credential Access": {
		Origin:   "BZAR Module",
		TypeHint: "ATTACK::Credential_Access",
		Severity: 7.9,
		Phase:    DA,
	},

	"BZAR Defense Evasion": {
		Origin:   "BZAR Module",
		TypeHint: "ATTACK::Defense_Evasion",
		Severity: 4.8,
		Phase:    LLR,
	},

	"BZAR Discovery": {
		Origin:   "BZAR Module",
		TypeHint: "ATTACK::Discovery",
		Severity: 4.2,
		Phase:    DA,
	},

	"BZAR Execution": {
		Origin:   "BZAR Module",
		TypeHint: "ATTACK::Execution",
		Severity: 7.2,
		Phase:    EF,
	},

	"BZAR Lateral Movement": {
		Origin:   "BZAR Module",
		TypeHint: "ATTACK::Lateral_Movement",
		Severity: 7,
		Phase:    LM,
	},

	"BZAR Lateral Movement and Execution": {
		Origin:   "BZAR Module",
		TypeHint: "ATTACK::Lateral_Movement_and_Execution",
		Severity: 8,
		Phase:    LM,
	},

	"BZAR Lateral Movement Extracted File": {
		Origin:   "BZAR Module",
		TypeHint: "ATTACK::Lateral_Movement_Extracted_File",
		Severity: 8,
		Phase:    LM,
	},

	"BZAR Lateral Movement Multiple Attempts": {
		Origin:   "BZAR Module",
		TypeHint: "ATTACK::Lateral_Movement_Multiple_Attempts",
		Severity: 7,
		Phase:    LM,
	},
	//"BZAR::Lateral_Movement_T1021.002": {
	//	Origin:     "BZAR MITRE Lateral Movement",
	//	TypeHint:   "ATTACK::Lateral_Movement",
	//	Sub:        "T1021.002 Remote Services: SMB/Windows Admin Shares",
	//	Severity:   high,
	//	Confidence: 0.8,
	//	Phase:      LM,
	//},
	//
	//"BZAR::Lateral_Movement_T1021.002_T1570": {
	//	Origin:     "BZAR MITRE Lateral Movement",
	//	TypeHint:   "ATTACK::Lateral_Movement",
	//	Sub:        "T1021.002 Remote Services: SMB/Windows Admin Shares + T1570 Lateral Tool Transfer",
	//	Severity:   high,
	//	Confidence: 0.8,
	//	Phase:      LM,
	//},
	//
	//"BZAR::Lateral_Movement_and_Execution": {
	//	Origin:     "BZAR MITRE Lateral Movement and Execution",
	//	TypeHint:   "ATTACK::Lateral_Movement_and_Execution",
	//	Severity:   high,
	//	Confidence: 0.8,
	//	Phase:      LM,
	//},
	//
	//"BZAR::Lateral_Movement_Multiple_Attempts": {
	//	Origin:     "BZAR MITRE Lateral Movement Multiple Attempts",
	//	TypeHint:   "ATTACK::Lateral_Movement_Multiple_Attempts",
	//	Severity:   high,
	//	Confidence: 0.8,
	//	Phase:      LM,
	//},
	//
	//"BZAR::Credential_Access_T1003.006": {
	//	Origin:     "BZAR MITRE Credential Access",
	//	TypeHint:   "ATTACK::Credential_Access",
	//	Sub:        "T1003.006 OS Credential Dumping: DCSync",
	//	Severity:   low,
	//	Confidence: 0.8,
	//	Phase:      IA,
	//},
	//
	//"BZAR::Execution_T1569.002": {
	//	Origin:     "BZAR MITRE Execution",
	//	TypeHint:   "ATTACK::Execution",
	//	Sub:        "T1569.002 System Services: Service Execution",
	//	Severity:   medium,
	//	Confidence: 0.8,
	//	Phase:      Exec,
	//},
	//
	//"BZAR::Execution_T1047": {
	//	Origin:     "BZAR MITRE Execution",
	//	TypeHint:   "ATTACK::Execution",
	//	Sub:        "T1047 WMI",
	//	Severity:   medium,
	//	Confidence: 0.8,
	//	Phase:      Exec,
	//},
	//
	//"BZAR::Execution_T1053.002": {
	//	Origin:     "BZAR MITRE Execution",
	//	TypeHint:   "ATTACK::Execution",
	//	Sub:        "T1053.002 Scheduled Task/Job: At",
	//	Severity:   medium,
	//	Confidence: 0.8,
	//	Phase:      Exec,
	//},
	//
	//"BZAR::Execution_T1053.005": {
	//	Origin:     "BZAR MITRE Execution",
	//	TypeHint:   "ATTACK::Execution",
	//	Sub:        "T1053.005 Scheduled Task/Job: Scheduled Task",
	//	Severity:   low,
	//	Confidence: 0.8,
	//	Phase:      Exec,
	//},
	//
	//"BZAR_Discovery": {
	//	Origin:     "DCE RPC Discovery",
	//	TypeHint:   "ATTACK::Discovery",
	//	Severity:   medium,
	//	Confidence: 0.8,
	//	Phase:      RD,
	//},
	//
	//"BZAR::Discovery_T1016": {
	//	Origin:     "BZAR MITRE Discovery",
	//	TypeHint:   "ATTACK::Discovery",
	//	Sub:        "T1016 System Network Configuration Discovery",
	//	Severity:   medium,
	//	Confidence: 0.8,
	//	Phase:      RD,
	//},
	//
	//"BZAR::Discovery_T1018": {
	//	Origin:     "BZAR MITRE Discovery",
	//	TypeHint:   "ATTACK::Discovery",
	//	Sub:        "T1018 Remote System Discovery",
	//	Severity:   medium,
	//	Confidence: 0.8,
	//	Phase:      RD,
	//},
	//
	//"BZAR::Discovery_T1033": {
	//	Origin:     "BZAR MITRE Discovery",
	//	TypeHint:   "ATTACK::Discovery",
	//	Sub:        "T1033 System Owner/User Discovery",
	//	Severity:   medium,
	//	Confidence: 0.8,
	//	Phase:      RD,
	//},
	//
	//"BZAR::Discovery_T1049": {
	//	Origin:     "BZAR MITRE Discovery",
	//	TypeHint:   "ATTACK::Discovery",
	//	Sub:        "T1049 System Network Connections Discovery",
	//	Severity:   low,
	//	Confidence: 0.8,
	//	Phase:      RD,
	//},
	//
	//"BZAR::Discovery_T1069": {
	//	Origin:     "BZAR MITRE Discovery",
	//	TypeHint:   "ATTACK::Discovery",
	//	Sub:        "T1069 Permission Groups Discovery",
	//	Severity:   medium,
	//	Confidence: 0.8,
	//	Phase:      RD,
	//},
	//
	//"BZAR::Discovery_T1082": {
	//	Origin:     "BZAR MITRE Discovery",
	//	TypeHint:   "ATTACK::Discovery",
	//	Sub:        "T1082 System Information Discovery",
	//	Severity:   medium,
	//	Confidence: 0.8,
	//	Phase:      RD,
	//},
	//
	//"BZAR::Discovery_T1083": {
	//	Origin:     "BZAR MITRE Discovery",
	//	TypeHint:   "ATTACK::Discovery",
	//	Sub:        "T1083 File and Directory Discovery",
	//	Severity:   medium,
	//	Confidence: 0.8,
	//	Phase:      RD,
	//},
	//
	//"BZAR::Discovery_T1087": {
	//	Origin:     "BZAR MITRE Discovery",
	//	TypeHint:   "ATTACK::Discovery",
	//	Sub:        "T1087 Account Discovery",
	//	Severity:   low,
	//	Confidence: 0.8,
	//	Phase:      RD,
	//},
	//
	//"BZAR::Discovery_T1124": {
	//	Origin:     "BZAR MITRE Discovery",
	//	TypeHint:   "ATTACK::Discovery",
	//	Sub:        "T1124 System Time Discovery",
	//	Severity:   low,
	//	Confidence: 0.8,
	//	Phase:      RD,
	//},
	//
	//"BZAR::Discovery_T1135": {
	//	Origin:     "BZAR MITRE Discovery",
	//	TypeHint:   "ATTACK::Discovery",
	//	Sub:        "T1135 Network Share Discovery",
	//	Severity:   low,
	//	Confidence: 0.8,
	//	Phase:      RD,
	//},
}
