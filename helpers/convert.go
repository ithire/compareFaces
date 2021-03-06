package helpers

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/chai2010/webp"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

type FileType int

const (
	PNG FileType = iota
	JPG
	GIF
	WEBP
	BMP
	TIFF
	ERR
)

type Convert struct {
}

func (c Convert) GetFileType(input string) FileType {
	switch input {
	case "jpg":
		fallthrough
	case "jpeg":
		return JPG
	case "png":
		return PNG
	case "gif":
		return GIF
	case "bmp":
		return BMP
	case "webp":
		return WEBP
	case "tiff":
		return TIFF
	default:
		return ERR
	}
}

func (c Convert) GetFileExtension(input FileType) string {
	switch input {
	case JPG:
		return "jpg"
	case PNG:
		return "png"
	case GIF:
		return "gif"
	case BMP:
		return "bmp"
	case WEBP:
		return "webp"
	case TIFF:
		return "tiff"
	default:
		return ""
	}
}

func (c Convert) Convert(files []string, outputDir string, fileType FileType) {
	var wg sync.WaitGroup
	for _, currPath := range files {
		wg.Add(1)
		go c.ConvertFile(&wg, currPath, outputDir, fileType)
	}

	wg.Wait()
}

func (c Convert) ConvertFile(wg *sync.WaitGroup, currPath string, outputDir string, fileType FileType) {
	// call done when finished
	defer wg.Done()

	ext := strings.ToLower(filepath.Ext(currPath))
	newExt := c.GetFileExtension(fileType)

	_, filename := filepath.Split(currPath)
	filenameNoExt := filename[0 : len(filename)-len(ext)]
	newFileName := filenameNoExt + "." + newExt
	newFilePath := outputDir + "/" + newFileName

	// validate file type
	startType := c.GetFileType(ext[1:])
	if startType == ERR {
		fmt.Println(errors.New("input file type not valid"))
	}

	// open files
	file, err := os.Open(outputDir + "/" + currPath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	defer os.Remove(outputDir + "/" + currPath)

	outFile := c.OpenOrCreate(newFilePath)
	defer outFile.Close()

	// decode
	imageData, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(errors.New("error decoding image"))
	}

	// encode in new type
	switch fileType {
	case JPG:
		err := jpeg.Encode(outFile, imageData, nil)
		if err != nil {
			fmt.Println(errors.New("error converting to jpeg"))
		}
	case PNG:
		err := png.Encode(outFile, imageData)
		if err != nil {
			fmt.Println(err)
		}
	case WEBP:
		if err := webp.Encode(outFile, imageData, nil); err != nil {
			fmt.Println(errors.New("error converting to webp"))
		}
	case GIF:
		err := gif.Encode(outFile, imageData, nil)
		if err != nil {
			fmt.Println(errors.New("error converting to gif"))
		}
	case BMP:
		err := bmp.Encode(outFile, imageData)
		if err != nil {
			fmt.Println(errors.New("error converting to bmp"))
		}
	case TIFF:
		err := tiff.Encode(outFile, imageData, nil)
		if err != nil {
			fmt.Println(errors.New("error converting to tiff"))
		}
	}
}

func (c Convert) OpenOrCreate(filename string) *os.File {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println(errors.New("error creating output file"))
		}
		return file
	} else {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(errors.New("error opening output file"))
		}
		return file
	}
}
