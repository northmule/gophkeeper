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
	manager       repository.Repository
	storage       storage.DBQuery
	session       storage.SessionManager
	log           *logger.Logger
	cfg           *config.Config
	accessService access.AccessService
	cryptService  service.CryptService
}

func NewAppRoutes(repositoryManager repository.Repository, storage storage.DBQuery, session storage.SessionManager, log *logger.Logger, cfg *config.Config, accessService access.AccessService, cryptService service.CryptService) AppRoutes {
	instance := AppRoutes{
		manager:       repositoryManager,
		storage:       storage,
		session:       session,
		log:           log,
		cfg:           cfg,
		accessService: accessService,
		cryptService:  cryptService,
	}
	return instance
}

// DefiningAppRoutes маршруты приложения
func (ar *AppRoutes) DefiningAppRoutes() chi.Router {

	// Обработчики
	healthHandler := NewHealthHandler(ar.log)
	registrationHandler := NewRegistrationHandler(ar.manager, ar.session, ar.accessService, ar.log)
	transactionHandler := NewTransactionHandler(ar.storage, ar.log)

	itemsListHandler := NewItemsListHandler(ar.accessService, ar.manager, ar.log)
	cardDataHandler := NewCardDataHandler(ar.accessService, ar.manager, ar.log)
	textDataHandler := NewTextDataHandler(ar.accessService, ar.manager, ar.log)
	fileDataHandler := NewFileDataHandler(ar.accessService, ar.manager, ar.cfg, ar.log)
	itemDataHandler := NewItemDataHandler(ar.accessService, ar.manager, ar.log)
	keysDataHandler := NewKeysDataHandler(ar.accessService, ar.cryptService, ar.manager, ar.cfg, ar.log)
	decryptDataHandler := NewDecryptDataHandler(ar.accessService, ar.manager, ar.log)

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
			jwtTokenObject := ar.accessService.FillJWTToken()

			// Проверка токена, заполнения данных о пользователе
			r.Use(jwtauth.Verify(jwtTokenObject, ar.accessService.FindTokenByRequest))
			r.Use(jwtauth.Authenticator(jwtTokenObject))

			// приём от клиента публичного ключа
			r.Post("/save_public_key", keysDataHandler.HandleSaveClientPublicKey)

			// приём от клиента приватного ключа(aes используется для шифрования данных)
			r.Post("/save_client_private_key", keysDataHandler.HandleSaveClientPrivateKey)

			// Клиент забирает публичный ключ сервера
			r.Post("/download_server_public_key", keysDataHandler.HandleDownloadServerPublicKey)

			// список сохранённых данных
			r.With(
				decryptDataHandler.HandleEncryptData, // шифрует исходящий запрос
			).Get("/items_list", itemsListHandler.HandleItemsList)

			// Получить данные по uuid
			r.With(
				decryptDataHandler.HandleEncryptData, // шифрует исходящий запрос
			).Get("/item_get/{uuid}", itemDataHandler.HandleItem)

			// добавить/изменить данные банковской карты
			r.With(
				decryptDataHandler.HandleDecryptData, // расшифровка тела запроса
				NewValidatorHandler(new(cardDataRequest), ar.log).HandleValidation,
			).Post("/save_card_data", cardDataHandler.HandleSave)

			// добавить/изменить текстовые данные
			r.With(
				decryptDataHandler.HandleDecryptData, // расшифровка тела запроса
				NewValidatorHandler(new(textDataRequest), ar.log).HandleValidation,
			).Post("/save_text_data", textDataHandler.HandleSave)

			// инициализация приёма файла, базовые данные о файле
			r.With(
				decryptDataHandler.HandleDecryptData, // расшифровка тела запроса
				NewValidatorHandler(new(fileDataInitRequest), ar.log).HandleValidation,
			).Post("/file_data/init", fileDataHandler.HandleInit)

			// приём данных файла
			r.With(
				decryptDataHandler.HandleDecryptData, // расшифровка тела запроса
			).Post("/file_data/load/{file_uuid}/{part}", fileDataHandler.HandleAction)

			// отдача файла клиенту
			r.With(
				decryptDataHandler.HandleEncryptData, // шифрует исходящий запрос
			).Post("/file_data/get/{file_uuid}/{part}", fileDataHandler.HandleGetAction)

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
