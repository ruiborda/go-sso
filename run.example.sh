#!/bin/bash

# Script para configurar las variables de entorno necesarias
# para el desarrollo local del servicio SSO con Google

# Credenciales de Google OAuth2
export GOOGLE_CLIENT_ID="742682450289-...apps.googleusercontent.com"
export GOOGLE_CLIENT_SECRET="GOCSPX-.."
export GOOGLE_REDIRECT_URL="http://localhost:8080/callback"

# Ejecuta la aplicación
go run main.go