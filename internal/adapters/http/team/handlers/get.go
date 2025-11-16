package teamhandlers

import (
	httperror "PRService/internal/adapters/http/error"
	teamhttp "PRService/internal/adapters/http/team"
	"PRService/internal/domain/team"
	"encoding/json"
	"errors"
	"net/http"
)

type GetTeamRequest struct {
	TeamName string `query:"team_name"`
}

type GetTeamResponse struct {
	Team teamhttp.TeamDTO `json:"team"`
}

func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodGet {
		httperror.WriteBadRequest(w, "bad method")
		return
	}

	// 1. get team
	teamName := r.URL.Query().Get("team_name")
	ctx := r.Context()

	t, err := h.Services.Team.GetByName(ctx, teamName)
	if err != nil {
		if errors.Is(err, team.ErrTeamNotFound) {
			httperror.WriteErrorResponse(w, http.StatusBadRequest, httperror.ErrorCodeNotFound, "not found")
			return
		}
		h.logger.Errorf("get team: %w", err)
		httperror.WriteErrorResponse(w, http.StatusInternalServerError, httperror.ErrorCodeInternal, "internal error")
		return
	}

	// 2. get users
	users, err := h.Services.User.GetByIDs(ctx, t.Members)
	if err != nil {
		h.logger.Errorf("get team: %w", err)
		httperror.WriteErrorResponse(w, http.StatusInternalServerError, httperror.ErrorCodeInternal, "internal error")
		return
	}

	// 3. make response
	members := MembersFromUsers(users)
	resp := GetTeamResponse{
		Team: teamhttp.TeamDTO{
			TeamName: teamName,
			Members:  members,
		},
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorf("get team: json encode: %w", err)
	}

}
