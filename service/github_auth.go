package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// GitHubAuthService maneja la autenticación con GitHub OAuth
type GitHubAuthService struct {
	OAuthConfig      *oauth2.Config
	OAuthStateString string
}

// GitHubUserInfo contiene los campos del usuario de GitHub
type GitHubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`

	// Campos adicionales para compatibilidad con la plantilla existente
	PictureBase64 string
	GivenName     string
	FamilyName    string
	Locale        string
	VerifiedEmail string
}

// NewGitHubAuthService crea una nueva instancia del servicio de autenticación de GitHub
func NewGitHubAuthService() *GitHubAuthService {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	redirectURL := os.Getenv("GITHUB_REDIRECT_URL")

	if redirectURL == "" {
		redirectURL = "http://localhost:8080/github.callback"
	}

	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"user:email", // Para obtener el email del usuario
		},
		Endpoint: github.Endpoint,
	}

	return &GitHubAuthService{
		OAuthConfig:      oauthConfig,
		OAuthStateString: "githubstate", // Puedes implementar una generación aleatoria segura aquí
	}
}

// GetAuthURL genera la URL para iniciar el flujo OAuth2
func (s *GitHubAuthService) GetAuthURL() string {
	return s.OAuthConfig.AuthCodeURL(s.OAuthStateString)
}

// HandleCallback procesa el callback de GitHub y obtiene la información del usuario
func (s *GitHubAuthService) HandleCallback(r *http.Request) (*UserInfo, error) {
	if r.FormValue("state") != s.OAuthStateString {
		return nil, fmt.Errorf("estado inválido")
	}

	token, err := s.OAuthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener el token: %w", err)
	}

	client := s.OAuthConfig.Client(context.Background(), token)

	// Obtener la información del usuario
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("error al obtener la información del usuario: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error al obtener la información del usuario, código de estado: %d", resp.StatusCode)
	}

	var githubUser GitHubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, fmt.Errorf("error al decodificar la respuesta: %w", err)
	}

	// Si el email es nulo, intentar obtenerlo desde la API de emails
	if githubUser.Email == "" {
		emails, err := getGitHubEmails(client)
		if err == nil && len(emails) > 0 {
			// Usar el email primario y verificado
			for _, email := range emails {
				if email.Primary && email.Verified {
					githubUser.Email = email.Email
					break
				}
			}
			// Si no hay email primario y verificado, usar el primero
			if githubUser.Email == "" && len(emails) > 0 {
				githubUser.Email = emails[0].Email
			}
		}
	}

	// Convertir a UserInfo para mantener compatibilidad con la plantilla existente
	user := &UserInfo{
		Email:         githubUser.Email,
		ID:            fmt.Sprintf("%d", githubUser.ID),
		Name:          githubUser.Name,
		GivenName:     githubUser.Login, // GitHub no tiene GivenName, usamos el login
		FamilyName:    "",               // GitHub no tiene FamilyName
		Picture:       githubUser.AvatarURL,
		Locale:        githubUser.Location, // GitHub no tiene locale, usamos location
		VerifiedEmail: "true",              // Asumimos que el email es verificado
	}

	// Descargar y convertir la imagen de perfil a base64
	if githubUser.AvatarURL != "" {
		if data, ct, err := downloadImage(githubUser.AvatarURL); err != nil {
			user.PictureBase64 = user.Picture // fallback a URL
		} else {
			user.PictureBase64 = imageToBase64DataURI(data, ct)
		}
	} else {
		user.PictureBase64 = user.Picture
	}

	return user, nil
}

// GitHubEmail representa la estructura de un email en la respuesta de la API de GitHub
type GitHubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

// getGitHubEmails obtiene los emails del usuario de GitHub
func getGitHubEmails(client *http.Client) ([]GitHubEmail, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error al obtener los emails, código de estado: %d", resp.StatusCode)
	}

	var emails []GitHubEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, err
	}

	return emails, nil
}
