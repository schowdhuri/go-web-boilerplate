package admin

import (
	"net/http"
)

func CreateAdminSessionMiddleware(authService *AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Parse the session cookie
			cookie, err := r.Cookie("session")
			if err != nil || cookie.Value == "" {
				// Handle invalid or missing session
				http.Redirect(w, r, "/admin/signin?error=unauthorized", http.StatusSeeOther)
				return
			}

			renewed_cookie, err := authService.ValidateSession(cookie.Value)
			if err != nil {
				authService.DeleteSession(cookie.Value)
				// Handle invalid or expired session
				http.Redirect(w, r, "/admin/signin", http.StatusSeeOther)

				// If the session is valid, proceed to the next handler
				next.ServeHTTP(w, r)
				return
			}

			if renewed_cookie != nil {
				http.SetCookie(w, renewed_cookie)
			}

			next.ServeHTTP(w, r)
		})
	}
}
