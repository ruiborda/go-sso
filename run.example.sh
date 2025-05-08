#!/bin/bash

# Script para configurar las variables de entorno necesarias
# para el desarrollo local del servicio SSO con Google y GitHub

# Credenciales de Google OAuth2
export GOOGLE_CLIENT_ID="11111-...apps.googleusercontent.com"
export GOOGLE_CLIENT_SECRET="aaaaaaaaaaaaa-.."
export GOOGLE_REDIRECT_URL="http://localhost:8080/callback"

# Credenciales de GitHub OAuth
export GITHUB_CLIENT_ID="xxxxxxxxxx"
export GITHUB_CLIENT_SECRET="xxxxxxxxxxxxx"
export GITHUB_REDIRECT_URL="http://localhost:8080/github.callback"

# Ejecuta la aplicación
go run main.go