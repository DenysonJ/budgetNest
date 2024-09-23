package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
)

type UserController struct {
	//Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		//Inject services
	}
}

func (r *UserController) Show(ctx http.Context) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"Hello": "Goravel",
	})
}

func (r *UserController) Store(ctx http.Context) http.Response {
	user := models.User{
		Name:     ctx.Request().Input("name"),
		Email:    ctx.Request().Input("email"),
		Password: ctx.Request().Input("password"),
	}

	result := facades.Orm().Query().Create(&user)

	if result != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": result,
		})
	}

	return ctx.Response().Json(http.StatusCreated, http.Json{
		"message": "User created successfully",
	})
}
