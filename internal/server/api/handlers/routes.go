package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/northmule/gophkeeper/internal/server/config"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/northmule/gophkeeper/internal/server/repository"
	"github.com/northmule/gophkeeper/internal/server/services/access"
	"github.com/northmule/gophkeeper/internal/server/storage"
)

type AppRoutes struct {
	manager repository.Repository
	storage storage.DBQuery
	session storage.SessionManager
	log     *logger.Logger
	cfg     *config.Config
}

func NewAppRoutes(repositoryManager repository.Repository, storage storage.DBQuery, session storage.SessionManager, log *logger.Logger, cfg *config.Config) AppRoutes {
	instance := AppRoutes{
		manager: repositoryManager,
		storage: storage,
		session: session,
		log:     log,
		cfg:     cfg,
	}
	return instance
}

// DefiningAppRoutes маршруты приложения
func (ar *AppRoutes) DefiningAppRoutes() chi.Router {

	//Сервисы
	accessService := access.NewAccess(ar.cfg)

	// Обработчики
	healthHandler := NewHealthHandler(ar.log)
	registrationHandler := NewRegistrationHandler(ar.manager, ar.session, accessService, ar.log)
	transactionHandler := NewTransactionHandler(ar.storage, ar.log)

	itemsListHandler := NewItemsListHandler(accessService, ar.manager, ar.log)
	cardDataHandler := NewCardDataHandler(accessService, ar.manager, ar.log)
	textDataHandler := NewTextDataHandler(accessService, ar.manager, ar.log)

	r := chi.NewRouter()

	// Общие мидлвары
	r.Use(middleware.RequestLogger(ar.log))
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Compress(ar.cfg.Value().HTTPCompressLevel))

	r.Route("/api/v1", func(r chi.Router) {
		// Для авторизованных api
		r.Group(func(r chi.Router) {
			jwtTokenObject := accessService.FillJWTToken()
			// список сохранённых данных
			r.With(
				jwtauth.Verify(jwtTokenObject, accessService.FindTokenByRequest),
				jwtauth.Authenticator(jwtTokenObject),
			).Get("/items_list", itemsListHandler.HandleItemsList)

			// добавить/изменить данные банковской карты
			r.With(
				jwtauth.Verify(jwtTokenObject, accessService.FindTokenByRequest),
				jwtauth.Authenticator(jwtTokenObject),
				NewValidatorHandler(new(cardDataRequest), ar.log).HandleValidation,
			).Post("/save_card_data", cardDataHandler.HandleSave)

			// добавить/изменить текстовые данные
			r.With(
				jwtauth.Verify(jwtTokenObject, accessService.FindTokenByRequest),
				jwtauth.Authenticator(jwtTokenObject),
				NewValidatorHandler(new(textDataRequest), ar.log).HandleValidation,
			).Post("/save_text_data", textDataHandler.HandleSave)
		})

		// Общедоступное api
		r.Group(func(r chi.Router) {
			// состояние сервера
			r.Get("/health", healthHandler.HandleGetHealth)

			// регистрация пользователя
			r.With(
				NewValidatorHandler(new(registrationRequest), ar.log).HandleValidation,
				transactionHandler.Transaction,
			).Post("/register", registrationHandler.HandleRegistration)

			// аутентификация пользователя
			r.With(
				NewValidatorHandler(new(authenticationRequest), ar.log).HandleValidation,
			).Post("/login", registrationHandler.HandleAuthentication)
		})

	})

	return r
}
