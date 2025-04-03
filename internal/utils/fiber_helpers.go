package utils

import (
	"log"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/starks97/alcohol-tracker-api/internal/state"
)

type FiberHelpers interface {
	SetCookie(ctx *fiber.Ctx, cookieName string, cookieValue string, expiration time.Duration) error
}

type FiberHelper struct {
	appState *state.AppState
}

func NewFiberHelper(appState *state.AppState) FiberHelpers {
	return &FiberHelper{
		appState: appState,
	}
}

func (fc *FiberHelper) SetCookie(ctx *fiber.Ctx, cookieName string, cookieValue string, expiration time.Duration) error {
	clientOriginURL, err := url.Parse(fc.appState.Config.ClientOrigin)
	if err != nil {
		log.Fatalf("Invalid CLIENT_ORIGIN: %v", err)
	}
	cookieDomain := clientOriginURL.Hostname()

	cookie := &fiber.Cookie{
		Name:     cookieName,
		Value:    cookieValue,
		Domain:   cookieDomain,
		Path:     "/",
		Secure:   true,
		HTTPOnly: true,
		Expires:  time.Now().Add(expiration),
	}
	ctx.Cookie(cookie)

	return nil

}
