package controllers

import (
	"strconv"

	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

// InvoiceController ...
type InvoiceController struct{}

var invoiceModel = new(models.InvoiceModel)
var invoiceForm = new(forms.InvoiceForm)

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
// @Router /invoice [post]
func (ctrl InvoiceController) Create(c *gin.Context) {
	userID := getUserID(c)

	var form forms.CreateInvoiceForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := invoiceForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	id, err := invoiceModel.Create(userID, form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invoice could not be created"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice created", "id": id})
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
// @Router /invoice [post]
func (ctrl InvoiceController) All(c *gin.Context) {
	userID := getUserID(c)

	results, err := invoiceModel.All(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get invoices"})
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
// @Router /invoice [post]
func (ctrl InvoiceController) One(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")

	getID, err := strconv.ParseInt(id, 10, 64)
	if getID == 0 || err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter"})
		return
	}

	data, err := invoiceModel.One(userID, getID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invoice not found"})
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
// @Router /invoice [post]
func (ctrl InvoiceController) Update(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")

	getID, err := strconv.ParseInt(id, 10, 64)
	if getID == 0 || err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter"})
		return
	}

	var form forms.CreateInvoiceForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := invoiceForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	err = invoiceModel.Update(userID, getID, form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Invoice could not be updated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice updated"})
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
// @Router /invoice [post]
func (ctrl InvoiceController) Delete(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")

	getID, err := strconv.ParseInt(id, 10, 64)
	if getID == 0 || err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter"})
		return
	}

	err = invoiceModel.Delete(userID, getID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Invoice could not be deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice deleted"})

}
