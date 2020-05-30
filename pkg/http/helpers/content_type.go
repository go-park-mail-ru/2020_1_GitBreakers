package helpers

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/pkg/errors"
	"net/http"
)

var allowedImagesContentTypes = map[string]string{
	"image/bmp":                "bmp",
	"image/gif":                "gif",
	"image/png":                "png",
	"image/jpeg":               "jpeg",
	"image/jpg":                "jpg",
	"image/svg+xml":            "svg",
	"image/webp":               "webp",
	"image/tiff":               "tiff",
	"image/vnd.microsoft.icon": "ico",
}

func CheckImageFileContentType(fileContent []byte) (string, error) {
	contentType := http.DetectContentType(fileContent)

	extension, ok := allowedImagesContentTypes[contentType]
	if !ok {
		return "", errors.WithMessage(entityerrors.Invalid(),
			"this content type is not allowed")
	}

	return extension, nil
}
