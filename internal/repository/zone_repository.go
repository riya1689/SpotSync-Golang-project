package repository 

import (
	"spotSync-golang-Project/internal/dto" 
	"spotSync-golang-Project/internal/models" 

	"gorm.io/gorm" 
)

type ZoneRepository interface { 
	CreateZone(zone *models.ParkingZone) error 
	GetAllZones() ([]dto.ZoneResponse, error) 
	GetZoneByID(id uint) (*dto.ZoneResponse, error) 
}

type zoneRepository struct { 
	db *gorm.DB 
}

func NewZoneRepository(db *gorm.DB) ZoneRepository { 
	return &zoneRepository{db} 
}

func (r *zoneRepository) CreateZone(zone *models.ParkingZone) error { 
	return r.db.Create(zone).Error 
}

func (r *zoneRepository) GetAllZones() ([]dto.ZoneResponse, error) { 
	var zones []dto.ZoneResponse 
	query := ` 
		SELECT 
			pz.id, pz.name, pz.type, pz.total_capacity, pz.price_per_hour, pz.created_at, pz.updated_at,
			(pz.total_capacity - COALESCE(r.active_count, 0)) AS available_spots
		FROM parking_zones pz
		LEFT JOIN (
			SELECT zone_id, COUNT(*) as active_count
			FROM reservations
			WHERE status = 'active'
			GROUP BY zone_id
		) r ON pz.id = r.zone_id
	`
	err := r.db.Raw(query).Scan(&zones).Error 
	return zones, err 
}

func (r *zoneRepository) GetZoneByID(id uint) (*dto.ZoneResponse, error) { 
	var zone dto.ZoneResponse 
	query := `
		SELECT 
			pz.id, pz.name, pz.type, pz.total_capacity, pz.price_per_hour, pz.created_at, pz.updated_at,
			(pz.total_capacity - COALESCE(r.active_count, 0)) AS available_spots
		FROM parking_zones pz
		LEFT JOIN (
			SELECT zone_id, COUNT(*) as active_count
			FROM reservations
			WHERE status = 'active'
			GROUP BY zone_id
		) r ON pz.id = r.zone_id
		WHERE pz.id = ?
	` 
	err := r.db.Raw(query, id).Scan(&zone).Error 
	if err != nil { 
		return nil, err 
	}
	if zone.ID == 0 { 
		return nil, gorm.ErrRecordNotFound 
	}
	return &zone, nil 
}
