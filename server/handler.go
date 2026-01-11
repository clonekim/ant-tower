package server

import (
	"net/http"
	"strings"
	"twn-monitor/sysagent"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	Hub *WsHub
}

func NewHandler(hub *WsHub) *Handler {
	return &Handler{Hub: hub}
}

func (h *Handler) GetUptime(c *gin.Context) {
	uptime, err := sysagent.GetUptime()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get system uptime")
		respondError(c, http.StatusInternalServerError, "Failed to get system uptime")
		return
	}

	respondSuccess(c, "Uptime retrieved", gin.H{
		"uptime_seconds": uptime,
	})
}

func (h *Handler) KillProcess(c *gin.Context) {
	var req struct {
		PID int32 `json:"pid"`
	}

	if !bindJSONOrError(c, &req) {
		return
	}

	if err := sysagent.KillProcess(req.PID); err != nil {
		log.Error().Err(err).Int32("pid", req.PID).Msg("Failed to kill process")

		// Handle Permission Error
		if strings.Contains(err.Error(), "Access") || strings.Contains(err.Error(), "denied") {
			respondError(c, http.StatusForbidden, "Permission denied. You can only terminate your own processes.")
		} else {
			respondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	log.Info().Int32("pid", req.PID).Msg("Process killed successfully")

	respondSuccess(c, "Process terminated", gin.H{
		"pid": req.PID,
	})
}

func (h *Handler) ControlPower(c *gin.Context) {
	var req struct {
		Action string `json:"action"`
	}

	if !bindJSONOrError(c, &req) {
		return
	}

	log.Warn().Str("action", req.Action).Msg("System power control requested")

	if err := sysagent.ControlPower(req.Action); err != nil {
		log.Error().Err(err).Str("action", req.Action).Msg("Failed to execute power control")
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(c, "Power command executed", gin.H{
		"action": req.Action,
	})
}

func (h *Handler) GetProcessList(c *gin.Context) {
	list, err := sysagent.GetProcessSnapshots()
	if err != nil {
		log.Error().Err(err).Msg("Failed to list processes")
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(c, "Process list retrieved", gin.H{
		"count":     len(list),
		"processes": list,
	})
}

func (h *Handler) GetCurrentUser(c *gin.Context) {

	respondSuccess(c, "Logged-in user retrieved", gin.H{
		"username": sysagent.GetCurrentUser(),
	})
}
