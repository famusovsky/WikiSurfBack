package app

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/securecookie"
)

// cookieHandler - структура, хранящая данные и обрабатывающая Cookie.
type cookieHandler struct {
	instance  *securecookie.SecureCookie
	name, val string
}

// getCookieHandler - функция, возвращающая cookieHandler.
func getCookieHandler(cookie, val string) cookieHandler {
	hashKey, blockKey := make([]byte, 32), make([]byte, 16)
	rand.Read(hashKey)
	rand.Read(blockKey)

	var s = securecookie.New(hashKey, blockKey)
	return cookieHandler{s, cookie, val}
}

// Set - функция, устанавливающая в http.Response куки с данным именем и значением.
func (c *cookieHandler) Set(ctx *fiber.Ctx, email string) {
	value := map[string]string{
		c.val: email,
	}
	if encoded, err := c.instance.Encode(c.name, value); err == nil {
		now := time.Now()
		cookie := &fiber.Cookie{
			Name:    c.name,
			Value:   encoded,
			Path:    "/",
			Secure:  true,
			Expires: now.Add(7 * 24 * time.Hour),
		}
		ctx.Cookie(cookie) // Probably fixed
	}
}

// Read - функция, получающая из http.Request данное значение куки с данным именем.
func (c *cookieHandler) Read(ctx *fiber.Ctx) (string, error) {
	cookie := string(ctx.Request().Header.Cookie(c.name))

	if cookie != "" {
		value := make(map[string]string)
		if err := c.instance.Decode(c.name, cookie, &value); err == nil {
			return value[c.val], nil
		}
	}

	return "", errors.New("cannot read cookie")
}

// Remove - функция, удаляющая куки с данным именем.
func (c *cookieHandler) Remove(ctx *fiber.Ctx) {
	ctx.ClearCookie(c.name)
}
