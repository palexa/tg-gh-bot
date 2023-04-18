package gh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ghActionTelegramBot/internal/adapters/api"
	"ghActionTelegramBot/internal/config"
	"ghActionTelegramBot/internal/domain/person"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"io"
	"log"
	"net/http"
)

type handler struct {
	app     *fiber.App
	service person.Service
}

func NewHandler(service person.Service) api.Handler {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	return &handler{app: app, service: service}
}

func (h *handler) Run() {
	h.Register()
	log.Fatal(h.app.Listen(":3000"))
}

func (h *handler) Register() {
	h.app.Get("/:userId", h.Root)
	h.app.Get("/login/github/:userId", h.GitHubLogin)
	h.app.Get("/login/github/callback/:userId", h.GitHubCallback)
}

func (h *handler) Root(c *fiber.Ctx) error {
	userId := c.Params("userId")
	return c.Render("index", fiber.Map{"Link": fmt.Sprintf("%s/login/github/%s", config.Cfg.BaseUrl, userId)})
}

func (h *handler) GitHubLogin(c *fiber.Ctx) error {
	userId := c.Params("userId")
	githubClientID := getGithubClientID()

	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		githubClientID,
		fmt.Sprintf("%s/login/github/callback/%s", config.Cfg.BaseUrl, userId),
	)
	err := c.Redirect(redirectURL, 301)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) GitHubCallback(c *fiber.Ctx) error {
	userId := c.Params("userId")
	fmt.Println(userId)

	code := c.Query("code")

	githubAccessToken := getGithubAccessToken(code)

	err := h.service.SetGHToken(&person.UpdatePersonDto{
		ID:          userId,
		AccessToken: githubAccessToken,
	})

	if err != nil {

	}

	githubData := getGithubData(githubAccessToken)

	return h.LoggedInHandler(c, githubData)
}

func (h *handler) LoggedInHandler(c *fiber.Ctx, githubData string) error {
	if githubData == "" {
		_, err := c.Write([]byte("UNAUTHORIZED!"))
		if err != nil {
			return err
		}
		return nil
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	var prettyJSON bytes.Buffer
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		log.Panic("JSON parse error")
	}

	return c.Send(prettyJSON.Bytes())
}

func getGithubClientID() string {

	githubClientID := config.Cfg.GitHub.ClientId

	return githubClientID
}

func getGithubClientSecret() string {

	githubClientSecret := config.Cfg.GitHub.ClientSecret

	return githubClientSecret
}

func getGithubAccessToken(code string) string {
	clientID := getGithubClientID()
	clientSecret := getGithubClientSecret()

	requestBodyMap := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := io.ReadAll(resp.Body)

	// Represents the response received from Github
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var ghresp githubAccessTokenResponse
	err := json.Unmarshal(respbody, &ghresp)
	if err != nil {

	}
	return ghresp.AccessToken
}

func getGithubData(accessToken string) string {
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := io.ReadAll(resp.Body)

	return string(respbody)
}
