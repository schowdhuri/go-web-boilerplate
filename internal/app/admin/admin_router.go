package admin

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"viabl.ventures/gossr/internal/templates"
	"viabl.ventures/gossr/internal/utils"
)

type AdminRouter struct {
	renderer      *templates.Renderer
	emailService  *utils.EmailService
	signinService *SigninService
}

func NewAdminRouter(renderer *templates.Renderer, emailService *utils.EmailService, signinService *SigninService) *AdminRouter {
	return &AdminRouter{renderer, emailService, signinService}
}

func (router *AdminRouter) GetRoutes(r chi.Router) {
	r.Get("/signin", router.signinInitialViewHandler)
	r.Post("/signin/validate", router.adminValidateLoginCode)
}

// Step 1: render the admin signin page to capture the email
func (router *AdminRouter) signinInitialViewHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":        "Admin Signin",
		"GetLoginCode": false,
	}
	router.renderer.RenderTemplate(w, "signin.html", data)
}

// Step 2: render the admin signin page to capture the OTP
func (router *AdminRouter) adminGenerateLoginCode(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the user's email
	r.ParseForm()
	email := r.PostFormValue("email")
	// TODO: validate recaptcha

	// Generate a one-time code (OTP)
	loginCode, err := router.signinService.GenerateCode(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send the OTP to the user's email
	err = router.emailService.SendMail(email, "One time code", "Your code is: <p><b>"+loginCode.Code+"</b><p>")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message or appropriate error
	data := map[string]interface{}{
		"Title": "Admin Signin",
		"Email": email,
	}
	router.renderer.RenderTemplate(w, "signin.html", data)
}

func (router *AdminRouter) adminValidateLoginCode(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the user's email and code
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := r.PostFormValue("email")
	code := r.PostFormValue("code")

	// Validate the code
	_, err = router.signinService.ValidateAndUseCode(email, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Respond with a success message or appropriate error
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}
