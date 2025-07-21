// internal/handlers/category.go
package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    
    "github.com/gorilla/mux"
    "reflection-app/internal/models"
    "reflection-app/internal/service"
)

type CategoryHandler struct {
    service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
    return &CategoryHandler{service: service}
}

func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
    // Get user ID from query parameter
    userIDStr := r.URL.Query().Get("user_id")
    if userIDStr == "" {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    categories, err := h.service.GetCategoriesByUser(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) GetSubCategories(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    categoryID, err := strconv.Atoi(vars["categoryId"])
    if err != nil {
        http.Error(w, "Invalid category ID", http.StatusBadRequest)
        return
    }
    
    subcategories, err := h.service.GetSubCategoriesByCategory(categoryID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(subcategories)
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
    var category models.Category
    
    if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    if err := h.service.CreateCategory(&category); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) CreateSubCategory(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    categoryID, err := strconv.Atoi(vars["categoryId"])
    if err != nil {
        http.Error(w, "Invalid category ID", http.StatusBadRequest)
        return
    }
    
    var subcategory models.SubCategory
    if err := json.NewDecoder(r.Body).Decode(&subcategory); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    subcategory.CategoryID = categoryID
    
    if err := h.service.CreateSubCategory(&subcategory); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(subcategory)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    categoryID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid category ID", http.StatusBadRequest)
        return
    }
    
    var category models.Category
    if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    category.ID = categoryID
    
    if err := h.service.UpdateCategory(&category); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    categoryID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid category ID", http.StatusBadRequest)
        return
    }
    
    if err := h.service.DeleteCategory(categoryID); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

func (h *CategoryHandler) UpdateSubCategory(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    subcategoryID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid subcategory ID", http.StatusBadRequest)
        return
    }
    
    var subcategory models.SubCategory
    if err := json.NewDecoder(r.Body).Decode(&subcategory); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    subcategory.ID = subcategoryID
    
    if err := h.service.UpdateSubCategory(&subcategory); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(subcategory)
}

func (h *CategoryHandler) DeleteSubCategory(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    subcategoryID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid subcategory ID", http.StatusBadRequest)
        return
    }
    
    if err := h.service.DeleteSubCategory(subcategoryID); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}