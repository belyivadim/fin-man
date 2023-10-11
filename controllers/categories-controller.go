package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/belyivadim/fin-man/database"
	"github.com/belyivadim/fin-man/models"
	"github.com/gin-gonic/gin"
)

func getCategoryService(c *gin.Context) database.CategoryService {
  categoryService, ok := c.Copy().MustGet("CategoryService").(database.CategoryService)
  if !ok {
    c.JSON(http.StatusInternalServerError, ApiError{Error: "Internal Server Error"})
    c.Abort()
    log.Fatal("[category-controller]: Failed to get CategoryService.")
  }

  return categoryService
}

func AddCategory(c *gin.Context) {
  var category models.Category
  err := c.ShouldBindJSON(&category)
  if err != nil {
    log.Println(err)
    c.JSON(http.StatusBadRequest, ApiError{Error: "Invalid inputs"})
    c.Abort()
    return
  }

  categoryService := getCategoryService(c)

  err = categoryService.CreateCategoryRecord(&category)
  if err != nil {
    c.JSON(http.StatusInternalServerError, ApiError{Error: "Error creating category"})
    c.Abort()
    return
  }

  c.JSON(http.StatusCreated, &category)
}

func RemoveCategory(c *gin.Context) {
  id, err := strconv.ParseUint(c.Param("id"), 0, 32)
  if err != nil {
    c.JSON(http.StatusBadRequest, ApiError{Error: "Invalid path param `id`"})
    c.Abort()
    return
  }

  categoryService := getCategoryService(c)

  err = categoryService.DeleteCategoryRecordById(uint(id))
  if err != nil {
    c.JSON(http.StatusInternalServerError, ApiError{Error: "Error deleting category"})
    c.Abort()
    return
  }

  c.JSON(http.StatusNoContent, gin.H{})
}

func UpdateCategory(c *gin.Context) {
  id, err := strconv.ParseUint(c.Param("id"), 0, 32)
  if err != nil {
    c.JSON(http.StatusBadRequest, ApiError{Error: "Invalid path param `id`"})
    c.Abort()
    return
  }

  var category models.Category
  err = c.ShouldBindJSON(&category)
  if err != nil {
    log.Println(err)
    c.JSON(http.StatusBadRequest, ApiError{Error: "Invalid inputs"})
    c.Abort()
    return
  }
  category.ID = uint(id)

  categoryService := getCategoryService(c)

  err = categoryService.UpdateCategoryRecord(&category)
  if err != nil {
    if errors.Is(err, &database.CategoryNotFoundError{}) {
      c.JSON(http.StatusNotFound, ApiError{Error: "Category with `id` is not found"})
    } else {
      c.JSON(http.StatusInternalServerError, ApiError{Error: "Error updating category"})
    }
    c.Abort()
    return
  }

  c.JSON(http.StatusNoContent, gin.H{})
}

func GetAllCategories(c *gin.Context) {
  categoryService := getCategoryService(c)
  categories, err := categoryService.GetAllCategories()
  if err != nil {
    c.JSON(http.StatusInternalServerError, ApiError{Error: "Error fetching categories"})
    c.Abort()
    return
  }

  c.JSON(http.StatusOK, categories)
}

func GetCategoryById(c *gin.Context) {
  id, err := strconv.ParseUint(c.Param("id"), 0, 32)
  if err != nil {
    c.JSON(http.StatusBadRequest, ApiError{Error: "Invalid path param `id`"})
    c.Abort()
    return
  }

  categoryService := getCategoryService(c)
  category, err := categoryService.GetCategoryById(uint(id))
  if err != nil {
    c.JSON(http.StatusInternalServerError, ApiError{Error: "Error fetching category"})
    c.Abort()
    return
  }

  c.JSON(http.StatusOK, category)
}
