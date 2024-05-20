package routes

import (
	"errors"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GlobalRoute(router *gin.Engine) {
	router.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	router.Use(func(ctx *gin.Context) {
		if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
			hub.Scope().SetTag("someRandomTag", "maybeYouNeedIt")
		}
		ctx.Next()
	})

	router.GET("/", func(ctx *gin.Context) {
		if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("unwantedQuery", "someQueryDataMaybe")
				hub.CaptureMessage("User provided unwanted query string, but we recovered just fine")
			})
		}
		// create error using New() function
		_ = errors.New("WRONG MESSAGE")
		log.Fatal("test")
		ctx.Status(http.StatusOK)
	})
}
