package handlers

import (
	"io"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/render"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/keys"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/config"
	"github.com/northmule/gophkeeper/internal/server/logger"
	service "github.com/northmule/gophkeeper/internal/server/services"
	"golang.org/x/net/context"
)

// Ожидаемая схема взаимодействия:
// 1. Клиент авторизуется
// 2. Отправляет серверу свой rsa публичный ключ
// 3. Клиент запрашивает публичный rsa ключ сервера (ключ общий для всех на сервере)
// 4. Клиент при помощи публичного rsa ключа сервера шифрует свой секретный aes ключ и направляет его серверу
// 5. Сервер дешефрует секретный ключ своим rsa приватным ключом
// 6. Дальнейший обмен данных шифрование и дешифрование проивзодится приватным ключом клиента

// KeysDataHandler обработка запросо с ключами
type KeysDataHandler struct {
	log             *logger.Logger
	userFinderByJWT UserFinderByJWT
	keySaver        KeySaver
	userFinder      UserFinder
	expectedAction  map[string]bool

	cfg            *config.Config
	publicKeyPath  string
	privateKeyPath string

	cryptService service.CryptService
}

// NewKeysDataHandler конструктор
func NewKeysDataHandler(userFinderByJWT UserFinderByJWT, cryptService service.CryptService, keySaver KeySaver, userFinder UserFinder, cfg *config.Config, log *logger.Logger) *KeysDataHandler {

	return &KeysDataHandler{
		userFinderByJWT: userFinderByJWT,
		keySaver:        keySaver,
		userFinder:      userFinder,
		log:             log,
		cfg:             cfg,
		publicKeyPath:   path.Join(cfg.Value().PathKeys, keys.PublicKeyFileName),
		privateKeyPath:  path.Join(cfg.Value().PathKeys, keys.PrivateKeyFileName),
		cryptService:    cryptService,
	}
}

// KeySaver сохранение ключей
type KeySaver interface {
	SetPublicKey(ctx context.Context, data string, userUUID string) error
	SetPrivateClientKey(ctx context.Context, data string, userUUID string) error
}

// HandleSaveClientPublicKey привязка публичного ключа клиента
func (h *KeysDataHandler) HandleSaveClientPublicKey(res http.ResponseWriter, req *http.Request) {
	var (
		err      error
		userUUID string
	)

	userUUID, err = h.userFinderByJWT.GetUserUUIDByJWTToken(req.Context())
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}

	err = req.ParseMultipartForm(4096)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
	keyData, _, err := req.FormFile(data_type.FileField)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
	defer keyData.Close()

	keyBytes, err := io.ReadAll(keyData)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}

	keyString := string(keyBytes)

	err = h.keySaver.SetPublicKey(req.Context(), keyString, userUUID)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
}

// HandleDownloadServerPublicKey возвращает клиенту публичный ключ сервера
func (h *KeysDataHandler) HandleDownloadServerPublicKey(res http.ResponseWriter, req *http.Request) {
	var (
		err      error
		userUUID string
		user     *models.User
	)
	userUUID, err = h.userFinderByJWT.GetUserUUIDByJWTToken(req.Context())
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}

	user, err = h.userFinder.FindOneByUUID(req.Context(), userUUID)
	if err != nil || user == nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}

	buff, err := os.ReadFile(h.publicKeyPath)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}

	_, err = res.Write(buff)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
}

// HandleSaveClientPrivateKey привязка приватного ключа клиента. Ключ приходит зашифрованный публичным ключом сервера
func (h *KeysDataHandler) HandleSaveClientPrivateKey(res http.ResponseWriter, req *http.Request) {
	var (
		err      error
		userUUID string
	)

	userUUID, err = h.userFinderByJWT.GetUserUUIDByJWTToken(req.Context())
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}

	err = req.ParseMultipartForm(4096)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
	keyData, _, err := req.FormFile(data_type.FileField)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
	defer keyData.Close()

	keyBytes, err := io.ReadAll(keyData)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
	// Расшифровываем ключ
	keyBytes, err = h.cryptService.DecryptRSA(keyBytes)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
	keyString := string(keyBytes)
	// Секретный ключ клиента сохраняется
	err = h.keySaver.SetPrivateClientKey(req.Context(), keyString, userUUID)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
}
