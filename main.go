package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "os"
)

type msg struct {
    Status string `json:"status"`
    Desc  string    `json:"desc"`
}

func listHandler(w http.ResponseWriter, r *http.Request) {
    
	// Allowed periods
	periods := map[string]bool {
		"1h": true,
		"1d": true,
		"1mo": true,
		"1y": true,
	}

	// Get the querystring
	q := r.URL.Query()

	// Get period
	period := q.Get("period")
	if period == "" || !periods[period] {
		e := msg{ Status: "error", Desc: "Unsupported period"}
		r, _ := json.MarshalIndent(e, "", " ") 
		w.WriteHeader(400)
		w.Write(r)
		return
	}
	
	// Get timezone
	tz := q.Get("tz")
	loc, error := time.LoadLocation(tz)
	if error != nil {
		e := msg{ Status: "error", Desc: "Unsupported timezone"}
		resp, _ := json.MarshalIndent(e, "", " ") 
		w.WriteHeader(400)
		w.Write(resp)
		return
	}

	// Get times
	const tform = "20060102T150405Z"
	t1 := q.Get("t1")
	t2 := q.Get("t2")
	time1, _ := time.ParseInLocation(tform, t1, loc)
	time2, _ := time.ParseInLocation(tform, t2, loc)
	
	// Build response
	var response []string
	hour := 1 * time.Hour
	switch period {
		case "1h":
			for t := time1.Round(hour); t.After(time2) == false; t = t.Add(hour) {
				//fmt.Println(t.Format(tform))
				response = append(response, t.Format(tform))
			}
		case "1d":
			for t := time1.Round(hour); t.After(time2) == false; t = t.AddDate(0, 0, 1) {
				fmt.Println(t.Format(tform))
				response = append(response, t.Format(tform))
			}
		case "1mo":
			for t := time1.Round(hour); t.After(time2) == false; t = t.AddDate(0, 1, 0) {
				fmt.Println(t.Format(tform))
				response = append(response, t.Format(tform))
			}
		case "1y":
			for t := time1.Round(hour); t.After(time2) == false; t = t.AddDate(1, 0, 0) {
				fmt.Println(t.Format(tform))
				response = append(response, t.Format(tform))
			}
	}
	
	// Write response
	resp, _ := json.MarshalIndent(response, "", " ") 
	w.Write(resp)
        
}

func main() {
    http.HandleFunc("/ptlist", listHandler)
    if len(os.Args) > 1 {
		port := os.Args[1]
		http.ListenAndServe(":" + port, nil)
	} else {
		fmt.Println("No port specified, using default (8080)")
		http.ListenAndServe(":8080", nil)
	}
}
