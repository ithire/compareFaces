package main

import (
	face_persons "compareFaces/face-persons"
	"compareFaces/helpers"
	"net/http"
	"strconv"
)

func faceCompare(w http.ResponseWriter, r *http.Request) {
	jpeg := helpers.NewJpg("images/")

	face1 := jpeg.GetJpg(r.FormValue("face1"))
	face2 := jpeg.GetJpg(r.FormValue("face2"))

	image1Name := jpeg.IncrementImageName()
	image2Name := jpeg.IncrementImageName()

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
	http.ListenAndServe(":7001", nil)
}
