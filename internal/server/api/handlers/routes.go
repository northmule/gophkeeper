package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/northmule/gophkeeper/internal/server/config"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/northmule/gophkeeper/internal/server/repository"
	service "github.com/northmule/gophkeeper/internal/server/services"
	"github.com/northmule/gophkeeper/internal/server/storage"
	"golang.org/x/net/context"
)

type AppRoutes struct {
	storage       storage.DBQuery
	session       storage.SessionManager
	log           *logger.Logger
	cfg           *config.Config
	accessService AccessService
	cryptService  service.CryptService

	userRepository     *repository.UserRepository
	cardDataRepository *repository.CardDataRepository
	fileDataRepository *repository.FileDataRepository
	metaDataRepository *repository.MetaDataRepository
	ownerRepository    *repository.OwnerRepository
	textDataRepository *repository.TextDataRepository
}

func NewAppRoutes(storage storage.DBQuery, session storage.SessionManager, log *logger.Logger, cfg *config.Config, accessService AccessService, cryptService service.CryptService) *AppRoutes {
	instance := AppRoutes{
		storage:       storage,
		session:       session,
		log:           log,
		cfg:           cfg,
		accessService: accessService,
		cryptService:  cryptService,
	}
	return &instance
}

// AccessService серис доступа
type AccessService interface {
	PasswordHash(password string) (string, error)
	FillJWTToken() *jwtauth.JWTAuth
	GetUserUUIDByJWTToken(ctx context.Context) (string, error)
	FindTokenByRequest(r *http.Request) string
}

// DefiningAppRoutes маршруты приложения
func (ar *AppRoutes) DefiningAppRoutes() chi.Router {

	// Обработчики
	healthHandler := NewHealthHandler(ar.log)
	registrationHandler := NewRegistrationHandler(ar.userRepository, ar.userRepository, ar.session, ar.accessService, ar.log)
	transactionHandler := NewTransactionHandler(ar.storage, ar.log)

	itemsListHandler := NewItemsListHandler(ar.accessService, ar.ownerRepository, ar.log)
	cardDataHandler := NewCardDataHandler(ar.accessService, ar.ownerRepository, ar.cardDataRepository, ar.metaDataRepository, ar.log)
	textDataHandler := NewTextDataHandler(ar.accessService, ar.ownerRepository, ar.metaDataRepository, ar.textDataRepository, ar.log)
	fileDataHandler := NewFileDataHandler(ar.accessService, ar.fileDataRepository, ar.ownerRepository, ar.metaDataRepository, ar.cfg, ar.log)
	itemDataHandler := NewItemDataHandler(ar.accessService, ar.cardDataRepository, ar.metaDataRepository, ar.fileDataRepository, ar.textDataRepository, ar.ownerRepository, ar.log)
	keysDataHandler := NewKeysDataHandler(ar.accessService, ar.cryptService, ar.userRepository, ar.userRepository, ar.cfg, ar.log)
	decryptDataHandler := NewDecryptDataHandler(ar.accessService, ar.userRepository, ar.log)

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

// SetTextDataRepository установка репозитария
func (ar *AppRoutes) SetTextDataRepository(textDataRepository *repository.TextDataRepository) *AppRoutes {
	ar.textDataRepository = textDataRepository
	return ar
}

// SetOwnerRepository установка репозитария
func (ar *AppRoutes) SetOwnerRepository(ownerRepository *repository.OwnerRepository) *AppRoutes {
	ar.ownerRepository = ownerRepository
	return ar
}

// SetMetaDataRepository установка репозитария
func (ar *AppRoutes) SetMetaDataRepository(metaDataRepository *repository.MetaDataRepository) *AppRoutes {
	ar.metaDataRepository = metaDataRepository
	return ar
}

// SetFileDataRepository установка репозитария
func (ar *AppRoutes) SetFileDataRepository(fileDataRepository *repository.FileDataRepository) *AppRoutes {
	ar.fileDataRepository = fileDataRepository
	return ar
}

// SetCardDataRepository установка репозитария
func (ar *AppRoutes) SetCardDataRepository(cardDataRepository *repository.CardDataRepository) *AppRoutes {
	ar.cardDataRepository = cardDataRepository
	return ar
}

// SetUserRepository установка репозитария
func (ar *AppRoutes) SetUserRepository(userRepository *repository.UserRepository) *AppRoutes {
	ar.userRepository = userRepository
	return ar
}
