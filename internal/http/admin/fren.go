package adminhttp

import (
	"net/http"

	"github.com/frenswifbenefits/myfren/internal/dto"
	"github.com/frenswifbenefits/myfren/internal/entity"
	"github.com/frenswifbenefits/myfren/internal/http/converters"
	api_types "github.com/frenswifbenefits/myfren/internal/openapi"
	"github.com/gin-gonic/gin"
)

func (s *ServerImpl) DeleteV1AdminFren(c *gin.Context) {
	if !s.auth(c) {
		return
	}

	req := &api_types.DeleteV1AdminFrenJSONBody{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(err))
		return
	}

	err = s.repository.DeleteFren(req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}

	c.JSON(http.StatusOK, dto.MakeSuccessAPIResponse())
}

func (s *ServerImpl) PostV1AdminFren(c *gin.Context) {
	if !s.auth(c) {
		return
	}

	req := &api_types.PostV1AdminFrenJSONBody{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(err))
		return
	}

	portfolio := entity.Fren{
		Name: req.Name,
	}
	for _, id := range req.PortfolioIds {
		portfolio.Portfolios = append(portfolio.Portfolios, entity.Portfolio{
			Id: id,
		})
	}
	id, err := s.repository.InsertFren(portfolio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}

	resp := api_types.IdResponse{
		Code:    dto.SuccessCode,
		Message: dto.SuccessMessage,
	}
	resp.Data.Id = id
	c.JSON(http.StatusOK, resp)
}

func (s *ServerImpl) PutV1AdminFren(c *gin.Context) {
	if !s.auth(c) {
		return
	}

	req := &api_types.PutV1AdminFrenJSONBody{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(err))
		return
	}

	fren, err := s.repository.GetFrenById(req.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MakeErrorAPIResponse(err))
		return
	}
	if req.Description != nil {
		fren.Description = *req.Description
	}
	if req.ImageBase64 != nil {
		fren.ImageBase64 = *req.ImageBase64
	}
	if req.Name != nil {
		fren.Name = *req.Name
	}

	// TODO: update list of portfolios too
	err = s.repository.UpdateFren(fren)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}

	c.JSON(http.StatusOK, dto.MakeSuccessAPIResponse())
}

func (s *ServerImpl) GetV1AdminFrenList(c *gin.Context) {
	if !s.auth(c) {
		return
	}

	portfolios, err := s.repository.GetFrens()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
		return
	}

	resp := api_types.FrenListResponse{
		Code:    dto.SuccessCode,
		Message: dto.SuccessMessage,
	}

	for _, portfolio := range portfolios {
		apiFren, err := converters.ConvertFren(portfolio)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
			return
		}
		resp.Data = append(resp.Data, apiFren)
	}

	c.JSON(http.StatusOK, resp)
}
