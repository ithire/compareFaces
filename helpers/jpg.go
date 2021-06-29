package helpers

import "C"
import (
	"fmt"
	"log"
	"os"
)

type Jpg struct {
	dir string
}

func NewJpg(dir string) *Jpg {
	return &Jpg{
		dir: dir,
	}
}

func (j Jpg) CreateJpg(name string, data []byte) bool {
	f, err := os.Create(j.dir + name)
	if err != nil {
		log.Fatal(err)
	}
	//defer
	defer f.Close()
	_, err2 := f.Write(data)
	if err2 != nil {
		return false
	}
	return true
}

func (j Jpg) DeleteJpg(name string) {
	err := os.Remove(j.dir + name)
	if err != nil {
		fmt.Println("don't delete")
	}
}
