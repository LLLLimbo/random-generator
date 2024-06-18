package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//go:embed index.html
var fsys embed.FS

var (
	appId    = flag.String("app-id", "random-generator-1", "app id")
	idc      = flag.String("idc", "http://localhost:29705", "idc host")
	ac       = flag.String("auth-center", "http://localhost:29706", "auth center host")
	c1       = flag.Bool("auto-redirect-on-unauthorized", false, "auto redirect on unauthorized")
	port     = flag.String("port", "3000", "app's running port")
	redirect = flag.String("redirect", "http://localhost:3000/index.html", "redirect url after got token")
)

func main() {
	r := gin.Default()
	flag.Parse()

	// Generates a random string
	rand.Seed(time.Now().UnixNano())
	randString := func(n int) string {
		var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		s := make([]rune, n)
		for i := range s {
			s[i] = letters[rand.Intn(len(letters))]
		}
		return string(s)
	}

	// Handler for /random route
	r.GET("/api/random", func(c *gin.Context) {
		sessionId, _ := c.Cookie("session_id")
		resp, err := SessionValidate(sessionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}

		if *c1 && !resp.Active {
			//redirect to login
			t := fmt.Sprintf("%s/idc/redirect/login?appId=%s", *idc, *appId)
			c.Redirect(http.StatusFound, t)
		}

		if !*c1 && !resp.Active {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"random": randString(6)})
	})

	r.GET("/api/token", func(c *gin.Context) {
		key := c.Query("key")
		resp, err := MakeGetTokenRequest(key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}
		log.Printf("Got token %s", resp.Token)
		c.SetCookie("token", resp.Token, 36000, "/", "localhost", false, true)
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/index.html", *redirect))
	})

	// serve static files
	fs := http.FS(fsys)
	r.StaticFS("/index.html", fs)

	err := r.Run(fmt.Sprintf(":%s", *port))
	if err != nil {
		panic(err)
	}
}
