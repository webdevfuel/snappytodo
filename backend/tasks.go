package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const (
	JSONFilePath = "./tasks.json"
)

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type CreateTaskBody struct {
	Name string `json:"name"`
}

type UpdateTaskBody struct {
	Name *string `json:"name"`
	Done *bool   `json:"done"`
}

func tasks(w http.ResponseWriter, _ *http.Request) {
	jsonFile, err := ioutil.ReadFile(JSONFilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error reading tasks json file %v", err)
		return
	}

	w.Write(jsonFile)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	body := CreateTaskBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error decoding body into struct %v", err)
		return
	}

	jsonFile, err := os.Open(JSONFilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error reading tasks jsonfile %v", err)
		return
	}

	tasks := []Task{}
	if err := json.NewDecoder(jsonFile).Decode(&tasks); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error decoding json file into struct %v", err)
		return
	}

	newTask := Task{
		Name: body.Name,
		Done: false,
		ID:   len(tasks) + 1,
	}

	tasks = append(tasks, newTask)

	j, err := json.Marshal(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling tasks %v", err)
		return
	}

	err = ioutil.WriteFile(JSONFilePath, j, 0755)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error writing tasks json file %v", err)
		return
	}

	j, err = json.Marshal(newTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling new task %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	body := UpdateTaskBody{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error decoding body into struct %v", err)
		return
	}

	jsonFile, err := os.Open(JSONFilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error reading tasks json file %v", err)
		return
	}

	tasks := []Task{}
	if err := json.NewDecoder(jsonFile).Decode(&tasks); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error decoding json file intro struct %v", err)
		return
	}

	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error converting taskID from string into integer %v", err)
		return
	}

	for i, task := range tasks {
		if task.ID == taskID {
			if body.Name != nil {
				task.Name = *body.Name
			}

			if body.Done != nil {
				task.Done = *body.Done
			}

			tasks[i] = task
		}
	}

	j, err := json.Marshal(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling tasks %v", err)
		return
	}

	err = ioutil.WriteFile(JSONFilePath, j, 0755)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error writing tasks json file %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	jsonFile, err := os.Open(JSONFilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error reading tasks json file %v", err)
		return
	}

	tasks := []Task{}
	if err := json.NewDecoder(jsonFile).Decode(&tasks); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error decoding json file into struct %v", err)
		return
	}

	prevLen := len(tasks)

	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error converting taskID from string into integer %v", err)
		return
	}

	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
		}
	}

	currLen := len(tasks)
	if currLen == prevLen {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	j, err := json.Marshal(tasks)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(JSONFilePath, j, 0755)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error writing tasks json file %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
