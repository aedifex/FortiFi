package database

type User struct {
	Id 			string	`json:"id" sql:"id"`
 	FirstName 	string	`json:"first_name" sql:"first_name"`
 	LastName 	string 	`json:"last_name" sql:"last_name"`
 	Email 		string	`json:"email" sql:"email"`
 	Password 	string	`json:"password" sql:"password"`
	FcmToken 	string 	`json:"fcm_token" sql:"fcm_token"`
}

type RefreshToken struct {
	Token	 	string	`sql:"token"`
    Id 			string	`sql:"id"`
    Expires 	string	`sql:"id"`
}

type Event struct {
	ThreatId 	int		`json:"threat_id" sql:"threat_id"`
	Id 			string	`json:"id" sql:"id"`
    Details 	string	`json:"details" sql:"details"`
    TS 			string	`json:"ts" sql:"ts"`
    Expires 	string	`json:"expires" sql:"expires"`
	Type 		string	`json:"type" sql:"event_type"`
	SrcIP 		string	`json:"src" sql:"src_ip"`
	DstIP 		string	`json:"dst" sql:"dst_ip"`
}

type WeeklyDistribution struct {
	Benign			int	`sql:"benign_count"`
	PortScan		int	`sql:"port_scan_count"`
	DDoS			int	`sql:"ddos_count"`
	PrevWeekTotal	int	`sql:"prev_week_total"`
}

type Device struct {
	Id 			int		`json:"id" sql:"id"`
	Name 		string	`json:"name" sql:"name"`
	IpAddress 	string	`json:"ip_address" sql:"ip_address"`
	MacAddress 	string	`json:"mac_address" sql:"mac_address"`
	UserId 		string	`json:"user_id" sql:"user_id"`
	DateAdded 	string	`json:"date_added" sql:"date_added"`
	IncidentCount	int		`json:"incident_count" sql:"incident_count"`
}
