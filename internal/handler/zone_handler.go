package handler

import (
	"net/http"
	"strconv"
	"spotSync-golang-Project/internal/dto"
	"spotSync-golang-Project/internal/service" 

	"github.com/go-playground/validator/v10" 
	"github.com/labstack/echo/v4" 
	"gorm.io/gorm" 
)

type ZoneHandler struct {
	zoneService service.ZoneService 
	validate    *validator.Validate
}

func NewZoneHandler(zoneService service.ZoneService, validate *validator.Validate) *ZoneHandler {
	return &ZoneHandler{zoneService, validate}
}

func (h *ZoneHandler) CreateZone(c echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return JSONError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
	}

	if err := h.validate.Struct(req); err != nil {
		return JSONError(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	res, err := h.zoneService.CreateZone(req)
	if err != nil {
		return JSONError(c, http.StatusInternalServerError, "Failed to create zone", err.Error())
	}

	return JSONSuccess(c, http.StatusCreated, "Parking zone created successfully", res)
}

func (h *ZoneHandler) GetAllZones(c echo.Context) error {
	zones, err := h.zoneService.GetAllZones()
	if err != nil {
		return JSONError(c, http.StatusInternalServerError, "Failed to fetch zones", err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "Parking zones retrieved successfully", zones)
}

func (h *ZoneHandler) GetZoneByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "Invalid zone ID", err.Error())
	}

	zone, err := h.zoneService.GetZoneByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return JSONError(c, http.StatusNotFound, "Zone not found", "Zone ID does not exist")
		}
		return JSONError(c, http.StatusInternalServerError, "Failed to fetch zone", err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "Parking zone retrieved successfully", zone)
}
