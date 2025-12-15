package auth

import (
	"encoding/json"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// SerializableCookie is a JSON-safe cookie representation
type SerializableCookie struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Domain   string  `json:"domain"`
	Path     string  `json:"path"`
	Expires  float64 `json:"expires,omitempty"`
	Secure   bool    `json:"secure"`
	HTTPOnly bool    `json:"httpOnly"`
}

// SaveCookies saves cookies to disk safely
func SaveCookies(page *rod.Page, path string) error {
	rawCookies, err := page.Browser().GetCookies()
	if err != nil {
		return err
	}

	var cookies []SerializableCookie
	for _, c := range rawCookies {
		cookies = append(cookies, SerializableCookie{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   c.Domain,
			Path:     c.Path,
			Expires:  float64(c.Expires), // ✅ explicit cast
			Secure:   c.Secure,
			HTTPOnly: c.HTTPOnly,
		})
	}

	data, err := json.MarshalIndent(cookies, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// LoadCookies restores cookies into the browser
func LoadCookies(page *rod.Page, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var cookies []SerializableCookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return err
	}

	var params []*proto.NetworkCookieParam
	for _, c := range cookies {
		param := &proto.NetworkCookieParam{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   c.Domain,
			Path:     c.Path,
			Secure:   c.Secure,
			HTTPOnly: c.HTTPOnly,
		}

		// Only set Expires if it exists
		if c.Expires > 0 {
			param.Expires = proto.TimeSinceEpoch(c.Expires) // ✅ explicit cast
		}

		params = append(params, param)
	}

	return page.Browser().SetCookies(params)
}
