package main

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Link struct {
	Short    string `gorm:"primaryKey"`
	Original string
}

const (
	charset  = "abcdefghijklmnopqrstuvwxyz0123456789"
	shortLen = 6
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	godotenv.Load()

	dsn := os.Getenv("DBASE_PROD")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Link{})

	engine := html.New("./templates", ".html")
	engine.Reload(true)

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./static")
	app.Use(compress.New())

	app.Get("/", func(c *fiber.Ctx) error {
		link := c.Query("link", c.FormValue("link"))
		link = strings.TrimSpace(link)

		if link == "" {
			if strings.Contains(c.Get("Accept"), "html") {
				return c.Render("index", fiber.Map{})
			}
			return c.SendString("Welcome to the UISM.link! Its a link shortener, simple, fast, free right from CLI (and browser for sure). To short a link please use this url format: https://uism.link/?link=*your link here*, you will get a short url that leads to the original one. Thanks for using!")
		}

		if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
			link = "https://" + link
		}

		accept := c.Get("Accept")
		if !(strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://")) || !strings.Contains(link, ".") {
			return respondError(c, 400, "Hmmm, this does not look like a link")
		}

		var existing Link
		err := db.First(&existing, "original = ?", link).Error
		if err == nil {
			shortURL := c.BaseURL() + "/" + existing.Short
			accept := c.Get("Accept")

			if strings.Contains(accept, "html") {
				return c.Render("index", fiber.Map{
					"ShortURL": shortURL,
				})
			}

			return c.SendString(shortURL)
		} else if err != gorm.ErrRecordNotFound {
			return respondError(c, 500, "Internal linkbase error!")
		}

		var short string
		for {
			b := make([]byte, shortLen)
			for i := range b {
				b[i] = charset[rand.Intn(len(charset))]
			}
			short = string(b)

			var exists Link
			if err := db.First(&exists, "short = ?", short).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					break
				}
			}
		}

		if err := db.Create(&Link{
			Short:    short,
			Original: link,
		}).Error; err != nil {
			return respondError(c, 500, "Error saving a new link to linkbase!")
		}

		shortURL := c.BaseURL() + "/" + short
		accept = c.Get("Accept")
		if strings.Contains(accept, "html") {
			return c.Render("index", fiber.Map{
				"ShortURL": shortURL,
			})
		}
		return c.SendString(shortURL)
	})

	app.Get("/api", func(c *fiber.Ctx) error {
		return c.Render("api", fiber.Map{})
	})

	app.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var link Link
		if err := db.First(&link, "short = ?", id).Error; err != nil {
			return c.Status(404).SendString("Not found")
		}
		return c.Redirect(link.Original, fiber.StatusMovedPermanently)
	})

	log.Fatal(app.Listen(":8080"))
}

func respondError(c *fiber.Ctx, status int, msg string) error {
	if strings.Contains(c.Get("Accept"), "html") {
		return c.Status(status).Render("index", fiber.Map{
			"Error": msg,
		})
	}
	return c.Status(status).SendString(msg)
}
