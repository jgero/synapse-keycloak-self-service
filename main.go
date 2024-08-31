package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

//go:embed form.html
var form string

func main() {
	e := echo.New()

	e.Use(AuthMiddleware)

	e.GET("/", func(c echo.Context) error {
		info := c.Get(UserInfoKey).(*UserInfo)
		template.Must(template.New("form").Parse(form)).Execute(c.Response().Writer, info)
		return nil
	})

	synapseHost := os.Getenv("SYNAPSE_HOST")
	synapseAccessToken := os.Getenv("SYNAPSE_ACCESS_TOKEN")
	synapseDomain := os.Getenv("SYNAPSE_DOMAIN")
	e.POST("/signup", func(c echo.Context) error {
		info := c.Get(UserInfoKey).(*UserInfo)
		err := c.Request().ParseForm()
		if err != nil {
			e.Logger.Error(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		username := c.Request().Form.Get("username")
		matrix := SynapseConnection{
			Host:        synapseHost,
			AccessToken: synapseAccessToken,
			Domain:      synapseDomain,
		}
		ok, err := matrix.IsUsernameAvailable(username)
		if err != nil {
			e.Logger.Error(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		if !ok {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Username %s is not available", username))
		}
		err = matrix.CreateAccount(info, username)
		if err != nil {
			e.Logger.Error(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.String(http.StatusOK, fmt.Sprintf("Account %s created successfully", username))
	})
	e.Logger.Fatal(e.Start(":8080"))
}
