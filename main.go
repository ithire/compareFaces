package main

import (
	face_persons "compareFaces/face-persons"
	"compareFaces/helpers"
	"net/http"
	"strconv"
)

func faceCompare(w http.ResponseWriter, r *http.Request) {
	jpeg := helpers.NewJpg("images/")

	face1, face1Reader := jpeg.GetJpg(r.FormValue("face1"))
	face2, face2Reader := jpeg.GetJpg(r.FormValue("face2"))

	image1Name, image1Type := jpeg.CreateImage(face1, face1Reader)
	image2Name, image2Type := jpeg.CreateImage(face2, face2Reader)

	if image1Type == "png" {
		image1Name = jpeg.PngToJpg(image1Name, jpeg.IncrementImageName())
	}
	if image2Type == "png" {
		image2Name = jpeg.PngToJpg(image2Name, jpeg.IncrementImageName())
	}

	defer jpeg.DeleteJpg(image1Name)
	defer jpeg.DeleteJpg(image2Name)

	var facePersons = face_persons.NewFacePerson("images", "empty.jpg", image1Name, image2Name)

	w.Write([]byte(strconv.FormatBool(facePersons.Run())))

}

func main() {
	http.HandleFunc("/face-compare", faceCompare)
	http.ListenAndServe(":7001", nil)
}
