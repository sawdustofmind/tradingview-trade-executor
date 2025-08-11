package adminhttp

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/frenswifbenefits/myfren/internal/dto"
	api_types "github.com/frenswifbenefits/myfren/internal/openapi"
	"github.com/gin-gonic/gin"
)

func (s *ServerImpl) PostV1AdminGenerateInviteToken(c *gin.Context) {
	req := &api_types.PostV1AdminGenerateInviteTokenJSONBody{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(err))
		return
	}

	if req.Count < 1 || req.Count > 20 {
		c.JSON(http.StatusBadRequest, dto.MakeBadRequestAPIResponse(errors.New("count must be between 1 and 20")))
		return
	}

	resp := api_types.InviteTokensResponse{
		Code:    dto.SuccessCode,
		Message: dto.SuccessMessage,
	}
	for i := 0; i < req.Count; i++ {
		bts := make([]byte, 32)
		_, err = rand.Read(bts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
			return
		}

		token := hex.EncodeToString(bts)

		err = s.repository.InsertInviteToken(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.MakeErrorAPIResponse(err))
			return
		}
		resp.Data = append(resp.Data, token)
	}

	c.JSON(http.StatusOK, resp)
}
