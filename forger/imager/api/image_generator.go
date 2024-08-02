package api

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"forger/imager/models"
	"forger/imager/utils"
	"image"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ImageRequest struct {
	ImageBase64 string `json:"image_base64"`
}

func BuildIcon(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var requestBody map[string]interface{}

	// Decode the base64-encoded body
	decodedBody, err := base64.StdEncoding.DecodeString(req.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Error decoding base64 body",
		}
	}

	// Unmarshal the decoded body into a map
	err = json.Unmarshal(decodedBody, &requestBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Error unmarshalling JSON body",
		}
	}

	var imgReq ImageRequest
	err = json.Unmarshal(decodedBody, &imgReq)
	if err != nil {
		return responseBuilder(0, req.Body, "invalid JSON body", "invalid JSON body")
	}

	imageData, err := base64.StdEncoding.DecodeString(imgReq.ImageBase64)
	if err != nil {
		return responseBuilder(0, nil, "invalid base64 encoding", "invalid base64 encoding")
	}

	srcImage, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return responseBuilder(0, nil, "invalid image data", "invalid image data")
	}

	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)

	err = utils.IOSmageResizer(zipWriter, srcImage, models.IOSResizeMetaList, utils.IOS)
	if err != nil {
		return responseBuilder(0, nil, "IOSmageResizer error", "IOSmageResizer error")
	}

	err = utils.IOSmageResizer(zipWriter, srcImage, models.AndroidResizeMetaList, utils.Android)
	if err != nil {
		return responseBuilder(0, nil, "IOSmageResizer error", "IOSmageResizer error")
	}

	err = zipWriter.Close()
	if err != nil {
		return responseBuilder(0, nil, "zipWriter error", "zipWriter error")
	}

	// Prepare response headers and content
	resp := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "text/plain",
			// "Content-Disposition":          `attachment; filename="images.zip"`,
			// "Access-Control-Allow-Origin":  "*",
			// "Access-Control-Allow-Headers": "*",
			// "Content-Length":               fmt.Sprint(zipBuffer.Len()),
		},
		Body: base64.StdEncoding.EncodeToString(zipBuffer.Bytes()),
	}

	return resp
}
