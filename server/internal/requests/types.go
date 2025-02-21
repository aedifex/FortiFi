package requests

import (
	db "github.com/aedifex/FortiFi/internal/database"
)
type UpdateFcmRequest struct {
	FcmToken	string	`json:"fcm_token"`
}

type CreateUserRequest struct {
	User	*db.User	`json:"user"`
}

type LoginUserRequest struct {
	User 	*db.User	`json:"user"`		
}

type NotifyIntrusionRequest struct {
	Event	*db.Event	`json:"event"`
}

type PiInitRequest struct {
	Id	string	`json:"id"`
}

type UpdateWeeklyDistributionRequest struct {
	Benign		int	`json:"benign"`
	PortScan	int	`json:"port_scan"`
	DDoS		int	`json:"ddos"`
}

type ResetWeeklyDistributionRequest struct {
	WeekTotal	int	`json:"week_total"`
}

type AddDeviceRequest struct {
	Name 		string	`json:"name"`
	IpAddress	string	`json:"ip_address"`
	MacAddress	string	`json:"mac_address"`
}