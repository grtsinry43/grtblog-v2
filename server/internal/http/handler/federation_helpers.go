package handler

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func parseFederationRequest(c *fiber.Ctx) (*http.Request, error) {
	host := string(c.Context().Host())
	scheme := c.Protocol()
	if scheme == "" {
		scheme = "https"
	}
	rawQuery := string(c.Context().QueryArgs().QueryString())

	u := &url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     c.Path(),
		RawQuery: rawQuery,
	}
	req, err := http.NewRequest(c.Method(), u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header = make(http.Header)
	c.Request().Header.VisitAll(func(k, v []byte) {
		req.Header.Add(string(k), string(v))
	})
	req.Host = host
	req.RequestURI = c.OriginalURL()
	return req, nil
}

func sameBaseURL(a string, b string) bool {
	ua, errA := url.Parse(strings.TrimRight(a, "/"))
	ub, errB := url.Parse(strings.TrimRight(b, "/"))
	if errA != nil || errB != nil {
		return false
	}
	return strings.EqualFold(ua.Scheme, ub.Scheme) && strings.EqualFold(ua.Host, ub.Host)
}
