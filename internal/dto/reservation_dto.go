package dto 

import "time" 

type CreateReservationRequest struct { 
	ZoneID       uint   `json:"zone_id" validate:"required"` 
	LicensePlate string `json:"license_plate" validate:"required,max=15"` 
}

type ReservationResponse struct { 
	ID           uint          `json:"id"`
	UserID       uint          `json:"user_id,omitempty"`
	ZoneID       uint          `json:"zone_id,omitempty"`
	LicensePlate string        `json:"license_plate"` 
	Status       string        `json:"status"` 
	Zone         *ZoneResponse `json:"zone,omitempty"` 
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at,omitempty"`
}
