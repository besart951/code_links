package httpapi

import (
	"errors"
	"net/http"

	"github.com/besart951/code-links/apps/auth-service/backend/internal/userinfo"
	"github.com/google/uuid"
)

func (s *Server) registerUserInfoRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/userinfo/me", s.handleUserInfoMe)
	mux.HandleFunc("POST /api/userinfo/lookup", s.handleUserInfoLookup)
}

func (s *Server) handleUserInfoMe(w http.ResponseWriter, r *http.Request) {
	snapshot, err := s.users.CurrentUser(r.Context(), bearerToken(r))
	if err != nil {
		writeUserInfoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, snapshot)
}

func (s *Server) handleUserInfoLookup(w http.ResponseWriter, r *http.Request) {
	var request lookupUsersRequest
	if !decodeJSON(w, r, &request) {
		return
	}

	userIDs := make([]uuid.UUID, 0, len(request.UserIDs))
	for _, rawUserID := range request.UserIDs {
		userID, err := uuid.Parse(rawUserID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid user id")
			return
		}
		userIDs = append(userIDs, userID)
	}

	cards, err := s.users.LookupUsers(r.Context(), bearerToken(r), userIDs)
	if err != nil {
		writeUserInfoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cards)
}

func writeUserInfoError(w http.ResponseWriter, err error) {
	var serviceError *userinfo.Error
	if !errors.As(err, &serviceError) {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	status := http.StatusInternalServerError
	switch serviceError.Kind {
	case userinfo.KindBadRequest:
		status = http.StatusBadRequest
	case userinfo.KindUnauthorized:
		status = http.StatusUnauthorized
	}
	writeError(w, status, serviceError.Message)
}
