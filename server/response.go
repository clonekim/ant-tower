package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// respondError sends a standardized error response.
func respondError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status":  "error",
		"message": message,
	})
}

// respondSuccess sends a standardized success response with data.
func respondSuccess(c *gin.Context, message string, data gin.H) {
	// Base response
	response := gin.H{
		"status":  "success",
		"message": message,
	}

	// Merge data into response (Flattened structure)
	for k, v := range data {
		response[k] = v
	}

	c.JSON(http.StatusOK, response)
}

// bindJSONOrError binds the request body and handles the error automatically.
// Returns true if binding was successful, false otherwise.
func bindJSONOrError(c *gin.Context, obj interface{}) bool {
	if err := c.BindJSON(obj); err != nil {
		respondError(c, http.StatusBadRequest, "Invalid request body")
		return false
	}
	return true
}
