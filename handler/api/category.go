package api

import (
	"encoding/json"
	"kanbanApp/entity"
	"kanbanApp/service"
	"log"
	"net/http"
	"strconv"
)

type CategoryAPI interface {
	GetCategory(w http.ResponseWriter, r *http.Request)
	CreateNewCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryWithTasks(w http.ResponseWriter, r *http.Request)
}

type categoryAPI struct {
	categoryService service.CategoryService
}

func NewCategoryAPI(categoryService service.CategoryService) *categoryAPI {
	return &categoryAPI{categoryService}
}

func (c *categoryAPI) GetCategory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)
	if userID == "" {
		errorResponse := entity.ErrorResponse{Error: "invalid user id"}
		convErr, _ := json.Marshal(errorResponse)
		w.WriteHeader(400)
		w.Write(convErr)
		return
	}

	convUserID, _ := strconv.Atoi(userID)

	check, err := c.categoryService.GetCategories(r.Context(), convUserID)
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

func (c *categoryAPI) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	var category entity.CategoryRequest

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	if category.Type == "" {
		errorResponse := entity.ErrorResponse{Error: "invalid category request"}
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
	checkCategory := entity.Category{
		UserID: convUserID,
		Type:   category.Type,
	}
	check, err := c.categoryService.StoreCategory(r.Context(), &checkCategory)
	if err != nil {
		errorResponse := entity.ErrorResponse{Error: "error internal server"}
		convErr, _ := json.Marshal(errorResponse)
		w.WriteHeader(500)
		w.Write(convErr)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]interface{}{"user_id": convUserID, "category_id": check.ID, "message": "success create new category"})

}

func (c *categoryAPI) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)
	categoryID := r.URL.Query().Get("category_id")

	convUserID, _ := strconv.Atoi(userID)
	convCategoryID, _ := strconv.Atoi(categoryID)
	check := c.categoryService.DeleteCategory(r.Context(), convCategoryID)
	if check != nil {
		errorResponse := entity.ErrorResponse{Error: "error internal server"}
		convErr, _ := json.Marshal(errorResponse)
		w.WriteHeader(500)
		w.Write(convErr)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{"user_id": convUserID, "category_id": convCategoryID, "message": "success delete category"})

}

func (c *categoryAPI) GetCategoryWithTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("get category task", err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categories, err := c.categoryService.GetCategoriesWithTasks(r.Context(), int(idLogin))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)

}
