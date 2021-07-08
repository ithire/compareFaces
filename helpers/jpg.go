package helpers

import "C"
import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"strconv"
)

type Jpg struct {
	dir string
}

func NewJpg(dir string) *Jpg {
	return &Jpg{
		dir: dir,
	}
}

func (j Jpg) CreateImage(data []byte, readerData io.Reader) (string, string) {
	var imageType = j.GetFormat(readerData)
	imageName := j.IncrementImageName()
	if imageType == "png" {
		imageName = j.IncrementPngImageName()
	}
	f, err := os.Create(j.dir + imageName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err2 := f.Write(data)
	if err2 != nil {
		return "", imageType
	}
	return imageName, imageType
}

func (j Jpg) DeleteJpg(name string) {
	err := os.Remove(j.dir + name)
	if err != nil {
		fmt.Println("don't delete")
	}
}

func (j Jpg) GetJpg(url string) ([]byte, io.ReadCloser) {
	res, err := http.Get(url)
	if err != nil {
		return nil, nil
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)

	return content, res.Body
}

// Guess image format from gif/jpeg/png/webp
func (j Jpg) GuessImageFormat(r io.Reader) (format string, err error) {
	_, format, err = image.DecodeConfig(r)
	return
}

// Guess image mime types from gif/jpeg/png/webp
func (j Jpg) GetFormat(r io.Reader) string {
	format, _ := j.GuessImageFormat(r)
	if format == "" {
		return ""
	}
	return mime.TypeByExtension(format)
}

func (j Jpg) PngToJpg(pngImageName string, newJpegImage string) string {
	pngImgFile, err := os.Open(j.dir + pngImageName)

	if err != nil {
		fmt.Println("PNG-file.png file not found!")
		return ""
	}

	defer pngImgFile.Close()
	defer j.DeleteJpg(pngImageName)

	// create image from PNG file
	imgSrc, err := png.Decode(pngImgFile)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	// create a new Image with the same dimension of PNG image
	newImg := image.NewRGBA(imgSrc.Bounds())

	// we will use white background to replace PNG's transparent background
	// you can change it to whichever color you want with
	// a new color.RGBA{} and use image.NewUniform(color.RGBA{<fill in color>}) function

	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// paste PNG image OVER to newImage
	draw.Draw(newImg, newImg.Bounds(), imgSrc, imgSrc.Bounds().Min, draw.Over)

	// create new out JPEG file
	jpgImgFile, err := os.Create(j.dir + newJpegImage)

	if err != nil {
		fmt.Println("Cannot create JPEG-file.jpg !")
		return ""
	}

	defer jpgImgFile.Close()

	var opt jpeg.Options
	opt.Quality = 80

	// convert newImage to JPEG encoded byte and save to jpgImgFile
	// with quality = 80
	err = jpeg.Encode(jpgImgFile, newImg, &opt)

	//err = jpeg.Encode(jpgImgFile, newImg, nil) -- use nil if ignore quality options

	if err != nil {
		fmt.Println(err)
		return ""
	}

	fmt.Println("Converted PNG file to JPEG file")

	return newJpegImage
}

func (j Jpg) IncrementImageName() string {
	return strconv.Itoa(rand.Int()) + ".jpg"
}

func (j Jpg) IncrementPngImageName() string {
	return strconv.Itoa(rand.Int()) + ".png"
}
