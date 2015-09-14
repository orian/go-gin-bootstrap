package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/soy"
	"github.com/robfig/soy/data"
	"github.com/robfig/soy/soyhtml"

	"net/http"
	"time"
)

func GetDummyEndpoint(c *gin.Context) {
	resp := map[string]string{"hello": "world"}
	c.JSON(200, resp)
}

func ProvideSession(c *gin.Context) {
	// decodes a session id from cookie or provides a new one
	if ck, err := c.Request.Cookie("session"); err == nil {
		c.Set("session", ck.Value)
		return
	}
	newSession := "alfa romeo"
	c.Set("session", newSession)
	// c.Writer.C
	cookie := http.Cookie{}
	cookie.Name = "session"
	cookie.Value = newSession
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(time.Hour * 24 * 7)
	cookie.HttpOnly = true
	http.SetCookie(c.Writer, &cookie)
}

func SoyMust(b *soy.Bundle) *soyhtml.Tofu {
	t, err := b.CompileToTofu()
	if err != nil {
		panic(err)
	}
	return t
}

var tofu = SoyMust(soy.NewBundle().
	WatchFiles(true).            // watch soy files, reload on changes
	AddTemplateDir("templates")) // load *.soy in all sub-directories)

func RenderSoy(c *gin.Context) {
	var testdata = []data.Map{
		{"names": data.List{}},
		{"names": data.List{data.String("Rob")}},
		{"names": data.List{data.String("Rob"), data.String("Joe")}},
	}
	if err := tofu.Render(c.Writer, "soy.examples.simple.helloNames", testdata[2]); err != nil {
		c.Error(err)
		return
	}
}

func RenderPage(c *gin.Context) {
	if err := tofu.Render(c.Writer, "soy.examples.simple.page", nil); err != nil {
		c.Error(err)
		return
	}
}

func main() {
	api := gin.Default()
	r := api.Use(ProvideSession)
	r.GET("/dummy", GetDummyEndpoint)
	r.GET("/soy", RenderSoy)
	r.GET("/page", RenderPage)
	api.Run(":8080")
}
