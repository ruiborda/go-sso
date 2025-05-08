package main

import (
	"fmt"
	"go-sso/service"
	"html/template"
	"log"
	"net/http"
	"os"
)

// TemplateData envía opcionalmente User al template
type TemplateData struct {
	User *service.UserInfo
}

// Plantilla HTML con Tailwind CSS y renderizado condicional
const indexTemplate = `<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Login con Google (SSR)</title>
  <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100 flex items-center justify-center min-h-screen">
  <div class="bg-white p-8 rounded-lg shadow-md max-w-md w-full">
    {{if .User}}
      <h1 class="text-2xl font-bold mb-4 text-center text-indigo-600">Bienvenido, {{.User.Name}}</h1>
      <div class="flex justify-center mb-6">
        <img src="{{.User.PictureBase64}}" alt="Avatar de {{.User.Name}}" class="rounded-full h-24 w-24" />
      </div>
      <div class="space-y-2 text-gray-700">
        <p><span class="font-semibold">ID:</span> {{.User.ID}}</p>
        <p><span class="font-semibold">Email:</span> {{.User.Email}}</p>
        <p><span class="font-semibold">Nombre:</span> {{.User.Name}}</p>
        <p><span class="font-semibold">Given Name:</span> {{.User.GivenName}}</p>
        <p><span class="font-semibold">Family Name:</span> {{.User.FamilyName}}</p>
        <p><span class="font-semibold">Locale:</span> {{.User.Locale}}</p>
        <p><span class="font-semibold">Verified Email:</span> {{.User.VerifiedEmail}}</p>
      </div>
    {{else}}
      <h1 class="text-3xl font-bold mb-6 text-center text-indigo-600">Iniciar sesión</h1>
      <div class="flex flex-col space-y-4 items-center">
        <a href="/login" class="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded inline-flex items-center w-64 justify-center">
          <span>Iniciar sesión con Google</span>
        </a>
        <a href="/github.login" class="bg-gray-800 hover:bg-gray-700 text-white font-semibold py-2 px-4 rounded inline-flex items-center w-64 justify-center">
          <span>Iniciar sesión con GitHub</span>
        </a>
      </div>
    {{end}}
  </div>
</body>
</html>`

var (
	// Parseamos la plantilla una sola vez
	tmpl = template.Must(template.New("index").Parse(indexTemplate))
	
	// Servicios de autenticación
	googleAuth *service.GoogleAuthService
	githubAuth *service.GitHubAuthService
)

func main() {
	// Inicializar los servicios de autenticación
	googleAuth = service.NewGoogleAuthService()
	githubAuth = service.NewGitHubAuthService()
	
	// Configurar las rutas
	http.HandleFunc("/", serveIndexPage)
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/callback", handleGoogleCallback)
	http.HandleFunc("/github.login", handleGithubLogin)
	http.HandleFunc("/github.callback", handleGithubCallback)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Printf("Servidor corriendo en http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// serveIndexPage muestra la plantilla sin usuario
func serveIndexPage(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.Execute(w, TemplateData{User: nil}); err != nil {
		http.Error(w, "Error al renderizar la plantilla", http.StatusInternalServerError)
	}
}

// handleGoogleLogin redirige al flujo OAuth de Google
func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleAuth.GetAuthURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleGoogleCallback procesa el callback de Google y renderiza la plantilla con datos de usuario
func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	user, err := googleAuth.HandleCallback(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error en el callback: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Println("Información del usuario de Google:", user.Name, user.Email)
	
	// Renderizamos la plantilla con .User poblado
	if err := tmpl.Execute(w, TemplateData{User: user}); err != nil {
		http.Error(w, "Error al renderizar la plantilla con datos de usuario", http.StatusInternalServerError)
	}
}

// handleGithubLogin redirige al flujo OAuth de GitHub
func handleGithubLogin(w http.ResponseWriter, r *http.Request) {
	url := githubAuth.GetAuthURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleGithubCallback procesa el callback de GitHub y renderiza la plantilla con datos de usuario
func handleGithubCallback(w http.ResponseWriter, r *http.Request) {
	user, err := githubAuth.HandleCallback(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error en el callback de GitHub: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Println("Información del usuario de GitHub:", user.Name, user.Email)
	
	// Renderizamos la plantilla con .User poblado
	if err := tmpl.Execute(w, TemplateData{User: user}); err != nil {
		http.Error(w, "Error al renderizar la plantilla con datos de usuario", http.StatusInternalServerError)
	}
}