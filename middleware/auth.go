package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hafiddna/infrastructure-be-helper/config"
	"github.com/hafiddna/infrastructure-be-helper/helper"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")
		if authorization == "" || len(authorization) < 7 {
			return helper.SendResponse(helper.ResponseStruct{
				Ctx:        c,
				StatusCode: fiber.StatusUnauthorized,
				Data:       nil,
				Message:    "Unauthorized",
			})
		}

		token := authorization[7:]
		aToken, err := helper.ValidateRS512Token(config.Config.App.JWT.PublicKey, token)
		if err != nil {
			return helper.SendResponse(helper.ResponseStruct{
				Ctx:        c,
				StatusCode: fiber.StatusUnauthorized,
				Message:    "Unauthorized",
				Data:       nil,
				Error:      err.Error(),
			})
		}

		if !aToken.Valid {
			return helper.SendResponse(helper.ResponseStruct{
				Ctx:        c,
				StatusCode: fiber.StatusUnauthorized,
				Data:       nil,
				Message:    "Unauthorized",
			})
		}

		mapStringClaims := make(map[string]interface{})
		for key, value := range aToken.Claims.(jwt.MapClaims) {
			mapStringClaims[key] = value
		}

		//var encryptedData helper.EncryptedData
		//tokenData := helper.JSONMarshal(mapStringClaims["data"])
		//err = helper.JSONUnmarshal([]byte(tokenData), &encryptedData)
		//if err != nil {
		//	return helper.SendResponse(helper.ResponseStruct{
		//		Ctx:        c,
		//		StatusCode: fiber.StatusUnauthorized,
		//		Message:    "Unauthorized",
		//		Error:      err.Error(),
		//	})
		//}
		//decryptedData, err := helper.DecryptAES256CBC(&encryptedData, []byte(config.Config.App.Secret.AuthKey))
		//if err != nil {
		//	return helper.SendResponse(helper.ResponseStruct{
		//		Ctx:        c,
		//		StatusCode: fiber.StatusInternalServerError,
		//		Message:    "Internal Server Error",
		//		Error:      err.Error(),
		//	})
		//}

		//mapDecryptedData := make(map[string]interface{})
		//mapDecryptedData["sub"] = mapStringClaims["sub"]
		//err = helper.JSONUnmarshal([]byte(decryptedData), &mapDecryptedData)
		//if err != nil {
		//	return helper.SendResponse(helper.ResponseStruct{
		//		Ctx:        c,
		//		StatusCode: fiber.StatusUnauthorized,
		//		Message:    "Unauthorized",
		//		Error:      err.Error(),
		//	})
		//}

		c.Locals("user", mapStringClaims)

		return c.Next()
	}
}

func AuthWSMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Query("X-Token")
		if authorization == "" || len(authorization) < 7 {
			return fiber.ErrUnauthorized
		}

		token := authorization[7:]
		aToken, err := helper.ValidateRS512Token(config.Config.App.JWT.PublicKey, token)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		if !aToken.Valid {
			return fiber.ErrUnauthorized
		}

		mapStringClaims := make(map[string]interface{})
		for key, value := range aToken.Claims.(jwt.MapClaims) {
			mapStringClaims[key] = value
		}

		//var encryptedData helper.EncryptedData
		//tokenData := helper.JSONMarshal(mapStringClaims["data"])
		//err = helper.JSONUnmarshal([]byte(tokenData), &encryptedData)
		//if err != nil {
		//	return fiber.ErrUnauthorized
		//}
		//decryptedData, err := helper.DecryptAES256CBC(&encryptedData, []byte(config.Config.App.Secret.AuthKey))
		//if err != nil {
		//	return fiber.ErrInternalServerError
		//	// TODO: Is this needed?
		//	//return helper.SendResponse(helper.ResponseStruct{
		//	//	Ctx:        c,
		//	//	StatusCode: fiber.StatusInternalServerError,
		//	//	Message:    "Internal Server Error",
		//	//	Error:      err.Error(),
		//	//})
		//}

		//mapDecryptedData := make(map[string]interface{})
		//mapDecryptedData["sub"] = mapStringClaims["sub"]
		//err = helper.JSONUnmarshal([]byte(decryptedData), &mapDecryptedData)
		//if err != nil {
		//	return fiber.ErrUnauthorized
		//}

		c.Locals("user", mapStringClaims)

		return c.Next()
	}
}

func RoleAuthMiddleware(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenRole, ok := c.Locals("user").(map[string]interface{})["data"].(map[string]interface{})["roles"].([]interface{})
		if !ok {
			return helper.SendResponse(helper.ResponseStruct{
				Ctx:        c,
				StatusCode: fiber.StatusForbidden,
				Message:    "Forbidden",
				Data:       nil,
			})
		}

		for _, v := range roles {
			if !helper.ArrayInterfaceContains(tokenRole, v) {
				return helper.SendResponse(helper.ResponseStruct{
					Ctx:        c,
					StatusCode: fiber.StatusForbidden,
					Message:    "Forbidden",
					Data:       nil,
				})
			}
		}

		return c.Next()
	}
}
