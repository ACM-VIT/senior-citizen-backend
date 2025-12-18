package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"acm-portal/internal/config"

	"github.com/gofiber/fiber/v2"
)

func ConsentHandler(c *fiber.Ctx) error {
	challenge := c.Query("consent_challenge")
	if challenge == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing consent_challenge"})
	}

	cfg := config.Get()
	if cfg.HydraAdminURL == "" || cfg.KratosAdminURL == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "missing auth service configuration"})
	}

	hydraAdminURL := strings.TrimSuffix(cfg.HydraAdminURL, "/")
	kratosAdminURL := strings.TrimSuffix(cfg.KratosAdminURL, "/")

	consentURL := fmt.Sprintf("%s/admin/oauth2/auth/requests/consent?consent_challenge=%s", hydraAdminURL, url.QueryEscape(challenge))
	resp, err := http.Get(consentURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch consent request"})
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid consent challenge"})
	}

	var consentReq struct {
		Subject                      string   `json:"subject"`
		RequestedScope               []string `json:"requested_scope"`
		RequestedAccessTokenAudience []string `json:"requested_access_token_audience"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&consentReq); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to decode consent request"})
	}

	identityID := consentReq.Subject

	kratosResp, err := http.Get(fmt.Sprintf("%s/admin/identities/%s", kratosAdminURL, identityID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch kratos identity"})
	}
	defer kratosResp.Body.Close()
	if kratosResp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid kratos identity"})
	}

	var kratosIdentity struct {
		ID     string                 `json:"id"`
		Traits map[string]interface{} `json:"traits"`
	}
	if err := json.NewDecoder(kratosResp.Body).Decode(&kratosIdentity); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to decode kratos identity"})
	}

	email, _ := kratosIdentity.Traits["email"].(string)
	if strings.TrimSpace(email) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing email in identity traits"})
	}
	acceptPayload := map[string]interface{}{
		"grant_scope":  consentReq.RequestedScope,
		"remember":     true,
		"remember_for": 3600,
		"session": map[string]interface{}{
			"id_token": map[string]interface{}{
				"email": email,
			},
			"access_token": map[string]interface{}{
				"email": email,
			},
		},
	}

	bodyBytes, _ := json.Marshal(acceptPayload)
	acceptResp, err := http.Post(
		fmt.Sprintf("%s/admin/oauth2/auth/requests/consent/accept?consent_challenge=%s", hydraAdminURL, url.QueryEscape(challenge)),
		"application/json",
		strings.NewReader(string(bodyBytes)),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to accept consent request"})
	}
	defer acceptResp.Body.Close()
	if acceptResp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "consent request not accepted"})
	}

	var acceptResult struct {
		RedirectTo string `json:"redirect_to"`
	}
	if err := json.NewDecoder(acceptResp.Body).Decode(&acceptResult); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to decode accept response"})
	}

	return c.Redirect(acceptResult.RedirectTo, http.StatusFound)
}
