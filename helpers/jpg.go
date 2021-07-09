package helpers

import "C"
import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Jpg struct {
	dir string
}

var magicTable = map[string]string{
	"\xff\xd8\xff":      "jpg",
	"\x89PNG\r\n\x1a\n": "png",
	"GIF87a":            "gif",
	"GIF89a":            "gif",
}

func (j Jpg) GetFormat(incipit []byte) string {
	incipitStr := string(incipit)
	for magic, mime := range magicTable {
		if strings.HasPrefix(incipitStr, magic) {
			return mime
		}
	}
	return ""
}

func NewJpg(dir string) *Jpg {
	return &Jpg{
		dir: dir,
	}
}

func (j Jpg) CreateImage(data []byte, readerData io.Reader) (string, string, string) {
	var imageType = j.GetFormat(data)
	var imageIncrement = j.IncrementImageName()
	imageName := imageIncrement + "." + imageType

	f, err := os.Create(j.dir + imageName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err2 := f.Write(data)
	if err2 != nil {
		return "", "", imageType
	}
	return imageName, imageIncrement, imageType
}

func (j Jpg) DeleteJpg(name string) {
	err := os.Remove(j.dir + name)
	if err != nil {
		fmt.Println("don't delete")
	}
}

func (j Jpg) GetJpg(url string) ([]byte, io.Reader) {
	res, err := http.Get(url)
	if err != nil {
		return nil, nil
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)

	return content, res.Body
}

func (j Jpg) IncrementImageName() string {
	return strconv.Itoa(rand.Int())
}
