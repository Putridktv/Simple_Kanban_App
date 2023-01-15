package api

import (
	"encoding/json"
	"fmt"
	"kanbanApp/entity"
	"kanbanApp/service"
	"log"
	"net/http"
	"strconv"
)

type TaskAPI interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) GetTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)
	if userID == "" {
		errorResponse := entity.ErrorResponse{Error: "invalid user id"}
		convErr, _ := json.Marshal(errorResponse)
		w.WriteHeader(400)
		w.Write(convErr)
		return
	}

	convUserID, _ := strconv.Atoi(userID)

	taskID := r.URL.Query().Get("task_id")
	convTaskID, _ := strconv.Atoi(taskID)

	if taskID == "" {
		check, err := t.taskService.GetTasks(r.Context(), convUserID)
		if err != nil {
			errorResponse := entity.ErrorResponse{Error: "error internal server"}
			convErr, _ := json.Marshal(errorResponse)
			w.WriteHeader(500)
			w.Write(convErr)
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(check)
	} else {
		check, err := t.taskService.GetTaskByID(r.Context(), convTaskID)
		if err != nil {
			errorResponse := entity.ErrorResponse{Error: "error internal server"}
			convErr, _ := json.Marshal(errorResponse)
			w.WriteHeader(500)
			w.Write(convErr)
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(check)
	}

}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	if task.Title == "" || task.Description == "" || task.CategoryID <= 0 {
		errorResponse := entity.ErrorResponse{Error: "invalid task request"}
		convErr, _ := json.Marshal(errorResponse)
		w.WriteHeader(400)
		w.Write(convErr)
		return
	}

	userID := r.Context().Value("id").(string)
	if userID == "" {
		errorResponse := entity.ErrorResponse{Error: "invalid user id"}
		convErr, _ := json.Marshal(errorResponse)
		w.WriteHeader(400)
		w.Write(convErr)
		return
	}

	convUserID, _ := strconv.Atoi(userID)
	checkTask := entity.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CategoryID:  task.CategoryID,
		UserID:      convUserID,
	}
	check, err := t.taskService.StoreTask(r.Context(), &checkTask)
	if err != nil {
		errorResponse := entity.ErrorResponse{Error: "error internal server"}
		convErr, _ := json.Marshal(errorResponse)
		w.WriteHeader(500)
		w.Write(convErr)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]interface{}{"user_id": check.UserID, "task_id": check.ID, "message": "success create new task"})

}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)
	taskID := r.URL.Query().Get("task_id")

	convUserID, _ := strconv.Atoi(userID)
	convTaskID, _ := strconv.Atoi(taskID)

	check := t.taskService.DeleteTask(r.Context(), convTaskID)
	if check != nil {
		errorResponse := entity.ErrorResponse{Error: "error internal server"}
		convErr, _ := json.Marshal(errorResponse)
		w.WriteHeader(500)
		w.Write(convErr)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{"user_id": convUserID, "task_id": convTaskID, "message": "success delete task"})

}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userID := r.Context().Value("id").(string)
	// fmt.Println("user id ", userID)
	if userID == "" {
		errorResponse := entity.ErrorResponse{Error: "invalid user id"}
		convErr, _ := json.Marshal(errorResponse)
		w.WriteHeader(400)
		w.Write(convErr)
		return
	}

	// fmt.Println("user id ", userID)
	convUserID, err := strconv.Atoi(userID)
	if err != nil {
		fmt.Println(err)
		return
	}

	checkTask := entity.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CategoryID:  task.CategoryID,
		UserID:      convUserID,
	}
	check, err := t.taskService.UpdateTask(r.Context(), &checkTask)
	if err != nil {
		errorResponse := entity.ErrorResponse{Error: "error internal server"}
		convErr, _ := json.Marshal(errorResponse)
		w.WriteHeader(500)
		w.Write(convErr)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{"user_id": check.UserID, "task_id": check.ID, "message": "success update task"})

}

func (t *taskAPI) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = entity.Task{
		ID:         task.ID,
		CategoryID: task.CategoryID,
		UserID:     int(idLogin),
	}

	_, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userId,
		"task_id": task.ID,
		"message": "success update task category",
	})
}
