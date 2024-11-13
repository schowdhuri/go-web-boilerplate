package admin

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"viabl.ventures/gossr/internal/config"
	"viabl.ventures/gossr/internal/templates"
	"viabl.ventures/gossr/internal/utils"
)

type AdminRouter struct {
	conf          *config.EnvVars
	renderer      *templates.Renderer
	emailService  *utils.EmailService
	signinService *AuthService
}

func NewAdminRouter(conf *config.EnvVars, renderer *templates.Renderer, emailService *utils.EmailService, signinService *AuthService) *AdminRouter {
	return &AdminRouter{conf, renderer, emailService, signinService}
}

func (router *AdminRouter) GetRoutes(r chi.Router) {
	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(CreateAdminSessionMiddleware(router.signinService))
		r.Get("/", router.adminHomeView)
		r.Post("/signout", router.signoutHandler)
	})
	// Public routes
	r.Get("/signin", router.signinInitialViewHandler)
	r.Post("/signin", router.adminGenerateLoginCode)
	r.Get("/signin/verify", router.adminValidateLoginCode)

}

// render the admin home page
func (router *AdminRouter) adminHomeView(w http.ResponseWriter, _ *http.Request) {
	data := map[string]interface{}{
		"Title": "Admin Dashboard",
	}
	router.renderer.RenderTemplate(w, "admin_home.html", data)
}

// Step 1: render the admin signin page to capture the email
func (router *AdminRouter) signinInitialViewHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Admin Signin :: Get Code",
	}
	router.renderer.RenderTemplate(w, "signin.html", data)
}

// Step 2: render the admin signin page to capture the OTP and email a signin link
func (router *AdminRouter) adminGenerateLoginCode(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Admin Signin :: Verify Code",
		"Email": nil,
		"Error": nil,
	}

	r.ParseForm()
	email := r.PostFormValue("email")
	// TODO: validate recaptcha

	// Generate a one-time code (OTP)
	loginCode, err := router.signinService.GenerateCode(email)
	if err != nil {
		data["Error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		router.renderer.RenderTemplate(w, "signin.html", data)
		return
	}

	// Send the OTP to the user's email
	err = router.emailService.SendMail(email, "Sign In Link",
		fmt.Sprintf(`Sign in using this link:
			<p>
			<a href='%s/admin/signin/verify?email=%s&code=%s'>Sign In</a>
			</p>`, router.conf.PublicUrl, base64.StdEncoding.EncodeToString([]byte(email)), loginCode.Code))
	if err != nil {
		data["Error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		router.renderer.RenderTemplate(w, "signin.html", data)
		return
	}

	// Respond with a success message or appropriate error
	data["Email"] = email
	router.renderer.RenderTemplate(w, "signin.html", data)
}

// Step 3: validate the OTP and sign the user in
func (router *AdminRouter) adminValidateLoginCode(w http.ResponseWriter, r *http.Request) {
	// Parse the request URL to get the user's email and code
	query := r.URL.Query()

	email, err := base64.StdEncoding.DecodeString(query.Get("email"))
	if err != nil {
		http.Redirect(w, r, "/admin/signin?error=Invalid link", http.StatusSeeOther)
		return
	}
	code := query.Get("code")
	if code == "" {
		http.Redirect(w, r, "/admin/signin?error=Invalid link", http.StatusSeeOther)
		return
	}

	// Validate the code
	_, err = router.signinService.ValidateAndUseCode(string(email), string(code))
	if err != nil {
		http.Redirect(w, r, "/admin/signin?error=Invalid code", http.StatusSeeOther)
		return
	}
	cookie, err := router.signinService.CreateSessionCookie(string(email))
	if err != nil {
		// Redirect to the signin page with an error message
		http.Redirect(w, r, "/admin/signin?error=Invalid code", http.StatusSeeOther)
		return
	}

	// set cookie
	http.SetCookie(w, cookie)
	// Redirect to the admin home page
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// Sign out the user by deleting the session cookie
func (router *AdminRouter) signoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/admin/signin?error=unauthorized", http.StatusSeeOther)
		return
	}
	err = router.signinService.DeleteSession(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/admin/signin?error=unauthorized", http.StatusSeeOther)
		return
	}
	// Delete the session cookie by setting it to expire 1h ago
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/admin/signin?info=you have been signed out", http.StatusSeeOther)
}
