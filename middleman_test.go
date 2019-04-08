package middleman

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMiddleman(t *testing.T) {
	// tests setup
	mockResp := "Hello World"
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResp))
	})
	middlewareCount := 5
	var mw []Middleware
	for i := 0; i < middlewareCount; i++ {
		mw = append(mw, NewFunctional(func() {}, func() {}))
	}
	mm := NewMiddleman(mw...)
	wrappedHandler := mm.Wrap(mockHandler)
	ts := httptest.NewServer(wrappedHandler)

	Convey("Test NewMiddleman()", t, func() {
		Convey("Test NewMiddleman() - middleware len", func() {
			So(len(mm.middleware), ShouldEqual, middlewareCount)
		})
	})
	Convey("Test Wrap()", t, func() {
		Convey("Test NewMiddleman() - middleware operations happen", func() {
			req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
			resp, _ := http.DefaultClient.Do(req)
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			So(string(bodyBytes), ShouldEqual, mockResp)
		})
	})
}
