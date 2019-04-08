package middleman

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLogger(t *testing.T) {
	// tests setup
	mockResp := "Hello World"
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResp))
	})

	Convey("Test NewLogger()", t, func() {
		Convey("Test NewLogger() - config values are set", func() {
			alltruelogger := NewLogger(true, true, true)
			So(alltruelogger.logContentLength, ShouldEqual, true)
			So(alltruelogger.logStatus, ShouldEqual, true)
			So(alltruelogger.logDuration, ShouldEqual, true)
			allfalselogger := NewLogger(false, false, false)
			So(allfalselogger.logContentLength, ShouldEqual, false)
			So(allfalselogger.logStatus, ShouldEqual, false)
			So(allfalselogger.logDuration, ShouldEqual, false)
		})
	})
	Convey("Test Wrap()", t, func() {
		Convey("Test Wrap() - server still responds to requests", func() {
			lmw := NewLogger(true, true, true)
			wrappedHandler := lmw.Wrap(mockHandler)
			ts := httptest.NewServer(wrappedHandler)
			req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
			resp, _ := http.DefaultClient.Do(req)
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			So(string(bodyBytes), ShouldEqual, mockResp)
		})
	})
}
