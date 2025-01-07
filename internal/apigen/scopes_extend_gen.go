package apigen 

import "github.com/gofiber/fiber/v2"

type AuthFunc func(c *fiber.Ctx, rules ...string) error

func RegisterAuthFunc(app *fiber.App, f AuthFunc) {
	
	app.Get("/api/v1/clusters", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return c.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		
		return c.Next()
	})
}
