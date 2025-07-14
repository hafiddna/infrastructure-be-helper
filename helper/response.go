package helper

import (
	"github.com/gofiber/fiber/v2"
)

type ResponseStruct struct {
	Ctx        *fiber.Ctx
	StatusCode int
	Message    string
	Error      interface{}
	Data       interface{}
}

type BaseResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Error      interface{} `json:"error"`
	Data       interface{} `json:"data"`
}

func SendResponse(baseResponse ResponseStruct) (err error) {
	newBaseResponse := BaseResponse{
		StatusCode: baseResponse.StatusCode,
	}

	if baseResponse.Message != "" {
		newBaseResponse.Message = baseResponse.Message
	}

	if baseResponse.Error != nil {
		newBaseResponse.Error = baseResponse.Error
	}

	if baseResponse.Data != nil {
		//if config.Config.App.Environment == "development" {
		newBaseResponse.Data = baseResponse.Data
		//} else {
		//	marshalledData := JSONMarshal(baseResponse.Data)
		//
		//	newBaseResponse.Data, err = EncryptAES256CBC([]byte(marshalledData), []byte(config.Config.App.Secret.DataEncryptionKey))
		//	if err != nil {
		//		log.Fatalf("Error encrypting data: %v", err)
		//	}
		//}
	} else {
		newBaseResponse.Data = nil
	}

	baseResponse.Ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	baseResponse.Ctx.Status(baseResponse.StatusCode)
	return baseResponse.Ctx.JSON(newBaseResponse)
}
