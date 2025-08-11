package userhttp

import (
	"errors"
	"net/http"
	"strings"

	"github.com/frenswifbenefits/myfren/internal/dto"
	"github.com/frenswifbenefits/myfren/internal/entity"

	"github.com/gin-gonic/gin"
)

// auth performs user auth
func (s *ServerImpl) auth(c *gin.Context) (*entity.Customer, bool) {
	customer, err := s.getCustomer(c)
	if err != nil {
		c.JSON(http.StatusForbidden, dto.MakeForbiddenAPIResponse(err))
		return nil, false
	}
	return customer, true
}

func (s *ServerImpl) getCustomer(c *gin.Context) (*entity.Customer, error) {
	authHeader := c.Request.Header["Authorization"]
	if len(authHeader) == 0 {
		return nil, errors.New("empty auth header")
	}
	token := strings.Trim(authHeader[0], "\n\r\t ")
	token = strings.TrimPrefix(token, "Bearer")
	token = strings.Trim(token, "\n\r\t ")

	customer, ok := s.cp.GetByToken(token)
	if !ok {
		return nil, errors.New("user not found")
	}

	return customer, nil
}
