package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	var (
		c *testClient
		r *gin.Engine
	)
	BeforeEach(func() {
		c = &testClient{}
		r = newRouter(c)
	})
	Specify("not post request returns 404", func() {
		body := strings.NewReader(`["http://127.0.0.1"]`)
		for _, method := range []string{
			http.MethodGet,
			http.MethodDelete,
			http.MethodPut,
			http.MethodConnect,
			http.MethodHead,
			http.MethodOptions,
			http.MethodPatch,
			http.MethodTrace,
		} {
			By("testing method " + method)
			req := httptest.NewRequest(method, "/", body)
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)
			c.expectNoGets()
			Expect(res.Code).To(Equal(http.StatusNotFound))
			Expect(res.Body.String()).To(BeEmpty())
		}
	})
	Specify("post request not to root returns 404", func() {
		body := strings.NewReader(`["http://127.0.0.1"]`)
		req := httptest.NewRequest(http.MethodPost, "/test", body)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		c.expectNoGets()
		Expect(res.Code).To(Equal(http.StatusNotFound))
		Expect(res.Body.String()).To(BeEmpty())
	})
	Specify("not array request body returns 400", func() {
		for _, body := range []string{
			``,
			`{}`,
			`123`,
			`"asd"`,
		} {
			By("testing body `" + body + "`")
			req := newReq(body)
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)
			c.expectNoGets()
			Expect(res.Code).To(Equal(http.StatusBadRequest))
			Expect(res.Body.String()).To(BeEmpty())
		}
	})
	Specify("empty json array body returns 200 with null json", func() {
		req := newReq(`[]`)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		c.expectNoGets()
		Expect(res.Code).To(Equal(http.StatusOK))
		Expect(res.Body.String()).To(Equal("null"))
	})
	Specify("when request is ok with different client returns", func() {
		c.addResponse(200, "text/html", `<html></html>`)
		c.addResponse(299, "text/html", `<html></html>`)
		c.addResponse(200, "text/html; charset=utf-8", `<html></html>`)
		c.addResponse(404, "text/html", `<html></html>`)
		c.addError("fail")
		c.addResponse(200, "TEXT/HTML", `<html></html>`)
		c.addResponse(200, "text/plain", `<html></html>`)
		req := newReq(`["http://yandex.ru","https://vk.com","http://example.com","https://yahoo.com","http://facebook.com","https://google.com","http://mail.ru"]`)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		c.expectGet("http://yandex.ru")
		c.expectGet("https://vk.com")
		c.expectGet("http://example.com")
		c.expectGet("https://yahoo.com")
		c.expectGet("http://facebook.com")
		c.expectGet("https://google.com")
		c.expectGet("http://mail.ru")
		c.expectNoGets()
		Expect(res.Code).To(Equal(http.StatusOK))
		expResponseBody, err := json.Marshal([]page{
			{
				URL:      "http://yandex.ru",
				Meta:     meta{Status: 200, ContentType: "text/html", ContentLength: 13},
				Elements: []element{{TagName: "html", Count: 1}},
			},
			{
				URL:      "https://vk.com",
				Meta:     meta{Status: 299, ContentType: "text/html", ContentLength: 13},
				Elements: []element{{TagName: "html", Count: 1}},
			},
			{
				URL:      "http://example.com",
				Meta:     meta{Status: 200, ContentType: "text/html; charset=utf-8", ContentLength: 13},
				Elements: []element{{TagName: "html", Count: 1}},
			},
			{
				URL:  "https://yahoo.com",
				Meta: meta{Status: 404},
			},
			{
				URL:  "https://google.com",
				Meta: meta{Status: 200, ContentType: "TEXT/HTML", ContentLength: 13},
			},
			{
				URL:  "http://mail.ru",
				Meta: meta{Status: 200, ContentType: "text/plain", ContentLength: 13},
			},
		})
		if err != nil {
			Fail("failed to json marshal expected response body: " + err.Error())
		}
		Expect(res.Body.String()).To(BeIdenticalTo(string(expResponseBody)), "response body")
	})
})

// newReq creates new http post / request with body
func newReq(body string) *http.Request {
	return httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
}

// tcReturn is a test client return
type tcReturn struct {
	res *http.Response
	err error
}

// testClient is test client for testing server
type testClient struct {
	gets    []string
	returns []tcReturn
}

// get implements client interface, it stores url, returns stored return and throws it
func (c *testClient) get(url string) (*http.Response, error) {
	if len(c.returns) == 0 {
		panic("no returns left")
	}
	c.gets = append(c.gets, url)
	r := c.returns[0]
	c.returns = c.returns[1:]
	return r.res, r.err
}

// expectNoGets tests that no get calls done
func (c *testClient) expectNoGets() {
	ExpectWithOffset(1, c.gets).To(BeEmpty(), "http client gets")
}

// expectGet tests that get call with url done and throws it from stored gets
func (c *testClient) expectGet(url string) {
	ExpectWithOffset(1, c.gets).ToNot(BeEmpty(), "http client gets")
	ExpectWithOffset(1, c.gets[0]).To(Equal(url), "get url")
	c.gets = c.gets[1:]
}

// addResponse adds test client return with response without error
func (c *testClient) addResponse(statusCode int, contentType string, body string) {
	r := httptest.NewRecorder()
	r.Code = statusCode
	r.Header().Set("Content-Type", contentType)
	r.Body = bytes.NewBuffer([]byte(body))
	c.returns = append(c.returns, tcReturn{res: r.Result(), err: nil})
}

// addResponse adds test client return with nil response and with error
func (c *testClient) addError(err string) {
	c.returns = append(c.returns, tcReturn{res: nil, err: errors.New(err)})
}
