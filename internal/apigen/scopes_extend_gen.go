package apigen 

import "github.com/gofiber/fiber/v2"

type AuthFunc func(c *fiber.Ctx, rules ...string) error

func RegisterAuthFunc(app *fiber.App, f AuthFunc) {
	
	app.Post("/api/v1/clusters", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Get("/api/v1/clusters", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Put("/api/v1/clusters/:ID", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Delete("/api/v1/clusters/:ID", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Get("/api/v1/clusters/:ID", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Get("/api/v1/clusters/:ID/diagnostics", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Put("/api/v1/clusters/:ID/diagnostics/config", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Get("/api/v1/clusters/:ID/diagnostics/config", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Get("/api/v1/clusters/:ID/snapshot-config", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Put("/api/v1/clusters/:ID/snapshot-config", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Get("/api/v1/clusters/:ID/snapshots", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Post("/api/v1/clusters/:ID/snapshots", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Delete("/api/v1/clusters/:ID/snapshots/:snapshotId", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Post("/api/v1/clusters/:ID/snapshots/:snapshotId", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Post("/api/v1/databases", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Get("/api/v1/databases", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Delete("/api/v1/databases/:ID", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Get("/api/v1/databases/:ID", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Put("/api/v1/databases/:ID", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Get("/api/v1/databases/:ID/ddl-progress", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Post("/api/v1/databases/:ID/ddl-progress/:ddlID/cancel", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Post("/api/v1/databases/:ID/query", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
	app.Post("/api/v1/databases/test-connection", func(c *fiber.Ctx) error { 
		if c.Get("Authorization") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		} 
		if err := f(c); err != nil {
			return err
		}
		
		return c.Next()
	})
}
