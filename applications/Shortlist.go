package applications

import (
	"PostJson/db"
	"PostJson/structures"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

// swagger:model Applications
type test struct {
	U_ID string `json:"user_id" validate:"uuid"`
}

func Shortlist(w http.ResponseWriter, r *http.Request) {

	var t []test
	var app structures.Applications
	var apps []structures.Applications
	// dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	// json.Unmarshal(dataFromWeb, &dataToCompare)

	app.Job_ID = r.URL.Query().Get("job_id")
	if len(app.Job_ID) > 0 {
		// populating add for validation
		app.U_ID = app.Job_ID
		validate := validator.New()
		err := validate.Struct(app)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}

		// selecting U_ID only to be shortlisted
		result := db.Conn.Model(&apps).Where("Job_ID = ?", app.Job_ID).Scan(&t)

		if result.Error != nil {
			w.WriteHeader(400)
			fmt.Fprintln(w, result.Error)
			return
		}

		if len(t) == 0 {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return")
			return
		}

		data, _ := json.Marshal(t)
		json.Unmarshal(data, &dataToCompare)

		// change url with python's url later. It is Path parameter after url
		posturl := "http://127.0.0.1:8000/" + app.Job_ID

		fmt.Println(string(data))

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

		w.WriteHeader(http.StatusOK)

		json.Marshal(dataToCompare)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dataToCompare)

		fmt.Fprintln(w, "Shortlisting started!!")
	}
}
