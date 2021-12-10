package model

const (
	critical = 10
	high     = 8
	medium   = 5
	low      = 2
	IA       = "Initial Access"
	RD       = "Reconnaissance & Discovery"
	LM       = "Lateral Movement"
	CnC      = "Command and Control"
	Exec     = "Execution"
)

type NoticeType struct {
	Resource   string  `json:"resource"`
	Origin     string  `json:"origin"`
	TypeHint   string  `json:"type_hint"`
	Sub        string  `json:"sub"`
	Severity   float64 `json:"severity"`
	Confidence float64 `json:"confidence"`
	Phase      string  `json:"phase"`
}

var noticeNoteMap = map[string]NoticeType{
	"SQL Injection": {
		Origin:     "SQL Injection Detection",
		TypeHint:   "HTTP::SQL_Injection",
		Severity:   high,
		Confidence: 0.7,
		Phase:      IA,
	},

	"SSH Brute": {
		Origin:     "SSH Bruteforce Detection",
		TypeHint:   "SSH::Password_Guessing",
		Severity:   low,
		Confidence: 0.9,
		Phase:      IA,
	},

	"Zeek_Scan": {
		Origin:     "Scan Detection",
		TypeHint:   "Scan::",
		Severity:   low,
		Confidence: 0.8,
		Phase:      RD,
	},

	"SMB": {
		Origin:     "SMB EternalBlue Detection",
		TypeHint:   "EternalSafety::",
		Severity:   high,
		Confidence: 0.8,
		Phase:      LM,
	},

	"BZAR::Lateral_Movement_T1021.002": {
		Origin:     "BZAR MITRE Lateral Movement",
		TypeHint:   "ATTACK::Lateral_Movement",
		Sub:        "T1021.002 Remote Services: SMB/Windows Admin Shares",
		Severity:   high,
		Confidence: 0.8,
		Phase:      LM,
	},

	"BZAR::Lateral_Movement_T1021.002_T1570": {
		Origin:     "BZAR MITRE Lateral Movement",
		TypeHint:   "ATTACK::Lateral_Movement",
		Sub:        "T1021.002 Remote Services: SMB/Windows Admin Shares + T1570 Lateral Tool Transfer",
		Severity:   high,
		Confidence: 0.8,
		Phase:      LM,
	},

	"BZAR::Lateral_Movement_and_Execution": {
		Origin:     "BZAR MITRE Lateral Movement and Execution",
		TypeHint:   "ATTACK::Lateral_Movement_and_Execution",
		Severity:   high,
		Confidence: 0.8,
		Phase:      LM,
	},

	"BZAR::Lateral_Movement_Multiple_Attempts": {
		Origin:     "BZAR MITRE Lateral Movement Multiple Attempts",
		TypeHint:   "ATTACK::Lateral_Movement_Multiple_Attempts",
		Severity:   high,
		Confidence: 0.8,
		Phase:      LM,
	},

	"BZAR::Credential_Access_T1003.006": {
		Origin:     "BZAR MITRE Credential Access",
		TypeHint:   "ATTACK::Credential_Access",
		Sub:        "T1003.006 OS Credential Dumping: DCSync",
		Severity:   low,
		Confidence: 0.8,
		Phase:      IA,
	},

	"BZAR::Execution_T1569.002": {
		Origin:     "BZAR MITRE Execution",
		TypeHint:   "ATTACK::Execution",
		Sub:        "T1569.002 System Services: Service Execution",
		Severity:   medium,
		Confidence: 0.8,
		Phase:      Exec,
	},

	"BZAR::Execution_T1047": {
		Origin:     "BZAR MITRE Execution",
		TypeHint:   "ATTACK::Execution",
		Sub:        "T1047 WMI",
		Severity:   medium,
		Confidence: 0.8,
		Phase:      Exec,
	},

	"BZAR::Execution_T1053.002": {
		Origin:     "BZAR MITRE Execution",
		TypeHint:   "ATTACK::Execution",
		Sub:        "T1053.002 Scheduled Task/Job: At",
		Severity:   medium,
		Confidence: 0.8,
		Phase:      Exec,
	},

	"BZAR::Execution_T1053.005": {
		Origin:     "BZAR MITRE Execution",
		TypeHint:   "ATTACK::Execution",
		Sub:        "T1053.002 Scheduled Task/Job: Scheduled Task",
		Severity:   low,
		Confidence: 0.8,
		Phase:      Exec,
	},

	"BZAR_Discovery": {
		Origin:     "DCE RPC Discovery",
		TypeHint:   "ATTACK::Discovery",
		Severity:   medium,
		Confidence: 0.8,
		Phase:      RD,
	},

	"BZAR::Discovery_T1016": {
		Origin:     "BZAR MITRE Discovery",
		TypeHint:   "ATTACK::Discovery",
		Sub:        "T1016 System Network Configuration Discovery",
		Severity:   medium,
		Confidence: 0.8,
		Phase:      RD,
	},

	"BZAR::Discovery_T1018": {
		Origin:     "BZAR MITRE Discovery",
		TypeHint:   "ATTACK::Discovery",
		Sub:        "T1018 Remote System Discovery",
		Severity:   medium,
		Confidence: 0.8,
		Phase:      RD,
	},

	"BZAR::Discovery_T1033": {
		Origin:     "BZAR MITRE Discovery",
		TypeHint:   "ATTACK::Discovery",
		Sub:        "T1033 System Owner/User Discovery",
		Severity:   medium,
		Confidence: 0.8,
		Phase:      RD,
	},

	"BZAR::Discovery_T1049": {
		Origin:     "BZAR MITRE Discovery",
		TypeHint:   "ATTACK::Discovery",
		Sub:        "T1049 System Network Connections Discovery",
		Severity:   low,
		Confidence: 0.8,
		Phase:      RD,
	},

	"BZAR::Discovery_T1069": {
		Origin:     "BZAR MITRE Discovery",
		TypeHint:   "ATTACK::Discovery",
		Sub:        "T1069 Permission Groups Discovery",
		Severity:   medium,
		Confidence: 0.8,
		Phase:      RD,
	},

	"BZAR::Discovery_T1082": {
		Origin:     "BZAR MITRE Discovery",
		TypeHint:   "ATTACK::Discovery",
		Sub:        "T1082 System Information Discovery",
		Severity:   medium,
		Confidence: 0.8,
		Phase:      RD,
	},

	"BZAR::Discovery_T1083": {
		Origin:     "BZAR MITRE Discovery",
		TypeHint:   "ATTACK::Discovery",
		Sub:        "T1083 File and Directory Discovery",
		Severity:   medium,
		Confidence: 0.8,
		Phase:      RD,
	},

	"BZAR::Discovery_T1087": {
		Origin:     "BZAR MITRE Discovery",
		TypeHint:   "ATTACK::Discovery",
		Sub:        "T1087 Account Discovery",
		Severity:   low,
		Confidence: 0.8,
		Phase:      RD,
	},

	"BZAR::Discovery_T1124": {
		Origin:     "BZAR MITRE Discovery",
		TypeHint:   "ATTACK::Discovery",
		Sub:        "T1124 System Time Discovery",
		Severity:   low,
		Confidence: 0.8,
		Phase:      RD,
	},

	"BZAR::Discovery_T1135": {
		Origin:     "BZAR MITRE Discovery",
		TypeHint:   "ATTACK::Discovery",
		Sub:        "T1135 Network Share Discovery",
		Severity:   low,
		Confidence: 0.8,
		Phase:      RD,
	},
}
