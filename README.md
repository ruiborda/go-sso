# Go-SSO: Servicio de Inicio de Sesión Único (SSO) en Go Github y Google OAuth2

Este es un servicio simple de Inicio de Sesión Único (SSO) desarrollado en Go que permite la autenticación mediante los
proveedores OAuth2 de Google y GitHub.

## Características

- Autenticación con Google OAuth2
- Autenticación con GitHub OAuth2
- Interfaz web simple y responsiva usando Tailwind CSS
- Muestra de información del perfil del usuario tras el inicio de sesión exitoso
- Conversión de imagen de perfil a base64 para incrustarla

## Requisitos Previos

- docker y docker compose

Debes configurar las credenciales OAuth2 tanto para Google como para GitHub:

### Configuración de Google OAuth2

1. Ve a [Google Cloud Console](https://console.cloud.google.com/)
2. Crea un nuevo proyecto o selecciona uno existente
3. Navega a "APIs y servicios" > "Credenciales"
4. Haz clic en "Crear credenciales" > "ID de cliente de OAuth"
5. Configura la pantalla de consentimiento
6. Establece el tipo de aplicación como "Aplicación web"
7. Agrega `http://localhost:8080/callback` a las URI de redirección autorizadas
8. Copia el ID de cliente y el secreto de cliente generados

### Configuración de GitHub OAuth2

1. Ve a [GitHub Crear Nueva App](https://github.com/settings/apps/new)
2. Llena el nombre de la aplicación de GitHub
3. Llena la URL de la página principal
4. Llena la URL de redirección de autorización con `http://localhost:8080/github.callback`
5. Desactiva la casilla "Activo" del Webhook
6. En Permisos > selecciona "Cualquier cuenta"
7. Haz clic en "Crear aplicación de GitHub"

## Ejecución de la Aplicación

1. Clona este repositorio:
   ```bash
   git clone https://github.com/ruiborda/go-sso
   cd go-sso
   ```
2. Configurar Variables de entorno, crear un archivo `.env` en la raíz del proyecto
   ```bash
   cp example.env .env
   ```

2. Build the Docker image:
   ```bash
   docker compose build
   docker compose up -d
   ```

## Configuración

## Uso

1. Abre tu navegador y navega a `http://localhost:8080`
2. Verás una página de inicio de sesión con opciones para iniciar sesión con Google o GitHub
3. Haz clic en el método de inicio de sesión que prefieras
4. Autoriza a la aplicación para acceder a los datos de tu perfil
5. Tras una autenticación exitosa, serás redirigido nuevamente a la aplicación donde se mostrará la información de tu
   perfil

## Estructura del Proyecto

- `main.go`: Archivo principal de la aplicación que configura el servidor HTTP y las rutas
- `service/google_auth.go`: Implementación del servicio de autenticación con Google OAuth2
- `service/github_auth.go`: Implementación del servicio de autenticación con GitHub OAuth2
- `run.sh`: Script para ejecutar fácilmente la aplicación con las variables de entorno necesarias

