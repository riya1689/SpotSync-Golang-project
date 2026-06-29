package service 

import (
	"errors" 
	"spotSync-golang-Project/internal/dto" 
	"spotSync-golang-Project/internal/models" 
	"spotSync-golang-Project/internal/repository" 
)

type ReservationService interface { 
	ReserveSpot(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) 
	GetUserReservations(userID uint) ([]dto.ReservationResponse, error) 
	GetAllReservations() ([]dto.ReservationResponse, error) 
	CancelReservation(userID uint, reservationID uint, role string) error 
}

type reservationService struct { 
	resRepo repository.ReservationRepository 
}

func NewReservationService(repo repository.ReservationRepository) ReservationService { 
	return &reservationService{repo} 
}

func (s *reservationService) ReserveSpot(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	reservation := models.Reservation{ 
		UserID:       userID, 
		ZoneID:       req.ZoneID, 
		LicensePlate: req.LicensePlate, 
		Status:       "active", 
	}

	if err := s.resRepo.ReserveSpot(&reservation); err != nil { 
		return nil, err 
	}

	
	res, err := s.resRepo.GetReservationByID(reservation.ID) 
	if err != nil { 
		return nil, err 
	}

	return s.mapToResponse(res), nil 
}

func (s *reservationService) GetUserReservations(userID uint) ([]dto.ReservationResponse, error) { 
	reservations, err := s.resRepo.GetUserReservations(userID) 		
	if err != nil { 
		return nil, err 
	}

	var res []dto.ReservationResponse 
	for _, r := range reservations { 
		res = append(res, *s.mapToResponse(&r)) 
	}
	if res == nil { 
		res = []dto.ReservationResponse{} 
	}
	return res, nil 
}

func (s *reservationService) GetAllReservations() ([]dto.ReservationResponse, error) { 
	reservations, err := s.resRepo.GetAllReservations() 
	if err != nil { 
		return nil, err 
	}

	var res []dto.ReservationResponse 
	for _, r := range reservations { 
		res = append(res, *s.mapToResponse(&r)) 
	}
	if res == nil { 
		res = []dto.ReservationResponse{} 
	}
	return res, nil 
}

func (s *reservationService) CancelReservation(userID uint, reservationID uint, role string) error { 
	res, err := s.resRepo.GetReservationByID(reservationID) 
	if err != nil { 
		return err 
	}

	if role != "admin" && res.UserID != userID { 
		return errors.New("forbidden") 
	}

	if res.Status != "active" { 
		return errors.New("reservation is not active") 
	}

	return s.resRepo.UpdateReservationStatus(reservationID, "cancelled")
}

func (s *reservationService) mapToResponse(r *models.Reservation) *dto.ReservationResponse { 
	var zoneResp *dto.ZoneResponse 
	if r.Zone.ID != 0 {
		zoneResp = &dto.ZoneResponse{ 
			ID:            r.Zone.ID,
			Name:          r.Zone.Name,
			Type:          r.Zone.Type,
			TotalCapacity: r.Zone.TotalCapacity,
			PricePerHour:  r.Zone.PricePerHour,
		}
	}
	return &dto.ReservationResponse{ 
		ID:           r.ID,
		UserID:       r.UserID,
		ZoneID:       r.ZoneID,
		LicensePlate: r.LicensePlate,
		Status:       r.Status,
		Zone:         zoneResp,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}
