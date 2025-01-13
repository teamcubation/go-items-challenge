package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teamcubation/go-items-challenge/internal/domain/item"
)

func ValidateItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newItem item.Item
		if err := c.ShouldBindJSON(&newItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			c.Abort()
			return
		}
		c.Set("newItem", newItem)
		c.Next()
	}
}
