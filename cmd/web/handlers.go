package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"binme.haido.us/internal/models"
	"binme.haido.us/internal/validator"
	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	bins, err := app.bins.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Bins = bins

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) binView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	bin, err := app.bins.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Bin = bin

	app.render(w, r, http.StatusOK, "view.tmpl", data)

}

func (app *application) binCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = binCreateForm{
		Expires: 365,
	}
	app.render(w, r, http.StatusOK, "new.tmpl", data)
}

type binCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) binCreatePost(w http.ResponseWriter, r *http.Request) {
	var form binCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusOK, "new.tmpl", data)
		return
	}

	id, err := app.bins.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Bin successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/bins/%d", id), http.StatusSeeOther)
}

/**
* Auth-related
 */

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	templateData := app.newTemplateData(r)
	templateData.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl", templateData)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")

	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusOK, "signup.tmpl", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusOK, "signup.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Welcome to Bin Me! To start binning, please log in!")
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	templateData := app.newTemplateData(r)
	templateData.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl", templateData)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusOK, "login.tmpl", data)
		return
	}

	userId, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusOK, "login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), authenticatedUserIDKey.String(), userId)

	redirectUrl := app.sessionManager.PopString(r.Context(), redirectURLKey.String())
	if redirectUrl == "" {
		redirectUrl = "/bin/new"
	}
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), authenticatedUserIDKey.String())
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) accountView(w http.ResponseWriter, r *http.Request) {

	userID := app.sessionManager.GetInt(r.Context(), authenticatedUserIDKey.String())
	if userID == 0 {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	user, err := app.users.Get(userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	templateData := app.newTemplateData(r)
	templateData.User = user
	app.render(w, r, http.StatusOK, "account.tmpl", templateData)
}

type updatePasswordForm struct {
	CurrentPassword         string `form:"currentPassword"`
	NewPassword             string `form:"newPassword"`
	NewPasswordConfirmation string `form:"newPasswordConfirmation"`
	validator.Validator     `form:"-"`
}

func (app *application) accountUpdatePasswordView(w http.ResponseWriter, r *http.Request) {
	templateData := app.newTemplateData(r)
	templateData.Form = updatePasswordForm{}
	app.render(w, r, http.StatusOK, "password.tmpl", templateData)
}

func (app *application) accountUpdatePasswordPost(w http.ResponseWriter, r *http.Request) {
	var form updatePasswordForm
	err := app.decodePostForm(r, &form)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "This field cannot be blank")

	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")

	form.CheckField(validator.NotBlank(form.NewPasswordConfirmation), "newPasswordConfirmation", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.NewPasswordConfirmation, form.NewPassword), "newPasswordConfirmation", "Passwords do not match")

	if !form.Valid() {
		templateData := app.newTemplateData(r)
		templateData.Form = form
		app.render(w, r, http.StatusOK, "password.tmpl", templateData)
		return
	}

	userId := app.sessionManager.GetInt(r.Context(), authenticatedUserIDKey.String())
	err = app.users.PasswordUpdate(userId, form.CurrentPassword, form.NewPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("currentPassword", "Incorrect password")
			templateData := app.newTemplateData(r)
			templateData.Form = form
			app.render(w, r, http.StatusOK, "password.tmpl", templateData)
			return
		}
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your password has been updated!")
	http.Redirect(w, r, "/account", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
