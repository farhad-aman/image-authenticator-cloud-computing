package handler

import (
	"github.com/farhad-aman/image-authenticator-cloud-computing/publisher/datastore"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type StatusResponse struct {
	State string `json:"state"`
}

type StatusRequest struct {
	National string `json:"national" validate:"required,numeric,max=10"`
}

func Status(c echo.Context) error {
	req := new(StatusRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": err.(validator.ValidationErrors),
		})
	}

	user, err := datastore.FetchUserData(req.National)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	if user.IP != c.RealIP() {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Access denied. Invalid IP address",
		})
	}

	return c.JSON(http.StatusOK, user)
}
