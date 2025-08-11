package userhttp

import (
	"net/http"

	"github.com/frenswifbenefits/myfren/internal/dto"
	"github.com/frenswifbenefits/myfren/internal/http/converters"
	api_types "github.com/frenswifbenefits/myfren/internal/openapi"
	"github.com/gin-gonic/gin"
)

func (s *ServerImpl) GetV1UserFrenList(c *gin.Context) {
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
