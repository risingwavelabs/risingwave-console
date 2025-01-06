package middleware

import (
	"fmt"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/risingwavelabs/wavekit/internal/service"
)

var re = regexp.MustCompile(`[\r?\n| ]`)

var loggerConfig = logger.Config{
	CustomTags: map[string]logger.LogFunc{
		"statusColored": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
			code := c.Response().StatusCode()
			var color string
			if 200 <= code && code < 300 {
				color = fiber.DefaultColors.Green
			} else if 400 <= code && code < 500 {
				color = fiber.DefaultColors.Yellow
			} else {
				color = fiber.DefaultColors.Red
			}
			return output.WriteString(fmt.Sprintf("%s%3d%s", color, code, fiber.DefaultColors.Reset))
		},
		"methodColored": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
			method := c.Method()
			var color string
			switch method {
			case fiber.MethodGet:
				color = fiber.DefaultColors.White
			case fiber.MethodPost:
				color = fiber.DefaultColors.Green
			case fiber.MethodPut:
				color = fiber.DefaultColors.Yellow
			case fiber.MethodDelete:
				color = fiber.DefaultColors.Red
			case fiber.MethodPatch:
				color = fiber.DefaultColors.Cyan
			case fiber.MethodHead:
				color = fiber.DefaultColors.Magenta
			case fiber.MethodOptions:
				color = fiber.DefaultColors.Blue
			default:
				color = fiber.DefaultColors.Reset
			}
			return output.WriteString(fmt.Sprintf("%s%s%s", color, method, fiber.DefaultColors.Reset))
		},
		"flatReqBody": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
			b := re.ReplaceAllString(string(c.Body()), "")
			if len(b) > 512 {
				b = b[:512] + "...(truncated)"
			}
			return output.WriteString(b)
		},
		"flatResBody": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
			b := re.ReplaceAllString(string(c.Response().Body()), "")
			if len(b) > 512 {
				b = b[:512] + "...(truncated)"
			}
			return output.WriteString(b)
		},
		"user": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
			user, ok := c.Locals(UserContextKey).(*service.User)
			if !ok {
				return output.WriteString("unknown user")
			}
			return output.WriteString(fmt.Sprintf("{\"id\": \"%d\"}", user.ID))
		},
	},
	TimeZone: "Asia/Shanghai",
	Format:   "${time} | ${ip} | ${statusColored} | ${latency} | ${methodColored} | ${path} | ${locals:requestid} |${cyan} ${user} ${reset}|${blue} ${queryParams} |  ${flatReqBody} ${reset}| ${reqHeaders} |${green} ${flatResBody} ${reset}|${red} ${error} ${reset}|\n",
}

func NewLogger() func(*fiber.Ctx) error {
	return logger.New(loggerConfig)
}
