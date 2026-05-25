package userinfo

import (
	"context"
	"errors"

	"github.com/besart951/code-links/apps/auth-service/backend/internal/domain"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/token"
	"github.com/google/uuid"
)

type Store interface {
	FindUserByID(ctx context.Context, userID uuid.UUID) (User, []string, error)
	ListUserRoles(ctx context.Context, userID uuid.UUID) ([]AdminRole, []AdminPermission, error)
	LookupUserCards(ctx context.Context, userIDs []uuid.UUID) ([]UserCard, error)
}

type TokenParser interface {
	Parse(rawToken string) (token.Claims, error)
}

type Service struct {
	store  Store
	tokens TokenParser
}

func NewService(store Store, tokens TokenParser) *Service {
	return &Service{store: store, tokens: tokens}
}

type Kind string

const (
	KindBadRequest   Kind = "bad_request"
	KindUnauthorized Kind = "unauthorized"
	KindInternal     Kind = "internal"
)

type Error struct {
	Kind    Kind
	Message string
}

func (e *Error) Error() string { return e.Message }

func serviceError(kind Kind, message string) *Error {
	return &Error{Kind: kind, Message: message}
}

func (s *Service) CurrentUser(ctx context.Context, bearerToken string) (UserSnapshot, error) {
	userID, err := s.userIDFromBearer(bearerToken)
	if err != nil {
		return UserSnapshot{}, serviceError(KindUnauthorized, "invalid access token")
	}
	return s.snapshot(ctx, userID)
}

func (s *Service) LookupUsers(ctx context.Context, bearerToken string, userIDs []uuid.UUID) ([]UserCard, error) {
	if _, err := s.userIDFromBearer(bearerToken); err != nil {
		return nil, serviceError(KindUnauthorized, "invalid access token")
	}
	if len(userIDs) > 100 {
		return nil, serviceError(KindBadRequest, "too many users requested")
	}

	cards, err := s.store.LookupUserCards(ctx, userIDs)
	if err != nil {
		return nil, serviceError(KindInternal, "could not load users")
	}
	return cards, nil
}

func (s *Service) snapshot(ctx context.Context, userID uuid.UUID) (UserSnapshot, error) {
	user, licenses, err := s.store.FindUserByID(ctx, userID)
	if errors.Is(err, domain.ErrNotFound) {
		return UserSnapshot{}, serviceError(KindUnauthorized, "invalid access token")
	}
	if err != nil {
		return UserSnapshot{}, serviceError(KindInternal, "could not load user")
	}

	roles, permissions, err := s.store.ListUserRoles(ctx, user.ID)
	if err != nil {
		return UserSnapshot{}, serviceError(KindInternal, "could not load user")
	}

	return UserSnapshot{
		ID:            user.ID.String(),
		Email:         user.Email,
		Name:          user.Name,
		Status:        user.Status,
		EmailVerified: user.EmailVerifiedAt != nil,
		Licenses:      licenses,
		Roles:         roles,
		Permissions:   permissions,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}, nil
}

func (s *Service) userIDFromBearer(bearerToken string) (uuid.UUID, error) {
	if bearerToken == "" {
		return uuid.Nil, errors.New("missing bearer token")
	}
	claims, err := s.tokens.Parse(bearerToken)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(claims.Subject)
}
