package main

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	UserInfoKey = "userinfo"
	cookieName  = "access_token"
)

// does auth with keycloak and sets the userinfo in the context
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	keycloak := KeycloakConnector{
		Host:     os.Getenv("KEYCLOAK_HOST"),
		Realm:    os.Getenv("KEYCLOAK_REALM"),
		ClientId: os.Getenv("KEYCLOAK_CLIENT_ID"),
	}
	return func(c echo.Context) error {
		if cookie, err := c.Request().Cookie(cookieName); err != nil {
			if err != http.ErrNoCookie {
				c.Logger().Error(err)
				return err
			}
		} else if cookie != nil {
			info, err := keycloak.GetUserinfo(cookie.Value)
			if err != nil {
				c.Logger().Error(err)
				return err
			}
			c.Set(UserInfoKey, &info)
			return next(c)
		}
		// no cookie, do authentication steps
		authCode := c.Request().URL.Query().Get("code")
		if authCode == "" {
			return c.Redirect(http.StatusFound, keycloak.GetAuthUrl())
		} else {
			token, err := keycloak.GetAccesToken(authCode)
			if err != nil {
				c.Logger().Error(err)
				return err
			}
			c.SetCookie(&http.Cookie{
				Name:    cookieName,
				Value:   token.AccessToken,
				Expires: time.Now().Add(time.Duration(int(token.ExpiresIn)) * time.Second),
			})
			return c.Redirect(http.StatusFound, "/")
		}
	}
}
