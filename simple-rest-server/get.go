package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Search struct {
	Word       string    `json:"search"`
	TimeBefore time.Time `json:"before"`
}

func htmlParagraph(text *string, paragraph string) {
	*text += "\n<p>" + paragraph + "</p>"
}

func Index(respWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	htmlOut := `Todo List API.`
	htmlParagraph(&htmlOut, "/add")
	htmlParagraph(&htmlOut, "/delete")
	htmlParagraph(&htmlOut, "/search")
	htmlParagraph(&htmlOut, "/list")
	htmlParagraph(&htmlOut, "/checkoff")

	respWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	respWriter.Write([]byte(htmlOut))
}

func ListTask(respWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	output, _ := json.MarshalIndent(allTasks, "", "  ")
	fmt.Fprintln(respWriter, string(output))
}

func SearchTask(respWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var query Search
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(body, &query)
	if err != nil {
		fmt.Println(err)
	}

	var indicies []int
	fmt.Fprintln(respWriter, "Searching for "+query.Word)
	timezero := time.Time{}
	if query.Word != "" {
		indicies = append(indicies, searchByWord(query.Word)[:]...)
	} else if query.TimeBefore != timezero {
		indicies = append(indicies, searchByTime(query.TimeBefore)[:]...)
	} else {
		fmt.Println("Undefined search query")
	}

	if len(indicies) == 0 {
		fmt.Fprintln(respWriter, "No Results!!!!!")
	}
	for _, index := range indicies {
		accessTasks.Lock()
		prettyJson, err := json.MarshalIndent(allTasks[index], "", "  ")
		accessTasks.Unlock()
		if err != nil {
			log.Print(err)
			return
		}
		fmt.Fprintln(respWriter, string(prettyJson))
	}
}

func searchByWord(query string) (result []int) {
	accessTasks.Lock()
	for i, task := range allTasks {
		if strings.Contains(task.Task, query) {
			result = append(result, i)
		}
	}
	accessTasks.Unlock()
	return
}

func searchByTime(query time.Time) (result []int) {
	accessTasks.Lock()
	for i, task := range allTasks {
		if task.TimeAdded.Before(query) || task.TimeAdded == query {
			result = append(result, i)
		}
	}
	accessTasks.Unlock()
	return
}
