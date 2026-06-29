package handler

import ( 
	"net/http" 
	"strconv" 
	"spotSync-golang-Project/internal/dto" 
	"spotSync-golang-Project/internal/repository" 
	"spotSync-golang-Project/internal/service" 

	"github.com/go-playground/validator/v10" 
	"github.com/labstack/echo/v4" 
	"gorm.io/gorm" 
)

type ReservationHandler struct { 
	resService service.ReservationService 
	validate   *validator.Validate 
}

func NewReservationHandler(resService service.ReservationService, validate *validator.Validate) *ReservationHandler { 
	return &ReservationHandler{resService, validate} 
}

func (h *ReservationHandler) ReserveSpot(c echo.Context) error { 
	userID := c.Get("user_id").(uint) 

	var req dto.CreateReservationRequest 
	if err := c.Bind(&req); err != nil { 
		return JSONError(c, http.StatusBadRequest, "Invalid request payload", err.Error()) 
	}

	if err := h.validate.Struct(req); err != nil { 
		return JSONError(c, http.StatusBadRequest, "Validation failed", err.Error()) 
	}

	res, err := h.resService.ReserveSpot(userID, req) 
	if err != nil { 
		if err == repository.ErrZoneFull { 
			return JSONError(c, http.StatusConflict, "Reservation failed", err.Error()) 
		}
		if err == gorm.ErrRecordNotFound { 
			return JSONError(c, http.StatusNotFound, "Zone not found", "Invalid zone ID") 
		}
		return JSONError(c, http.StatusInternalServerError, "Internal server error", err.Error()) 
	}

	return JSONSuccess(c, http.StatusCreated, "Reservation confirmed successfully", res) 
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error { 
	userID := c.Get("user_id").(uint) 

	res, err := h.resService.GetUserReservations(userID) 
	if err != nil { 
		return JSONError(c, http.StatusInternalServerError, "Failed to fetch reservations", err.Error()) 
	}

	return JSONSuccess(c, http.StatusOK, "My reservations retrieved successfully", res) 
}

func (h *ReservationHandler) GetAllReservations(c echo.Context) error { 
	res, err := h.resService.GetAllReservations() 
	if err != nil { 
		return JSONError(c, http.StatusInternalServerError, "Failed to fetch reservations", err.Error()) 
	}

	return JSONSuccess(c, http.StatusOK, "All reservations retrieved successfully", res) 
}

func (h *ReservationHandler) CancelReservation(c echo.Context) error { 
	userID := c.Get("user_id").(uint) 
	role := c.Get("role").(string) 

	idParam := c.Param("id") 
	resID, err := strconv.ParseUint(idParam, 10, 32) 
	if err != nil { 
		return JSONError(c, http.StatusBadRequest, "Invalid reservation ID", err.Error()) 
	}

	err = h.resService.CancelReservation(userID, uint(resID), role) 
	if err != nil { 
		if err.Error() == "forbidden" { 
			return JSONError(c, http.StatusForbidden, "Cancel failed", "You can only cancel your own reservations") 
		}
		if err == gorm.ErrRecordNotFound { 
			return JSONError(c, http.StatusNotFound, "Reservation not found", "Invalid reservation ID") 
		}
		if err.Error() == "reservation is not active" { 
			return JSONError(c, http.StatusBadRequest, "Cancel failed", err.Error()) 
		}
		return JSONError(c, http.StatusInternalServerError, "Internal server error", err.Error()) 
	}

	return JSONSuccess(c, http.StatusOK, "Reservation cancelled successfully", nil) 
}
