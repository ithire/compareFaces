package main

import (
	face_persons "compareFaces/face-persons"
	"compareFaces/helpers"
	"net/http"
	"strconv"
)

func faceCompare(w http.ResponseWriter, r *http.Request) {
	jpeg := helpers.NewJpg("images/")
	convert := helpers.Convert{}

	face1, face1Reader := jpeg.GetJpg(r.FormValue("face1"))
	face2, face2Reader := jpeg.GetJpg(r.FormValue("face2"))

	image1Name, image1Increment, image1Type := jpeg.CreateImage(face1, face1Reader)
	image2Name, image2Increment, image2Type := jpeg.CreateImage(face2, face2Reader)

	if image1Type == "png" {
		convert.Convert([]string{image1Name}, "images", 1)
	}
	if image2Type == "png" {
		convert.Convert([]string{image2Name}, "images", 1)
	}

	defer jpeg.DeleteJpg(image1Increment + ".jpg")
	defer jpeg.DeleteJpg(image2Increment + ".jpg")

	var facePersons = face_persons.NewFacePerson("images", "empty.jpg", image1Increment+".jpg", image2Increment+".jpg")

	w.Write([]byte(strconv.FormatBool(facePersons.Run())))

}

func main() {
	http.HandleFunc("/face-compare", faceCompare)
	http.ListenAndServe(":7001", nil)
}
