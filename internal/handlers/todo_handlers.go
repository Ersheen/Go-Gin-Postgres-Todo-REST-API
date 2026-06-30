package handlers

import (
	"net/http"
	"strconv"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type UpdatedTodo struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		useridinterface, exist := ctx.Get("user_id")
		if !exist {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user id not found"})
			return
		}
		userid := useridinterface.(string)

		var input CreateTodoInput

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo, err := repository.CreateTodo(pool, input.Title, input.Completed, userid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, todo)
	}
}

func GetAllTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		useridinterface, exist := ctx.Get("user_id")
		if !exist {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user id not found"})
			return
		}
		userid := useridinterface.(string)

		strstatus := ctx.Query("status")
		var boolstatus bool
		if strstatus != "" {
			var err error
			boolstatus, err = strconv.ParseBool(strstatus)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			todo, err := repository.GetTodoByQuery(pool, boolstatus, userid)
			if err != nil {
				if err.Error() == "no todo found" {
					ctx.JSON(http.StatusNotFound, gin.H{"error": "no todo found"})
					return
				} else {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

			}

			ctx.JSON(http.StatusOK, todo)
		} else {
			todo, err := repository.GetAllTodos(pool, boolstatus, userid)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			ctx.JSON(http.StatusOK, todo)
		}
	}
}

func GetTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		useridinterface, exist := ctx.Get("user_id")
		if !exist {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user id not found"})
			return
		}
		userid := useridinterface.(string)
		idstr := ctx.Param("id")

		idint, err := strconv.Atoi(idstr)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo, err := repository.GetTodo(pool, idint, userid)

		if err != nil {
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, todo)
	}
}

func UpdateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		useridinterface, exist := ctx.Get("user_id")
		if !exist {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user id not found"})
			return
		}
		userid := useridinterface.(string)

		idstr := ctx.Param("id")
		idInt, err := strconv.Atoi(idstr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updatedInput UpdatedTodo
		if err := ctx.ShouldBindJSON(&updatedInput); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if updatedInput.Title == nil && updatedInput.Completed == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
			return
		}

		todo, err := repository.UpdateTodo(pool, updatedInput.Title, updatedInput.Completed, idInt, userid)
		if err != nil {
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "todo does not exist"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, todo)
	}
}

func DeleteTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		useridinterface, exist := ctx.Get("user_id")
		if !exist {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user id not found"})
			return
		}
		userid := useridinterface.(string)

		idStr := ctx.Param("id")
		idInt, err := strconv.Atoi(idStr)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "not a integer"})
			return
		}

		err = repository.DeleteTodo(pool, idInt, userid)

		if err != nil {
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})

	}
}
