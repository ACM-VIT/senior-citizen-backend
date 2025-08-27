package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Hydra endpoints (admin API)
var hydraAdminURL = "http://127.0.0.1:4445" // Hydra admin API
var kratosAdminURL = "http://127.0.0.1:4434" // Kratos admin API

// ConsentHandler handles Hydra consent flow
func ConsentHandler(c *fiber.Ctx) error {
	// 1. Get consent challenge from query
	challenge := c.Query("consent_challenge")
	if challenge == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing consent_challenge",
		})
	}

	// 2. Fetch consent request from Hydra
	resp, err := http.Get(fmt.Sprintf("%s/oauth2/auth/requests/consent?consent_challenge=%s", hydraAdminURL, url.QueryEscape(challenge)))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch consent request",
		})
	}
	defer resp.Body.Close()

	var consentReq struct {
		Subject string `json:"subject"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&consentReq); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to decode consent request",
		})
	}

	identityID := consentReq.Subject

	// 3. Fetch identity from Kratos
	kratosResp, err := http.Get(fmt.Sprintf("%s/admin/identities/%s", kratosAdminURL, identityID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch kratos identity",
		})
	}
	defer kratosResp.Body.Close()

	var kratosIdentity struct {
		ID     string                 `json:"id"`
		Traits map[string]interface{} `json:"traits"`
	}
	if err := json.NewDecoder(kratosResp.Body).Decode(&kratosIdentity); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to decode kratos identity",
		})
	}

	// 4. Extract traits (email, admin)
	email, _ := kratosIdentity.Traits["email"].(string)
	admin, _ := kratosIdentity.Traits["admin"].(bool)

	// 5. Accept consent and include traits as claims
	acceptPayload := map[string]interface{}{
		"grant_scope": []string{"openid", "profile", "email", "admin"},
		"remember":    true,
		"remember_for": 3600,
		"session": map[string]interface{}{
			"id_token": map[string]interface{}{
				"email": email,
				"admin": admin,
			},
			"access_token": map[string]interface{}{
				"email": email,
				"admin": admin,
			},
		},
	}

	bodyBytes, _ := json.Marshal(acceptPayload)
	acceptResp, err := http.Post(
		fmt.Sprintf("%s/oauth2/auth/requests/consent/accept?consent_challenge=%s", hydraAdminURL, url.QueryEscape(challenge)),
		"application/json",
		strings.NewReader(string(bodyBytes)),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to accept consent request",
		})
	}
	defer acceptResp.Body.Close()

	var acceptResult struct {
		RedirectTo string `json:"redirect_to"`
	}
	if err := json.NewDecoder(acceptResp.Body).Decode(&acceptResult); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to decode accept response",
		})
	}

	// 6. Redirect user back to Hydra
	return c.Redirect(acceptResult.RedirectTo, http.StatusFound)
}