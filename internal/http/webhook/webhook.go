package webhookhttp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/frenswifbenefits/myfren/internal/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (s *ServerImpl) PostV1Tv(c *gin.Context) {

	s.logger.Info("Start webhook from tv")
	if s.tvWhitelistEnabled {
		ip := c.ClientIP()
		if !slices.Contains(TvWhitelist, ip) {
			s.logger.Warn("skip Webhook from TV by Whitelist")
			c.JSON(http.StatusNotFound, dto.MakeNotFoundAPIResponse(errors.New("wrong ip")))
			return
		}
	}

	defer c.JSON(http.StatusOK, dto.MakeSuccessAPIResponse())
	signal := &dto.Signal{}
	bodyContent, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.logger.Error("failed to read signal body", zap.Error(err))
		return
	}
	err = json.Unmarshal(bodyContent, &signal)
	if err != nil {
		s.logger.Error("cannot parse signal body", zap.String("content", string(bodyContent)), zap.Error(err))
		return
	}

	s.logger.Info("Process webhook from tv parsed", zap.Any("signal", signal))

	signal.Exchange = strings.ToUpper(signal.Exchange)
	signal.Symbol = strings.ToUpper(signal.Symbol)

	s.si.ExecuteSignal(c, *signal)
}
