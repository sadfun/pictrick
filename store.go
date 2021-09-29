package main

import (
	"bufio"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"io"
	"pictrick/mongo"
	"sort"
)

func Store(c *fiber.Ctx) error  {
	file, err := c.FormFile("document")

	if err != nil {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusNoContent)
	}

	if file.Size > Config.MaxFileSize {
		return c.SendStatus(fiber.StatusRequestEntityTooLarge)
	}

	contentType := utils.CopyString(file.Header.Get("Content-Type"))
	if contentType == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if Config.Formats[sort.SearchStrings(Config.Formats, contentType)] != contentType {
		return c.SendStatus(fiber.StatusUnsupportedMediaType)
	}

	f, err := file.Open()
	if err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	payload, err := io.ReadAll(bufio.NewReader(f))
	if err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	id, err := mongo.SaveFile(payload, file.Header.Get("Content-Type"), c.IP())
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id.String(),
	})
}
