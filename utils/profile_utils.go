package utils

import (
	"strings"
)

func ProfileSlugify(firstName, lastName string) string {

	fullName := strings.TrimSpace(firstName + " " + lastName)
	slug := strings.ToLower(fullName)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}