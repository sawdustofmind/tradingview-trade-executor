package adminhttp

import (
	"net/http"

	"github.com/frenswifbenefits/myfren/internal/dto"
	api_types "github.com/frenswifbenefits/myfren/internal/openapi"
	"github.com/gin-gonic/gin"
)

func (s *ServerImpl) GetV1AdminCustomerList(c *gin.Context) {
	customers := s.cp.GetAll()
	resp := api_types.CustomerListResponse{
		Code:    dto.SuccessCode,
		Message: dto.SuccessMessage,
	}
	for _, customer := range customers {
		resp.Data = append(resp.Data, api_types.Customer{
			Id:       customer.Id,
			Pnl:      "0",
			Username: customer.Username,
		})
	}
	c.JSON(http.StatusOK, resp)
}
