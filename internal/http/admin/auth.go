package adminhttp

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/frenswifbenefits/myfren/internal/config"
	"github.com/frenswifbenefits/myfren/internal/dto"
	api_types "github.com/frenswifbenefits/myfren/internal/openapi"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (s *ServerImpl) PostV1AdminLogin(c *gin.Context) {
	req := &api_types.PostV1AdminLoginJSONBody{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(err))
		return
	}

	if req.Password == "" {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(errors.New("password param must exist")))
		return
	}

	idx := slices.IndexFunc(s.authConf, func(user config.AdminUserConfig) bool {
		return user.Username == req.Username
	})
	if idx == -1 {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(errors.New("username or password is invalid")))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(s.authConf[idx].Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(err))
		return
	}

	bts := make([]byte, 32)
	_, err = rand.Read(bts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}

	token := hex.EncodeToString(bts)
	resp := api_types.TokenResponse{
		Code:    dto.SuccessCode,
		Message: dto.SuccessMessage,
	}
	resp.Data.Token = token

	s.authTokensMu.Lock()
	s.authTokens[token] = struct{}{}
	s.authTokensMu.Unlock()
	c.JSON(http.StatusOK, resp)
}

// auth performs user auth
func (s *ServerImpl) auth(c *gin.Context) bool {
	authHeader := c.Request.Header["Authorization"]
	if len(authHeader) == 0 {
		c.JSON(http.StatusForbidden, dto.MakeForbiddenAPIResponse(errors.New("empty auth header")))
		return false
	}
	token := strings.Trim(authHeader[0], "\n\r\t ")
	token = strings.TrimPrefix(token, "Bearer")
	token = strings.Trim(token, "\n\r\t ")

	s.authTokensMu.RLock()
	_, tokenFound := s.authTokens[token]
	s.authTokensMu.RUnlock()

	if !tokenFound {
		c.JSON(http.StatusForbidden, dto.MakeForbiddenAPIResponse(errors.New("empty auth header")))
		return false
	}

	return true
}
