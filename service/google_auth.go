package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2API "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

// GoogleAuthService maneja la autenticación con Google OAuth2
type GoogleAuthService struct {
	OAuthConfig     *oauth2.Config
	OAuthStateString string
}

// UserInfo contiene los campos que mostraremos en la plantilla
type UserInfo struct {
	Email         string
	ID            string
	Name          string
	GivenName     string
	FamilyName    string
	Picture       string
	PictureBase64 string
	Locale        string
	VerifiedEmail string
}

// NewGoogleAuthService crea una nueva instancia del servicio de autenticación de Google
func NewGoogleAuthService() *GoogleAuthService {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/callback"
	}

	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleAuthService{
		OAuthConfig:     oauthConfig,
		OAuthStateString: "random", // Puedes implementar una generación aleatoria segura aquí
	}
}

// GetAuthURL genera la URL para iniciar el flujo OAuth2
func (s *GoogleAuthService) GetAuthURL() string {
	return s.OAuthConfig.AuthCodeURL(s.OAuthStateString)
}

// HandleCallback procesa el callback de Google y obtiene la información del usuario
func (s *GoogleAuthService) HandleCallback(r *http.Request) (*UserInfo, error) {
	if r.FormValue("state") != s.OAuthStateString {
		return nil, fmt.Errorf("estado inválido")
	}

	token, err := s.OAuthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener el token: %w", err)
	}

	client := s.OAuthConfig.Client(context.Background(), token)
	oauth2Service, err := oauth2API.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error al crear el servicio OAuth2: %w", err)
	}

	userinfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("error al obtener la información del usuario: %w", err)
	}

	user := &UserInfo{
		Email:         userinfo.Email,
		ID:            userinfo.Id,
		Name:          userinfo.Name,
		GivenName:     userinfo.GivenName,
		FamilyName:    userinfo.FamilyName,
		Picture:       userinfo.Picture,
		Locale:        userinfo.Locale,
		VerifiedEmail: fmt.Sprintf("%t", userinfo.VerifiedEmail),
	}

	// Descargar y convertir la imagen de perfil a base64
	if userinfo.Picture != "" {
		if data, ct, err := downloadImage(userinfo.Picture); err != nil {
			user.PictureBase64 = user.Picture // fallback a URL
		} else {
			user.PictureBase64 = imageToBase64DataURI(data, ct)
		}
	} else {
		user.PictureBase64 = user.Picture
	}

	return user, nil
}

// downloadImage descarga la imagen desde una URL y devuelve los bytes y el contentType
func downloadImage(url string) ([]byte, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("error al descargar la imagen, código de estado: %d", resp.StatusCode)
	}

	// Leer content type real
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg" // fallback
	}

	data, err := io.ReadAll(resp.Body)
	return data, contentType, err
}

// imageToBase64DataURI convierte una imagen a formato data URI con base64
func imageToBase64DataURI(imageData []byte, contentType string) string {
	base64Data := base64.StdEncoding.EncodeToString(imageData)
	return fmt.Sprintf("data:%s;base64,%s", contentType, base64Data)
}