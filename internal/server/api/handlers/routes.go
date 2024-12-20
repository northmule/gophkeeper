package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/northmule/gophkeeper/internal/server/config"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/northmule/gophkeeper/internal/server/repository"
	service "github.com/northmule/gophkeeper/internal/server/services"
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
	cryptService, err := service.NewCrypt(ar.cfg)
	if err != nil {
		ar.log.Error(err)
	}

	// Обработчики
	healthHandler := NewHealthHandler(ar.log)
	registrationHandler := NewRegistrationHandler(ar.manager, ar.session, accessService, ar.log)
	transactionHandler := NewTransactionHandler(ar.storage, ar.log)

	itemsListHandler := NewItemsListHandler(accessService, ar.manager, ar.log)
	cardDataHandler := NewCardDataHandler(accessService, ar.manager, ar.log)
	textDataHandler := NewTextDataHandler(accessService, ar.manager, ar.log)
	fileDataHandler := NewFileDataHandler(accessService, ar.manager, ar.cfg, ar.log)
	itemDataHandler := NewItemDataHandler(accessService, ar.manager, ar.log)
	keysDataHandler := NewKeysDataHandler(accessService, cryptService, ar.manager, ar.cfg, ar.log)
	decryptDataHandler := NewDecryptDataHandler(accessService, ar.manager, ar.log)

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
			// Начальный jwt объект
			jwtTokenObject := accessService.FillJWTToken()

			// приём от клиента публичного ключа
			r.With(
				jwtauth.Verify(jwtTokenObject, accessService.FindTokenByRequest),
				jwtauth.Authenticator(jwtTokenObject),
			).Post("/save_public_key", keysDataHandler.HandleSaveClientPublicKey)

			// приём от клиента приватного ключа(aes используется для шифрования данных)
			r.With(
				jwtauth.Verify(jwtTokenObject, accessService.FindTokenByRequest),
				jwtauth.Authenticator(jwtTokenObject),
			).Post("/save_client_private_key", keysDataHandler.HandleSaveClientPrivateKey)

			// Клиент забирает публичный ключ сервера
			r.With(
				jwtauth.Verify(jwtTokenObject, accessService.FindTokenByRequest),
				jwtauth.Authenticator(jwtTokenObject),
			).Post("/download_server_public_key", keysDataHandler.HandleDownloadServerPublicKey)

			// список сохранённых данных
			r.With(
				jwtauth.Verify(jwtTokenObject, accessService.FindTokenByRequest),
				jwtauth.Authenticator(jwtTokenObject),
			).Get("/items_list", itemsListHandler.HandleItemsList)

			// Получить данные по uuid
			r.With(
				jwtauth.Verify(jwtTokenObject, accessService.FindTokenByRequest),
				jwtauth.Authenticator(jwtTokenObject),
			).Get("/item_get/{uuid}", itemDataHandler.HandleItem)

			// добавить/изменить данные банковской карты
			r.With(
				jwtauth.Verify(jwtTokenObject, accessService.FindTokenByRequest),
				jwtauth.Authenticator(jwtTokenObject),
				decryptDataHandler.HandleDecryptData, // расшифровка тела запроса
				NewValidatorHandler(new(cardDataRequest), ar.log).HandleValidation,
			).Post("/save_card_data", cardDataHandler.HandleSave)

			// добавить/изменить текстовые данные
			r.With(
				jwtauth.Verify(jwtTokenObject, accessService.FindTokenByRequest),
				jwtauth.Authenticator(jwtTokenObject),
				decryptDataHandler.HandleDecryptData, // расшифровка тела запроса
				NewValidatorHandler(new(textDataRequest), ar.log).HandleValidation,
			).Post("/save_text_data", textDataHandler.HandleSave)

			// инициализация приёма файла, базовые данные о файле
			r.With(
				jwtauth.Verify(jwtTokenObject, accessService.FindTokenByRequest),
				jwtauth.Authenticator(jwtTokenObject),
				decryptDataHandler.HandleDecryptData, // расшифровка тела запроса
				NewValidatorHandler(new(fileDataInitRequest), ar.log).HandleValidation,
			).Post("/file_data/init", fileDataHandler.HandleInit)

			// приём данных файла
			r.With(
				jwtauth.Verify(jwtTokenObject, accessService.FindTokenByRequest),
				jwtauth.Authenticator(jwtTokenObject),
			).Post("/file_data/{action}/{file_uuid}/{part}", fileDataHandler.HandleAction)

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
