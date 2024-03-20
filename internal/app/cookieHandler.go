package app

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/securecookie"
)

// cookieHandler - структура, заведующая куками.
type cookieHandler struct {
	instance *securecookie.SecureCookie
}

const (
	cookieName = "user-info"
	valName    = "email"
)

// getCookieHandler - функция, возвращающая cookieHandler.
func getCookieHandler() cookieHandler {
	hashKey, blockKey := make([]byte, 32), make([]byte, 16)
	rand.Read(hashKey)
	rand.Read(blockKey)

	var s = securecookie.New(hashKey, blockKey)
	return cookieHandler{s}
}

// Set - функция, устанавливающая в http.Response куки, содержащий email пользователя.
func (c *cookieHandler) Set(ctx *fiber.Ctx, email string) {
	value := map[string]string{
		valName: email,
	}
	if encoded, err := c.instance.Encode(cookieName, value); err == nil {
		now := time.Now()
		cookie := &fiber.Cookie{
			Name:    cookieName,
			Value:   encoded,
			Path:    "/",
			Secure:  true,
			Expires: now.Add(7 * 24 * time.Hour),
		}
		ctx.Cookie(cookie) // TODO
	}
}

// Read - функция, получающая из куки http.Request email пользователя.
func (c *cookieHandler) Read(ctx *fiber.Ctx) (string, error) {
	cookie := string(ctx.Request().Header.Cookie(cookieName))

	if cookie != "" {
		value := make(map[string]string)
		if err := c.instance.Decode(cookieName, cookie, &value); err == nil {
			return value[valName], nil
		}
	}

	return "", errors.New("cannot read cookie")
}
