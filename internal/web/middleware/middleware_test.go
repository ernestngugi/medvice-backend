package middleware

import (
	"net/http"
	"testing"

	"github.com/ernestngugi/medvice-backend/internal/testutils"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMiddleware(t *testing.T) {

	Convey("TestMiddleware", t, func() {


		testRouter := gin.Default()
		testRouter.Use(DefaultMiddlewares()...)

		Convey("can set request id", func() {
			w, err := testutils.DoRequest(testRouter, http.MethodGet, "test-router", nil)
			So(err, ShouldBeNil)

			So(w.Header().Get("x-request-id"), ShouldNotBeEmpty)
		})

		Convey("can set application cors", func() {
			w, err := testutils.DoRequest(testRouter, http.MethodGet, "test-router", nil)
			So(err, ShouldBeNil)

			So(w.Header().Get("Access-Control-Allow-Origin"), ShouldEqual, "*")
			So(w.Header().Get("Access-Control-Allow-Credentials"), ShouldEqual, "true")
			So(w.Header().Get("Access-Control-Allow-Methods"), ShouldEqual, "POST, OPTIONS, GET, PUT, DELETE")
		})
	})
}
