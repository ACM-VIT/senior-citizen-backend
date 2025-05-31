package helpers

import (
	"fmt"
)

func ValidateRole(role string) error {
	validRoles := map[string]bool{
		"community": true,
		"jc":        true,
		"sc":        true,
		"board":     true,
	}
	if !validRoles[role] {
		return fmt.Errorf("invalid role: %s", role)
	}
	return nil
}

func ValidatePosition(position string) error {
	if position == "" {
		return nil
	}

	validPositions := map[string]bool{
		"chairperson": true,
		"vice-chairperson": true,
		"secretary": true,
		"cosecretary": true,
		"technical-director": true,
		"projects-lead": true,
		"research-lead": true,
		"contest-lead": true,
		"design-lead": true,
		"creative-lead": true,
		"acmw-chairperson": true,
		"acmw-vice-chairperson": true,
		"acmw-secretary": true,
	}
	if !validPositions[position] {
		return fmt.Errorf("invalid position: %s", position)
	}
	return nil
}
