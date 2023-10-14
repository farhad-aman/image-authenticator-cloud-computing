package handler

import (
	"github.com/farhad-aman/image-authenticator-cloud-computing/publisher/datastore"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,max=50"`
	Email    string `json:"email" validate:"required,email"`
	National string `json:"national" validate:"required,numeric,max=10"`
	Image1   string `json:"image1" validate:"required"`
	Image2   string `json:"image2" validate:"required"`
	IP       string `json:"ip"`
}

func Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":  "Failed to parse JSON data",
			"detail": err.Error(),
		})
	}
	req.IP = c.RealIP()

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": err.(validator.ValidationErrors),
		})
	}

	err := datastore.SaveUserData(req.Name, req.Email, datastore.EncodeNationalID(req.National), req.IP)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":  "Failed to save user data in the database",
			"detail": err.Error(),
		})
	}

	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")

	err = datastore.UploadToS3([]byte(req.Image1), datastore.EncodeNationalID(req.National)+"-1", accessKey, secretKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to upload image1 to S3",
			"details": err.Error(),
		})
	}

	err = datastore.UploadToS3([]byte(req.Image2), datastore.EncodeNationalID(req.National)+"-2", accessKey, secretKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to upload image2 to S3",
			"details": err.Error(),
		})
	}

	err = datastore.SendNationalToRabbit(datastore.EncodeNationalID(req.National))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to send national ID to RabbitMQ",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Registration successful",
	})
}
