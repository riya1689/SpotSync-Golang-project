package service

import (
	"spotSync-golang-Project/internal/dto"
	"spotSync-golang-Project/internal/models"
	"spotSync-golang-Project/internal/repository"
)

type ZoneService interface {
	CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) 
	GetAllZones() ([]dto.ZoneResponse, error) 
	GetZoneByID(id uint) (*dto.ZoneResponse, error)
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
}

func NewZoneService(repo repository.ZoneRepository) ZoneService {
	return &zoneService{repo}
}

func (s *zoneService) CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.zoneRepo.CreateZone(&zone); err != nil { 
		return nil, err 
	}

	return &dto.ZoneResponse{ 
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
		UpdatedAt:     zone.UpdatedAt,
	}, nil
}

func (s *zoneService) GetAllZones() ([]dto.ZoneResponse, error) { 
	return s.zoneRepo.GetAllZones() 
}

func (s *zoneService) GetZoneByID(id uint) (*dto.ZoneResponse, error) { 
	return s.zoneRepo.GetZoneByID(id) 
}
