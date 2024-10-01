package authentication

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/helpers"
)

func Authenticate(c *gin.Context) (*helpers.SignedDetails, error) {

	clientToken := c.Request.Header.Get("Authorization")

	if clientToken == "" {
		return nil, errors.New("token not found")
	}

	claims, err := helpers.ValidateToken(clientToken)

	if err != nil {
		return nil, err
	}

	return claims, nil
}
