package handler

import (
	"encoding/json"
	"github.com/farhad-aman/image-authenticator-cloud-computing/publisher/db"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

type RegistrationData struct {
	Name     string `json:"name" validate:"required,max=50"`
	Email    string `json:"email" validate:"required,email"`
	National string `json:"national" validate:"required,numeric,max=10"`
	Image1   string `json:"image1" validate:"required"`
	Image2   string `json:"image2" validate:"required"`
	IP       string `json:"ip"`
}

func Register(c echo.Context) error {
	var data RegistrationData
	err := json.NewDecoder(c.Request().Body).Decode(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to parse JSON data",
		})
	}
	data.IP = c.RealIP()

	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": err.(validator.ValidationErrors),
		})
	}

	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")

	err = db.UploadToS3([]byte(data.Image1), db.EncodeNationalID(data.National)+"-1", accessKey, secretKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to upload image1 to S3",
			"details": err.Error(),
		})
	}

	err = db.UploadToS3([]byte(data.Image2), db.EncodeNationalID(data.National)+"-2", accessKey, secretKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to upload image2 to S3",
			"details": err.Error(),
		})
	}

	err = db.SaveUserData(data.Name, data.Email, data.National, data.IP)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to save user data in the database",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Registration successful",
	})
}
