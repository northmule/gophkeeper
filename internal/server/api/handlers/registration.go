package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/api/rctx"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/northmule/gophkeeper/internal/server/repository"
	"github.com/northmule/gophkeeper/internal/server/services/access"
	"github.com/northmule/gophkeeper/internal/server/storage"
)

type RegistrationHandler struct {
	manager       repository.Repository
	session       storage.SessionManager
	accessService access.AccessService
	log           *logger.Logger
}

type registrationRequest struct {
	Login    string `json:"login" validate:"required,min=2,max=50"`
	Password string `json:"password" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
}

type authenticationRequest struct {
	Login    string `json:"login" validate:"required,min=2,max=50"`
	Password string `json:"password" validate:"required,min=3,max=100"`
}

func NewRegistrationHandler(manager repository.Repository, session storage.SessionManager, accessService access.AccessService, log *logger.Logger) *RegistrationHandler {
	instance := &RegistrationHandler{
		manager:       manager,
		session:       session,
		accessService: accessService,
		log:           log,
	}
	return instance
}

func (rr *registrationRequest) Bind(r *http.Request) error {
	return nil
}

func (ar *authenticationRequest) Bind(r *http.Request) error {
	return nil
}

// HandleRegistration регистрация пользователя
func (r *RegistrationHandler) HandleRegistration(res http.ResponseWriter, req *http.Request) {
	var err error
	request := &registrationRequest{}
	if err = render.Bind(req, request); err != nil {
		r.log.Info(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}

	user, err := r.manager.User().FindOneByLogin(req.Context(), request.Login)
	if err != nil {
		r.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
	if user != nil && user.Login == request.Login {
		r.log.Infof("The login '%s' is already occupied by the user with the uid '%s'", request.Login, user.UUID)
		_ = render.Render(res, req, ErrConflict(fmt.Errorf("the login '%s' is already occupied by the user", request.Login)))
		return
	}
	passwordHash, err := r.accessService.PasswordHash(request.Password)
	if err != nil {
		r.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
	var newUser models.User

	newUser.Login = request.Login
	newUser.Password = passwordHash
	newUser.Email = request.Email
	newUser.UUID = uuid.NewString()

	tx := req.Context().Value(rctx.TransactionCtxKey).(*storage.Transaction)
	userID, err := r.manager.User().TxCreateNewUser(req.Context(), tx, newUser)
	if err != nil {
		r.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}

	if userID == 0 {
		r.log.Error("An empty ID value when registering a user")
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
	newUser.ID = userID
	r.log.Infof("A new user '%s' has been registered with the uuid '%s'", newUser.Login, newUser.UUID)

	ctx := context.WithValue(req.Context(), rctx.UserCtxKey, newUser)
	req = req.WithContext(ctx)

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(nil)
}

// JWTClaims утверждение
type JWTClaims struct {
	jwt.RegisteredClaims
	UserUUID string `json:"user_uuid"`
}

// HandleAuthentication аунтификация пользователя
func (r *RegistrationHandler) HandleAuthentication(res http.ResponseWriter, req *http.Request) {

	var err error
	request := &authenticationRequest{}
	if err = render.Bind(req, request); err != nil {
		_ = render.Render(res, req, ErrBadRequest)
		return
	}

	user, err := r.manager.User().FindOneByLogin(req.Context(), request.Login)
	if err != nil {
		r.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}

	passwordHash, err := r.accessService.PasswordHash(request.Password)
	if err != nil {
		r.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}

	if user == nil || user.Password != passwordHash {
		r.log.Infof("Invalid username/password pair %s/****", request.Login)
		_ = render.Render(res, req, ErrUnauthorized)
		return
	}

	// Токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 300)), // todo в конфиг
		},
		UserUUID: user.UUID,
	})

	// подпись токена ключом
	tokenValue, err := token.SignedString([]byte("_secret_")) // todo в конфиг
	if err != nil {
		r.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}

	res.Header().Set("Authorization", "Bearer "+tokenValue)

	r.log.Infof("User %s has been authenticated and the %s token has been issued", user.UUID, tokenValue)
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(nil)
}
