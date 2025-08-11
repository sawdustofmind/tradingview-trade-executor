package userhttp

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/frenswifbenefits/myfren/internal/entity"
	"net/http"
	"slices"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"

	"github.com/frenswifbenefits/myfren/internal/dto"
	"github.com/frenswifbenefits/myfren/internal/openapi"
)

func (s *ServerImpl) GetV1UserSettings(c *gin.Context) {
	customer, ok := s.auth(c)
	if !ok {
		return
	}

	resp := api_types.SettingsResponse{
		Code:    dto.SuccessCode,
		Message: dto.SuccessMessage,
	}

	resp.Data.Settings.Username = customer.Username
	resp.Data.Settings.LegalName = customer.LegalName
	resp.Data.Settings.Gender = customer.Gender
	resp.Data.Settings.Country = customer.Country
	resp.Data.Settings.PhoneNumber = customer.PhoneNumber
	resp.Data.Settings.ImageBase64 = customer.ImageBase64
	if customer.BybitApiKey != nil {
		resp.Data.Settings.BybitApiKey = *customer.BybitApiKey
	}
	if customer.BybitTestApiKey != nil {
		resp.Data.Settings.BybitTestApiKey = *customer.BybitTestApiKey
	}
	if customer.BybitApiSecret != nil {
		resp.Data.Settings.BybitApiSecret = *customer.BybitApiSecret
		if len(resp.Data.Settings.BybitApiSecret) >= 4 {
			resp.Data.Settings.BybitApiSecret = resp.Data.Settings.BybitApiSecret[:4] + "***"
		}
	}
	if customer.BybitTestApiSecret != nil {
		resp.Data.Settings.BybitTestApiSecret = *customer.BybitTestApiSecret
	}

	c.JSON(http.StatusOK, resp)
}

func (s *ServerImpl) PutV1UserSettings(c *gin.Context) {
	customer, ok := s.auth(c)
	if !ok {
		return
	}

	req := &api_types.PutV1UserSettingsJSONBody{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(err))
		return
	}

	if req.Country != nil {
		customer.Country = *req.Country
	}
	if req.Gender != nil {
		customer.Gender = *req.Gender
	}
	if req.ImageBase64 != nil {
		customer.ImageBase64 = *req.ImageBase64
	}
	if req.LegalName != nil {
		customer.LegalName = *req.LegalName
	}
	if req.PhoneNumber != nil {
		customer.PhoneNumber = *req.PhoneNumber
	}

	isBybitApiKeyChanged := false
	isBybitTestApiKeyChanged := false
	if req.BybitApiKey != nil {
		customer.BybitApiKey = req.BybitApiKey
		isBybitApiKeyChanged = true
	}
	if req.BybitTestApiKey != nil {
		customer.BybitTestApiKey = req.BybitTestApiKey
		isBybitTestApiKeyChanged = true
	}
	if req.BybitApiSecret != nil {
		customer.BybitApiSecret = req.BybitApiSecret
		isBybitApiKeyChanged = true
	}
	if req.BybitTestApiSecret != nil {
		customer.BybitTestApiSecret = req.BybitTestApiSecret
		isBybitTestApiKeyChanged = true
	}

	bybitKeyErased := customer.BybitApiKey != nil && len(*customer.BybitApiKey) == 0 &&
		customer.BybitApiSecret != nil && len(*customer.BybitApiSecret) == 0
	if isBybitApiKeyChanged && !bybitKeyErased {
		if customer.BybitApiKey == nil || len(*customer.BybitApiKey) == 0 {
			c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(errors.New("empty bybit api key")))
			return
		}
		if customer.BybitApiSecret == nil || len(*customer.BybitApiSecret) == 0 {
			c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(errors.New("empty bybit api secret")))
			return
		}
		healthy, err := s.bs.Healthcheck(c, *customer, false)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.MakeBadApiKeyAPIResponse(err))
			return
		}
		if !healthy {
			c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(errors.New("api creds is unhealthy")))
			return
		}
	}
	bybitTestKeyErased := customer.BybitTestApiKey != nil && len(*customer.BybitTestApiKey) == 0 &&
		customer.BybitTestApiSecret != nil && len(*customer.BybitTestApiSecret) == 0
	if isBybitTestApiKeyChanged && !bybitTestKeyErased {
		if customer.BybitTestApiKey == nil || len(*customer.BybitTestApiKey) == 0 {
			c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(errors.New("empty bybit test api key")))
			return
		}
		if customer.BybitTestApiSecret == nil || len(*customer.BybitTestApiSecret) == 0 {
			c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(errors.New("empty bybit test api secret")))
			return
		}
		healthy, err := s.bs.Healthcheck(c, *customer, true)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.MakeBadApiKeyAPIResponse(err))
			return
		}
		if !healthy {
			c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(errors.New("api test creds is unhealthy")))
			return
		}
	}

	err = s.repository.UpdateCustomer(*customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeBadRequestAPIResponse(err))
		return
	}

	c.JSON(http.StatusOK, dto.MakeSuccessAPIResponse())
}

func (s *ServerImpl) PostV1UserRegister(c *gin.Context) {
	req := &api_types.PostV1UserRegisterJSONBody{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(err))
		return
	}
	if req.Username == "" {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(errors.New("name param must exist")))
		return
	}
	if req.Password == "" {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(errors.New("password param must exist")))
		return
	}
	if req.InviteToken == "" {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(errors.New("invite_token param must exist")))
		return
	}

	inviteTokens, err := s.repository.GetInviteTokens()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}
	if !slices.Contains(inviteTokens, req.InviteToken) {
		c.JSON(http.StatusBadRequest, dto.MakeErrorAPIResponse(errors.New("invite_token doesn't exist")))
		return
	}

	_, getUserError := s.repository.GetCustomerByUsername(req.Username)
	if getUserError == nil {
		c.JSON(http.StatusBadRequest, dto.MakeErrorAPIResponse(errors.New("customer username is already in use")))
		return
	}
	if errors.Is(getUserError, sql.ErrNoRows) {
		getUserError = nil
	}
	if getUserError != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(getUserError))
		return
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}
	_, err = s.repository.InsertCustomer(entity.Customer{
		Username: req.Username,
		Password: string(passHash),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}

	err = s.repository.DeleteInviteToken(req.InviteToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}

	c.JSON(http.StatusOK, dto.MakeSuccessAPIResponse())
}

func (s *ServerImpl) PostV1UserLogin(c *gin.Context) {
	req := &api_types.PostV1UserLoginJSONBody{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(err))
		return
	}
	if req.Username == "" {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(errors.New("name param must exist")))
		return
	}
	if req.Password == "" {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(errors.New("password param must exist")))
		return
	}

	customer, ok := s.cp.GetByUsername(req.Username)
	if !ok {
		c.JSON(http.StatusNotFound, dto.MakeErrorAPIResponse(errors.New("user not found")))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusNotFound, dto.MakeErrorAPIResponse(errors.New("user not found")))
		return
	}

	bts := make([]byte, 32)
	_, err = rand.Read(bts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}
	token := hex.EncodeToString(bts)
	s.cp.AttachToken(token, customer)

	resp := api_types.TokenResponse{
		Code:    dto.SuccessCode,
		Message: dto.SuccessMessage,
	}
	resp.Data.Token = token

	c.JSON(http.StatusOK, resp)
}
