package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

var (
	Client       = http.Client{}
	HttpProxyURL string
	ShowIndex    bool
	Port         string
)

func init() {
	flag.StringVar(&HttpProxyURL, "proxy.http", "", "HTTP Proxy URL.Basic auth url as http://user:pass@proxy.com")
	flag.BoolVar(&ShowIndex, "github.index", false, "Show Github Index If access root.")
	flag.StringVar(&Port, "port", "3000", "Server runing port.")
	flag.Parse()
}

func main() {
	if HttpProxyURL != "" {
		purl, err := url.Parse(HttpProxyURL)
		if err != nil {
			fmt.Println("Proxy URL ERROR!")
			panic(err)
		}

		Client.Transport = &http.Transport{
			Proxy: http.ProxyURL(purl),
		}
	}

	e := echo.New()
	e.HideBanner = true
	e.Any("/*", proxy)
	if !ShowIndex {
		e.GET("/", func(context echo.Context) error {
			return context.String(200, "It works")
		})
	}
	e.Start(":" + Port)
}

func proxy(c echo.Context) error {
	path := c.Request().URL.String()
	if path == "index" {
		path = "/"
	}
	realDownloadUrl := "https://github.com" + strings.TrimPrefix(path, "/https://github.com")
	log.Println("OLD:", c.Request().URL.String(), "=>", realDownloadUrl)
	req, err := http.NewRequest(c.Request().Method, realDownloadUrl, c.Request().Body)
	if err != nil {
		log.Println(err)
		return err
	}

	resp, err := Client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	c.Response().Header().Add("Content-Disposition", resp.Header.Get("Content-Disposition"))
	return c.Stream(200, resp.Header.Get("Content-Type"), resp.Body)
}
