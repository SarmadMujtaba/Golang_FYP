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
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

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
	tempFile, err := ioutil.TempFile("/home/sarmad/Desktop/FYP_Resumes", "upload-*.pdf")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
	w.WriteHeader(200)
}
