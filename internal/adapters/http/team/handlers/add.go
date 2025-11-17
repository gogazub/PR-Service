package teamhandlers

import (
	httperror "PRService/internal/adapters/http/error"
	teamhttp "PRService/internal/adapters/http/team"
	"PRService/internal/domain/team"
	team_usecase "PRService/internal/usecase/team"
	"encoding/json"
	"errors"
	"net/http"
)

type AddTeamRequest = teamhttp.TeamDTO

type AddTeamResponseDTO struct {
	Team teamhttp.TeamDTO `json:"team"`
}

func (h *Handler) AddTeam(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req AddTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteErrorResponse(w, http.StatusBadRequest, httperror.ErrorCodeBadRequest, "invalid json")
		return
	}

	cmd := team_usecase.CreateTeamAndUsersCommand{
		Name:    req.TeamName,
		Members: UsersFromMembers(req.Members, req.TeamName),
	}

	ctx := r.Context()
	t, usrs, err := h.CreateTeam(ctx, cmd)
	if err != nil {
		if errors.Is(err, team.ErrTeamExists) {
			httperror.WriteErrorResponse(w, http.StatusBadRequest, httperror.ErrorCodeTeamExists, "team is already exists")
			return
		}
		h.logger.Error("create team failed", "error", err)
		httperror.WriteErrorResponse(w, http.StatusInternalServerError, httperror.ErrorCodeInternal, "internal error")
		return
	}
	var resp AddTeamResponseDTO
	resp.Team = teamhttp.TeamToDTO(t, usrs)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		h.logger.Error("team add failed: JSON encoding", "error", err)
		httperror.WriteErrorResponse(w, http.StatusInternalServerError, httperror.ErrorCodeInternal, "internal error")
		return
	}
}
