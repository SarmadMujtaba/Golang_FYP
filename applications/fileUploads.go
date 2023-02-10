package applications

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func FileUpload(w http.ResponseWriter, r *http.Request) {
	// Parse our multipart form, 5 << 10 specifies a maximum. upload of 5 MB files.
	r.ParseMultipartForm(5 << 10)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename, the Header and the size of the file
	file, _, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
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
	str := "ABS-" + "*.pdf"
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

	// return that we have successfully uploaded our file!
	w.WriteHeader(200)
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
