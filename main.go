package main

import (
	"github.com/boombuler/barcode"
	"html/template"
	"image/png"
	"log"
	"net/http"
	"os"
	"snsmod/util"
	"strings"

	"github.com/gin-gonic/gin"
	//"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", func(ctx *gin.Context) {
		html := template.Must(template.ParseFiles("templates/navbar.html", "templates/index.html"))
		r.SetHTMLTemplate(html)
		ctx.HTML(http.StatusOK, "navbar.html", gin.H{})
	})

	r.GET("/mkqrcode", func(ctx *gin.Context) {
		html := template.Must(template.ParseFiles("templates/navbar.html", "templates/mkqrcode.html"))
		r.SetHTMLTemplate(html)
		ctx.HTML(http.StatusOK, "navbar.html", gin.H{})
	})

	r.POST("/save", func(ctx *gin.Context) {
		m := util.MySnsData{}

		m.Facebook = strings.TrimSpace(ctx.PostForm("facebook"))
		m.Twitter = strings.TrimSpace(ctx.PostForm("twitter"))
		m.Instagram = strings.TrimSpace(ctx.PostForm("instagram"))
		m.Line = strings.TrimSpace(ctx.PostForm("line"))

		_, err := util.SaveUserItem(m)
		if err != nil {
			log.Print("An error has occurred: %s\n", err)
		}

		// Firestore 書き込み
		ctx.Redirect(http.StatusFound, "settings")
	})

	r.GET("/settings", func(ctx *gin.Context) {
		//m := util.MySnsData{}
		m, _ := util.GetUserItem()
		//if err != nil {
		//	log.Print("An error has occurred: %s\n", err)
		//	//return
		//}
		html := template.Must(template.ParseFiles("templates/navbar.html", "templates/settings.html"))

		r.SetHTMLTemplate(html)
		ctx.HTML(http.StatusOK, "navbar.html", gin.H{"mySnsData": m})
	})

	r.GET("/qrcode/:snstype", func(ctx *gin.Context) {
		snsType := ctx.Param("snstype")
		url := ""

		m := util.MySnsData{}
		m, err := util.GetUserItem()
		if err != nil {
			log.Print("An error has occurred: %s\n", err)
			return
		}

		switch snsType {
		case "facebook":
			url = "fb://profile/" + m.Facebook
		case "twitter":
			url = "twitter://user?screen_name=" + m.Twitter
		case "instagram":
			url = "https://instagram.com/" + m.Instagram
		case "line":
			url = "https://line.me/ti/p/" + m.Line
		default:
			log.Print("An error has occurred: %s\n", err)
			return
		}

		qrCode, err := qr.Encode(url, qr.L, qr.Auto)
		if err != nil {
			log.Print("An error has occurred: %s\n", err)
			return
		}

		qrCode, err = barcode.Scale(qrCode, 512, 512)
		if err != nil {
			log.Print("An error has occurred: %s\n", err)
			return
		}

		file, err := os.Create("./assets/images/qrcode.png")
		if err != nil {
			log.Print("An error has occurred: %s\n", err)
			return
		}

		defer file.Close()
		png.Encode(file, qrCode)
		html := template.Must(template.ParseFiles("templates/navbar.html", "templates/qrcode.html"))

		r.SetHTMLTemplate(html)
		ctx.HTML(http.StatusOK, "navbar.html", gin.H{})
	})

	r.POST("/alldiscard", func(ctx *gin.Context) {
		err := util.AllDiscard()
		if err != nil {
			return
		}
		ctx.Redirect(302, "/settings")
	})

	return r
}

func main() {
	r := setupRouter()
	_ = r.Run(":8080")
}
