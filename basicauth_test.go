package middleman

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBasicAuthenticator(t *testing.T) {
	// tests setup
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World!"))
	})
	mockAuthedUsername := "whitelisted"
	mockAuthedPass := "whitelisted"
	ba := NewBasicAuthenticator(func(uname, pw string) (string, error) {
		if uname != mockAuthedUsername || pw != mockAuthedPass {
			return "", errors.New("any error will cause the unauthed err")
		}
		return uname, nil
	})
	wrappedHandler := ba.Wrap(mockHandler)
	ts := httptest.NewServer(wrappedHandler)

	Convey("Test BasicAuthenticator()", t, func() {
		Convey("Test NewBasicAuthenticator() - Proxies anonymous func's success", func() {
			proxiedPass, err := ba.authenticate(mockAuthedUsername, mockAuthedPass)
			So(proxiedPass, ShouldEqual, mockAuthedUsername)
			So(err, ShouldBeNil)
		})
		Convey("Test NewBasicAuthenticator() - Proxies anonymous func's failure", func() {
			proxiedFail, err := ba.authenticate(mockAuthedUsername, "badpass")
			So(proxiedFail, ShouldEqual, "")
			So(err, ShouldNotBeNil)
		})
	})
	Convey("Test Wrap()", t, func() {
		Convey("Test Wrap() - http request fails with empty credentials", func() {
			baseReq, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
			resp, _ := http.DefaultClient.Do(baseReq)
			So(resp.StatusCode, ShouldEqual, http.StatusUnauthorized)
		})
		Convey("Test Wrap() - http request failes if unauthed", func() {
			baseReq, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
			baseReq.SetBasicAuth(mockAuthedUsername, "badpass")
			resp, _ := http.DefaultClient.Do(baseReq)
			So(resp.StatusCode, ShouldEqual, http.StatusUnauthorized)
		})
		Convey("Test Wrap() - http request succeeds if authed", func() {
			baseReq, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
			baseReq.SetBasicAuth(mockAuthedUsername, mockAuthedPass)
			resp, _ := http.DefaultClient.Do(baseReq)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
		})
	})
	Convey("Test GetUname()", t, func() {
		Convey("Test GetUname() - read context value", func() {
			baseReq, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
			baseReq = baseReq.WithContext(context.WithValue(baseReq.Context(), unameCtxKey, mockAuthedUsername))
			So(GetUname(baseReq), ShouldEqual, mockAuthedUsername)
		})
	})
}
