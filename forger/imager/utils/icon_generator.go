package utils

import (
	"archive/zip"
	"fmt"
	"forger/imager/models"
	"image"
	"image/png"
	"path/filepath"
	"strconv"

	"github.com/disintegration/imaging"
)

const (
	IOS     string = "ios"
	Android string = "android"
)

func IOSmageResizer(zipWriter *zip.Writer, icon image.Image, resizeMetaList []models.ResizeMetaModel, platform string) error {

	for _, meta := range resizeMetaList {

		resized := imaging.Resize(icon, meta.Size, meta.Size, imaging.Lanczos)
		outputFileName := ""

		if platform == IOS {
			outputFileName = strconv.Itoa(meta.Size) + ".png"
		}

		if platform == Android {
			outputFileName = "ic_launcher.png"
		}

		fileInZip, err := zipWriter.Create(filepath.Join(platform, meta.DirName, outputFileName))
		if err != nil {
			return fmt.Errorf("error creating file in ZIP: %v", err)
		}

		err = png.Encode(fileInZip, resized)
		if err != nil {
			return fmt.Errorf("error encoding image to PNG: %v", err)
		}

		fmt.Println(zipWriter)
	}

	return nil
}
