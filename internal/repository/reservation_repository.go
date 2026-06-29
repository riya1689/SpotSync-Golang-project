package repository 

import ( 
	"errors" 
	"spotSync-golang-Project/internal/models" 

	"gorm.io/gorm" 
	"gorm.io/gorm/clause" 
)

var ErrZoneFull = errors.New("parking zone is full") 

type ReservationRepository interface { 
	ReserveSpot(reservation *models.Reservation) error 
	GetUserReservations(userID uint) ([]models.Reservation, error) 
	GetAllReservations() ([]models.Reservation, error) 
	GetReservationByID(id uint) (*models.Reservation, error) 
	UpdateReservationStatus(id uint, status string) error 
}

type reservationRepository struct { 
	db *gorm.DB 
}

func NewReservationRepository(db *gorm.DB) ReservationRepository { 
	return &reservationRepository{db} 
}

func (r *reservationRepository) ReserveSpot(reservation *models.Reservation) error { 
	return r.db.Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone 
		
		// 1. Lock the row for update to prevent race conditions 
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, reservation.ZoneID).Error; err != nil { 
			return err 
		}

		// 2. Count current 'active' reservations for this zone 
		var activeCount int64 
		if err := tx.Model(&models.Reservation{}).Where("zone_id = ? AND status = ?", reservation.ZoneID, "active").Count(&activeCount).Error; err != nil { 
			return err 
		}

		// 3. Check if active_count < zone.total_capacity 
		if activeCount >= int64(zone.TotalCapacity) { 
			return ErrZoneFull 
		}

		// 4. Create reservation 
		if err := tx.Create(reservation).Error; err != nil {
			return err
		}

		return nil 
	})
}

func (r *reservationRepository) GetUserReservations(userID uint) ([]models.Reservation, error) { 
	var reservations []models.Reservation 
	err := r.db.Preload("Zone").Where("user_id = ?", userID).Find(&reservations).Error 
	return reservations, err
}

func (r *reservationRepository) GetAllReservations() ([]models.Reservation, error) { 
	var reservations []models.Reservation 
	err := r.db.Preload("User").Preload("Zone").Find(&reservations).Error 
	return reservations, err
}

func (r *reservationRepository) GetReservationByID(id uint) (*models.Reservation, error) { 
	var reservation models.Reservation 
	err := r.db.Preload("Zone").First(&reservation, id).Error 
	return &reservation, err
}

func (r *reservationRepository) UpdateReservationStatus(id uint, status string) error { 
	return r.db.Model(&models.Reservation{}).Where("id = ?", id).Update("status", status).Error 
}
