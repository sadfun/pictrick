package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"pictrick/mongo"
)

func Get(c *fiber.Ctx) error  {
	param := c.Params("oid")
	if param == "" {
		return c.SendStatus(fiber.StatusNoContent)
	}

	id, err := uuid.Parse(param)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	payload, contentType, err := mongo.GetFilePayload(id)
	if err != nil {
		switch err {
		case mongo.FileNotFoundError:
			return c.SendStatus(fiber.StatusNotFound)

		default:
			log.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	if contentType == "" {
		c.Set("Content-Type", "application/octet-stream")
	} else {
		c.Set("Content-Type", contentType)
	}

	return c.Send(payload)
}