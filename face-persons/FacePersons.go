package face_persons

import (
	"fmt"
	"github.com/Arturbox/go-face"
	"path/filepath"
)

type FacePerson struct {
	dataDir     string
	defaultFace string
	face1       string
	face2       string
	//esim face.Recognizer

}

func NewFacePerson(dataDir string, face1 string, face2 string, face3 string) *FacePerson {
	return &FacePerson{
		dataDir:     dataDir,
		defaultFace: face1,
		face1:       face2,
		face2:       face3,
	}
}

func (f FacePerson) Run() bool {
	rec, err := face.NewRecognizer(f.dataDir)
	if err != nil {
		//Cannot initialize recognizer
		return false
	}
	defer rec.Close()

	//Recognizer Initialized

	optimizedFace1 := f.Optimize(rec, f.face1)
	optimizedFace2 := f.Optimize(rec, f.face2)

	if optimizedFace1 == nil || optimizedFace2 == nil {
		return false
	}

	var faces = f.MergeFaces(f.Optimize(rec, f.defaultFace), optimizedFace1)

	var samples []face.Descriptor
	var avengers []int32
	for i, fc := range faces {
		samples = append(samples, fc.Descriptor)
		// Each face is unique on that image so goes to its own category.
		avengers = append(avengers, int32(i))
	}

	rec.SetSamples(samples, avengers)

	avengerID := rec.Classify(optimizedFace2.Descriptor)
	fmt.Println(avengerID)
	if avengerID < 1 {
		//Can't classify
		return false
	}

	return true
}

func (f FacePerson) Optimize(rec *face.Recognizer, face string) *face.Face {
	image := filepath.Join(f.dataDir, face)
	item, err := rec.RecognizeSingleFile(image)
	if err != nil {
		fmt.Println("Can't recognize", err)
		//log.Fatalf("Can't recognize: %v", err)
	}

	return item
}

func (f FacePerson) MergeFaces(face1 *face.Face, face2 *face.Face) (faces []*face.Face) {
	faces = append(faces, face1)
	faces = append(faces, face2)

	return faces
}
