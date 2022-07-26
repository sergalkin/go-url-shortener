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

// Cookie - uuid cookie middleware that attempts to read uid cookie and set it cookieUUID
// if no cookie was read, it creates a new one and stores it in cookieUUID.
// finally it adds uid cookie to http.ResponseWriter
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

// setCookie - attempts to read cookie from http.Request and generate a new one if could not.
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

// GetUUID - gets uuid from cookieUUID.
func GetUUID() string {
	defer cookie.mu.Unlock()
	cookie.mu.Lock()

	return cookie.uuid
}
