package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"acm-portal/internal/config"
	"acm-portal/internal/loginstate"
	"acm-portal/internal/ory"

	"github.com/gofiber/fiber/v2"
)

type hydraLoginRequest struct {
	Skip    bool   `json:"skip"`
	Subject string `json:"subject"`
}

type hydraLoginAcceptResponse struct {
	RedirectTo string `json:"redirect_to"`
}

type kratosWhoamiResponse struct {
	Identity struct {
		ID string `json:"id"`
	} `json:"identity"`
}

func LoginHandler(c *fiber.Ctx) error {
	challenge := strings.TrimSpace(c.Query("login_challenge"))
	if challenge == "" {
		token := strings.TrimSpace(c.Query("lcid"))
		if token == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing login_challenge"})
		}
		resolved, ok := loginstate.ConsumeChallenge(token)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid or expired login challenge"})
		}
		challenge = resolved
	}

	cfg := config.Get()
	if cfg.HydraAdminURL == "" || cfg.KratosPublicURL == "" || cfg.KratosUIURL == "" || cfg.AuthPublicURL == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "missing auth service configuration"})
	}

	kratosPublicURL := strings.TrimSuffix(cfg.KratosPublicURL, "/")

	loginReqURL := ory.HydraAdminURL(cfg.HydraAdminURL, "/oauth2/auth/requests/login") + "?login_challenge=" + url.QueryEscape(challenge)
	loginResp, err := http.Get(loginReqURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch login request"})
	}
	defer loginResp.Body.Close()
	if loginResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(loginResp.Body)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid login challenge",
			"details": string(body),
		})
	}

	var loginReq hydraLoginRequest
	if err := json.NewDecoder(loginResp.Body).Decode(&loginReq); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to decode login request"})
	}

	if loginReq.Skip && strings.TrimSpace(loginReq.Subject) != "" {
		return acceptLogin(c, cfg.HydraAdminURL, challenge, loginReq.Subject)
	}

	whoamiReq, err := http.NewRequest(http.MethodGet, kratosPublicURL+"/sessions/whoami", nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to build kratos session request"})
	}
	if cookie := c.Get("Cookie"); cookie != "" {
		whoamiReq.Header.Set("Cookie", cookie)
	}

	whoamiResp, err := http.DefaultClient.Do(whoamiReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch kratos session"})
	}
	defer whoamiResp.Body.Close()

	if whoamiResp.StatusCode == http.StatusUnauthorized || whoamiResp.StatusCode == http.StatusForbidden {
		token, err := loginstate.SaveChallenge(challenge)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create login session"})
		}
		returnTo := strings.TrimSuffix(cfg.AuthPublicURL, "/") + "/api/v1/auth/login?lcid=" + url.QueryEscape(token)
		loginURL := strings.TrimSuffix(cfg.KratosUIURL, "/") + "/login?return_to=" + url.QueryEscape(returnTo)
		return c.Redirect(loginURL, http.StatusFound)
	}
	if whoamiResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(whoamiResp.Body)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "unexpected kratos session response",
			"details": string(body),
		})
	}

	var whoami kratosWhoamiResponse
	if err := json.NewDecoder(whoamiResp.Body).Decode(&whoami); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to decode kratos session"})
	}

	subject := strings.TrimSpace(whoami.Identity.ID)
	if subject == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "missing kratos identity"})
	}

	return acceptLogin(c, cfg.HydraAdminURL, challenge, subject)
}

func acceptLogin(c *fiber.Ctx, hydraAdminURL, challenge, subject string) error {
	payload := map[string]interface{}{
		"subject":      subject,
		"remember":     true,
		"remember_for": 3600,
	}

	body, _ := json.Marshal(payload)
	acceptURL := ory.HydraAdminURL(hydraAdminURL, "/oauth2/auth/requests/login/accept") + "?login_challenge=" + url.QueryEscape(challenge)
	req, err := http.NewRequest(http.MethodPut, acceptURL, strings.NewReader(string(body)))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to build login accept request"})
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to accept login request"})
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "login request not accepted",
			"details": string(body),
		})
	}

	var accept hydraLoginAcceptResponse
	if err := json.NewDecoder(resp.Body).Decode(&accept); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to decode login response"})
	}

	return c.Redirect(accept.RedirectTo, http.StatusFound)
}
