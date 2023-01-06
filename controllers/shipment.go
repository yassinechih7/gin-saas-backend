package controllers

import (
	"strconv"

	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

// ShipmentController ...
type ShipmentController struct{}

var shipmentModel = new(models.ShipmentModel)
var shipmentForm = new(forms.ShipmentForm)

//Create ...
// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /shipment [post]
func (ctrl ShipmentController) Create(c *gin.Context) {
	userID := getUserID(c)

	var form forms.CreateShipmentForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := shipmentForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	id, err := shipmentModel.Create(userID, form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Shipment could not be created"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipment created", "id": id})
}

// All ...
// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /shipment [post]
func (ctrl ShipmentController) All(c *gin.Context) {
	userID := getUserID(c)

	results, err := shipmentModel.All(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get shipments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}

// One ...
// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /shipment [post]
func (ctrl ShipmentController) One(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")

	getID, err := strconv.ParseInt(id, 10, 64)
	if getID == 0 || err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter"})
		return
	}

	data, err := shipmentModel.One(userID, getID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Shipment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// Update ...
// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /shipment [post]
func (ctrl ShipmentController) Update(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")

	getID, err := strconv.ParseInt(id, 10, 64)
	if getID == 0 || err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter"})
		return
	}

	var form forms.CreateShipmentForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := shipmentForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	err = shipmentModel.Update(userID, getID, form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Shipment could not be updated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipment updated"})
}

// Delete ...
// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /shipment [post]
func (ctrl ShipmentController) Delete(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")

	getID, err := strconv.ParseInt(id, 10, 64)
	if getID == 0 || err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter"})
		return
	}

	err = shipmentModel.Delete(userID, getID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Shipment could not be deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipment deleted"})

}
