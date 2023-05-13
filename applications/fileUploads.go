package applications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type test struct {
	U_ID string `json:"user_id" validate:"uuid"`
}

func FileUpload(w http.ResponseWriter, r *http.Request) {

	// var app structures.Applications

	user := strings.ReplaceAll(r.URL.Query().Get("user_id"), `"`, "")

	fmt.Println(user)
	if len(user) > 0 {
		// populating add for validation
		// app.Job_ID = app.U_ID
		// validate := validator.New()
		// err := validate.Struct(app)
		// if err != nil {
		// 	fmt.Println("invalid")
		// 	w.WriteHeader(400)
		// 	fmt.Fprintf(w, "Incorrect input!!")
		// 	return
		// }

		// Parse our multipart form, 5 << 10 specifies a maximum. upload of 5 MB files.
		r.ParseMultipartForm(5 << 10)
		// FormFile returns the first file for the given key `myFile`
		// it also returns the FileHeader so we can get the Filename, the Header and the size of the file
		file, _, err := r.FormFile("myFile")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			fmt.Println(r.Header.Get("Content-Type"))
			w.WriteHeader(400)
			return
		}
		defer file.Close()

		// // Only allow specific file uploading format (pdf, docx)
		// buff := make([]byte, 512)
		// _, err = file.Read(buff)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// filetype := http.DetectContentType(buff)
		// if filetype != "application/pdf" && filetype != "application/msword" {
		// 	w.WriteHeader(400)
		// 	fmt.Fprintf(w, "The provided file format is not allowed.")
		// 	return
		// }

		// Create a temporary file within our 'FYP_Resumes' directory that follows
		// a particular naming pattern (created directory manually)
		str := user + "_" + "*.pdf"

		tempFile, err := ioutil.TempFile("/app/Resumes", str)
		if err != nil {
			fmt.Println("msla 1: ", err)
		}
		defer tempFile.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("msla 2: ", err)
		}
		// write this byte array to our temporary file
		tempFile.Write(fileBytes)

		// sending file name to Python
		data, _ := json.Marshal(user)
		// change url with python's url later. It is Path parameter after url
		posturl := "http://host.docker.internal:8000/extract"

		r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(data))
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			fmt.Println(err)
			return
		}

		r.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(r)
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			fmt.Println(err)
			return
		}
		fmt.Println(resp.StatusCode)

		// return that we have successfully uploaded our file!
		w.WriteHeader(200)
		fmt.Fprintf(w, "Successfully Uploaded File\n")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("validation failed")
		fmt.Fprintf(w, "Please provide user ID with file.")
	}
}
