package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func newRouter(c client) *gin.Engine {
	s := server{client: c}
	r := gin.New()
	r.POST("/", s.postRoot)
	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
	return r
}

type server struct {
	client client
}

func (s server) postRoot(c *gin.Context) {
	urlsJSON, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var urls []string
	if err := json.Unmarshal(urlsJSON, &urls); err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var pages []page
	for _, url := range urls {
		res, err := s.client.get(url)
		if err != nil {
			c.Error(err)
			continue
		}
		var p page
		p.URL = url
		p.Meta.Status = res.StatusCode
		if res.StatusCode >= 200 && res.StatusCode < 300 {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				c.Error(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			if err := res.Body.Close(); err != nil {
				c.Error(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			p.Meta.ContentType = res.Header.Get("Content-Type")
			p.Meta.ContentLength = len(body)
			if p.Meta.ContentType == "text/html" || strings.HasPrefix(p.Meta.ContentType, "text/html;") {
				p.Elements, err = parse(body)
				if err != nil {
					c.Error(err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
			}
		}
		pages = append(pages, p)
	}
	c.JSON(http.StatusOK, pages)
}
