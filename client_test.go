package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHealth(t *testing.T) {
	Convey("Given a health check request send to "+constants.APIHeathEndpoint, t, func() {
		request := httptest.NewRequest("GET", constants.APIHeathEndpoint, nil)
		response := httptest.NewRecorder()
		route := &Route{
			Method:  "GET",
			Path:    constants.APIHeathEndpoint,
			Handler: http.HandlerFunc(Health),
		}
		Convey("When the request is handled by the router", func() {
			route.Test(response, request)
			Convey("Then we should get expected reponse and success http response code", func() {
				So(response.Code, ShouldEqual, 200)
				So(string(cleanResponse(response.Body.String())), ShouldEqual, healthExpectedResponse)
			})
		})

	})
}
