package applications

import (
	"PostJson/db"
	"PostJson/structures"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// swagger:model Applications
// type test struct {
// 	U_ID string `json:"user_id" validate:"uuid"`
// }

type myJSON struct {
	Users          []string
	RequiredSkills []string
}

func Shortlist(w http.ResponseWriter, r *http.Request) {

	var t []test
	var req_skills []structures.RequiredSkills
	var app structures.Applications
	var apps []structures.Applications
	var usrs []string
	var reqSkills []string
	// dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	// json.Unmarshal(dataFromWeb, &dataToCompare)

	app.Job_ID = strings.ReplaceAll(r.URL.Query().Get("job_id"), `"`, "")
	if len(app.Job_ID) > 0 {
		// populating add for validation
		app.U_ID = app.Job_ID
		// validate := validator.New()
		// err := validate.Struct(app)
		// if err != nil {
		// 	w.WriteHeader(400)
		// 	fmt.Fprintf(w, "Incorrect input!!")
		// 	return
		// }

		// selecting U_ID only to be shortlisted
		result := db.Conn.Model(&apps).Where("Job_ID = ?", app.Job_ID).Scan(&t)
		if result.Error != nil {
			w.WriteHeader(400)
			fmt.Fprintln(w, result.Error)
			return
		}

		// creating json array of users
		for _, v := range t {
			usrs = append(usrs, v.U_ID)
		}

		result2 := db.Conn.Model(&req_skills).Where("Job_ID = ?", app.Job_ID).Find(&req_skills)
		if result2.Error != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return")
			return
		}

		// creating json array of skills
		for _, v := range req_skills {
			reqSkills = append(reqSkills, v.Skill)
		}

		// combining json arrays
		jsondat := &myJSON{Users: usrs, RequiredSkills: reqSkills}
		encjson, _ := json.Marshal(jsondat)

		if len(t) == 0 {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return")
			return
		}

		data, _ := json.Marshal(t)
		json.Unmarshal(data, &dataToCompare)

		// change url with python's url later. It is Path parameter after url
		posturl := "http://34.93.204.130:8000/" + app.Job_ID

		// concurently sending request to python
		go SendRequest(posturl, encjson)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Shortlisting started!!")
	}
}

func SendRequest(url string, data []byte) {
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
}
