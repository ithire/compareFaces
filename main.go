package main

import (
	face_persons "compareFaces/face-persons"
	"compareFaces/helpers"
	"encoding/base64"
	"math/rand"
	"net/http"
	"strconv"
)

func faceCompare(w http.ResponseWriter, r *http.Request) {
	jpeg := helpers.NewJpg("images/")

	base64Face1 := r.FormValue("face1")
	base64Face2 := r.FormValue("face2")
	face1, _ := base64.StdEncoding.DecodeString(base64Face1)
	face2, _ := base64.StdEncoding.DecodeString(base64Face2)

	image1Name := strconv.Itoa(rand.Int()) + ".jpg"
	image2Name := strconv.Itoa(rand.Int()) + ".jpg"

	jpeg.CreateJpg(image1Name, face1)
	jpeg.CreateJpg(image2Name, face2)

	var facePersons = face_persons.NewFacePerson("images", "empty.jpg", image1Name, image2Name)

	w.Write([]byte(strconv.FormatBool(facePersons.Run())))

	defer jpeg.DeleteJpg(image1Name)
	defer jpeg.DeleteJpg(image2Name)

	return

}

func main() {
	http.HandleFunc("/face-compare", faceCompare)
	http.ListenAndServe(":8099", nil)
}
