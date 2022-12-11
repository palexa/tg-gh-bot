package gh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ghActionTelegramBot/internal/adapters/api"
	"ghActionTelegramBot/internal/config"
	"ghActionTelegramBot/internal/domain/person"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"net/http"
)

type handler struct {
	app     *fiber.App
	service person.Service
}

func NewHandler(service person.Service) api.Handler {
	app := fiber.New()
	return &handler{app: app, service: service}
}

func (h *handler) Run() {
	h.Register()
	log.Fatal(h.app.Listen(":3000"))
}

func (h *handler) Register() {
	h.app.Get("/", h.Root)
	h.app.Get("/login/github/", h.GitHubLogin)
	h.app.Get("/login/github/callback", h.GitHubCallback)
}

func (h *handler) Root(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	err := c.SendString(`<a href="/login/github/">LOGIN</a>`)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) GitHubLogin(c *fiber.Ctx) error {
	githubClientID := getGithubClientID()

	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		githubClientID,
		"http://localhost:3000/login/github/callback",
	)
	err := c.Redirect(redirectURL, 301)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) GitHubCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	githubAccessToken := getGithubAccessToken(code)

	githubData := getGithubData(githubAccessToken)

	return h.LoggedInHandler(c, githubData)
}

func (h *handler) LoggedInHandler(c *fiber.Ctx, githubData string) error {
	if githubData == "" {
		// Unauthorized users get an unauthorized message
		_, err := c.Write([]byte("UNAUTHORIZED!"))
		if err != nil {
			return err
		}
		//fmt.Fprintf(w, "UNAUTHORIZED!")
		return nil
	}

	// Set return type JSON
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Prettifying the json
	var prettyJSON bytes.Buffer
	// json.indent is a library utility function to prettify JSON indentation
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		log.Panic("JSON parse error")
	}

	// Return the prettified JSON as a string
	return c.Send(prettyJSON.Bytes())
	//fmt.Fprintf(w, string(prettyJSON.Bytes()))
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
	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken
}

func getGithubData(accessToken string) string {
	// Get request to a set URL
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

	// Make the request
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	// Read the response as a byte slice
	respbody, _ := io.ReadAll(resp.Body)

	// Convert byte slice to string and return
	return string(respbody)
}
