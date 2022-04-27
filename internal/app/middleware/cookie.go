package middleware

import (
	"net/http"
	"sync"

	"github.com/google/uuid"

	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
)

type cookieUUID struct {
	name string
	uuid string
	mu   sync.Mutex
}

var cookie cookieUUID

func init() {
	cookie = *New()
}

func New() *cookieUUID {
	return &cookieUUID{name: "uid"}
}

func Cookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		sha, err := setCookie(writer, request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		cookie.uuid = sha

		http.SetCookie(writer, &http.Cookie{
			Name:   cookie.name,
			Value:  cookie.uuid,
			Path:   "/",
			Secure: false,
			MaxAge: 300000,
		})

		next.ServeHTTP(writer, request)
	})
}

func setCookie(writer http.ResponseWriter, request *http.Request) (string, error) {
	defer cookie.mu.Unlock()
	cookie.mu.Lock()

	cookie.uuid = uuid.New().String()

	if cookieUserID, err := request.Cookie(cookie.name); err == nil {
		utils.Decode(cookieUserID.Value, &cookie.uuid)
	}

	sha, err := utils.Encode(cookie.uuid)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	return sha, nil
}

func GetUUID() string {
	defer cookie.mu.Unlock()
	cookie.mu.Lock()

	return cookie.uuid
}
