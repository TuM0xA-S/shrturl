package main

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/teris-io/shortid"
)

type Shrt struct {
	db      *redis.Client
	baseURL string
	expTime time.Duration
}

type ShrtCfg struct {
	DB             *redis.Client
	BaseURL        string
	ExpirationTime time.Duration
}

func NewShrt(ctx context.Context, cfg ShrtCfg) *Shrt {
	const api = "NewShrt"

	return &Shrt{db: cfg.DB, baseURL: cfg.BaseURL, expTime: cfg.ExpirationTime}
}

func (s *Shrt) Setup(e *echo.Echo) {
	e.GET("/", s.root)
	e.POST("/", s.makeShrt)
	e.GET("/:code", s.useShrt)
}

func (s *Shrt) root(c echo.Context) error {
	return c.String(http.StatusOK, "shrturl 0.01\n")
}

func (s *Shrt) makeShrt(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var req makeShrtReq
	if err := (&echo.DefaultJSONSerializer{}).Deserialize(c, &req); err != nil {
		return err
	}
	if u, err := url.Parse(req.URL); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if u.Host == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid url: no host")
	}

	code, err := shortid.Generate()
	if err != nil {
		return err
	}

	var resp makeShrtResp
	resp.URL = req.URL
	resp.ShrtURL = s.baseURL + code
	if s.expTime != 0 {
		resp.ValidUntil = time.Now().Add(s.expTime)
	}

	if err := s.db.Set(ctx, code, req.URL, s.expTime).Err(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Shrt) useShrt(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	code := c.Param("code")
	val, err := s.db.Get(ctx, code).Result()
	if err != nil {
		if err == redis.Nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return err
	}
	return c.Redirect(http.StatusFound, val)
}

type makeShrtReq struct {
	URL string `json:"url"`
}

type makeShrtResp struct {
	ShrtURL    string    `json:"shrt_url"`
	URL        string    `json:"url"`
	ValidUntil time.Time `json:"valid_until"`
}
