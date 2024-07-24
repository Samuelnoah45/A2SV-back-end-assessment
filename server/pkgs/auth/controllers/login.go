package auth_controllers

// imports
import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	graphqlClient "server/clients/graphql"
	authModel "server/pkgs/auth/models"
	"server/utilService"
)

// login controller
func Login(ctx *gin.Context) {
	var validate = validator.New()

	var input struct {
		Password string `json:"password" validate:"required,min=6"`
		Email    string `json:"email" validate:"required,email,min=6"`
	}

	// bind data from request body
	if dataBindError := ctx.ShouldBindJSON(&input); dataBindError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": dataBindError.Error()})
		return
	}

	// validate the input data
	if err := validate.Struct(input); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errorMsg string
		for _, e := range validationErrors {
			errorMsg = fmt.Sprintf("Field %s validation failed on the %s tag", e.Field(), e.Tag())
			break
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": errorMsg})
		return
	}
	var query struct {
		Users []struct {
			ID         string `json:"id"`
			First_name string `json:"first_name"`
			Last_name  string `json:"last_name"`
			Email      string `json:"email"`
			Password   string `json:"password"`
			User_roles []struct {
				Role_name string `json:"role_name"`
			} `json:"user_roles"`
		} `graphql:"users(where: {email: {_eq: $email}})"`
	}

	variables := map[string]interface{}{
		"email": input.Email,
	}
	err := graphqlClient.AnonymousClient().Query(context.Background(), &query, variables)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	// check if user exists
	if len(query.Users) == 0 {
		ctx.JSON(400, gin.H{"message": "There is no account with email address " + input.Email})
		return
	}

	// check if the account is verified

	user_roles := []string{}
	for _, role := range query.Users[0].User_roles {
		editedRole := strings.Replace(role.Role_name, "_", "-", -1)
		user_roles = append(user_roles, editedRole)
	}

	if len(query.Users) > 0 && utilService.ComparePasswords(query.Users[0].Password, input.Password) {
		var user authModel.User
		user.ID = query.Users[0].ID
		user.Email = query.Users[0].Email
		user.FirstName = query.Users[0].First_name
		user.LastName = query.Users[0].Last_name
		user.UserRoles = user_roles
		sendTokenAndUserData(ctx, user)
		return
	} else {
		ctx.JSON(400, gin.H{"message": "Invalid credentials"})
		return
	}
}
