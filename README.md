# Go-SSO: Single Sign-On Service in Go

This is a simple Single Sign-On (SSO) service built in Go that supports authentication through Google and GitHub OAuth2 providers.

## Features

- Authentication with Google OAuth2
- Authentication with GitHub OAuth2
- Simple responsive web interface using Tailwind CSS
- Display of user profile information after successful login
- Profile picture conversion to base64 for embedding

## Prerequisites

- Go 1.16 or newer
- OAuth2 credentials for Google and GitHub (client ID and client secret)

## Installation

1. Clone this repository:
   ```bash
   git clone <repository-url>
   cd go-sso
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

## Configuration

You need to set up OAuth2 credentials for both Google and GitHub:

### Google OAuth2 Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Navigate to "APIs & Services" > "Credentials"
4. Click "Create Credentials" > "OAuth Client ID"
5. Configure the consent screen
6. Set the application type to "Web Application"
7. Add `http://localhost:8080/callback` to the authorized redirect URIs
8. Copy the generated client ID and client secret

### GitHub OAuth2 Setup

1. Go to [GitHub Developer Settings](https://github.com/settings/developers)
2. Click on "New OAuth App"
3. Fill in the application details
4. Set the Homepage URL to `http://localhost:8080`
5. Set the Authorization callback URL to `http://localhost:8080/github.callback`
6. Register the application
7. Copy the generated client ID and client secret

## Running the Application

### Using the Provided Script

The easiest way to run the application is using the provided `run.sh` script, which sets up the required environment variables:

```bash
# Make the script executable
chmod +x run.sh

# Run the application
./run.sh
```

### Running Manually

If you prefer to set up the environment variables manually:

1. Configure the required environment variables:
   ```bash
   export GOOGLE_CLIENT_ID="your-google-client-id"
   export GOOGLE_CLIENT_SECRET="your-google-client-secret"
   export GOOGLE_REDIRECT_URL="http://localhost:8080/callback"
   
   export GITHUB_CLIENT_ID="your-github-client-id"
   export GITHUB_CLIENT_SECRET="your-github-client-secret"
   export GITHUB_REDIRECT_URL="http://localhost:8080/github.callback"
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

## Usage

1. Open your browser and navigate to `http://localhost:8080`
2. You'll see a login page with options to log in with Google or GitHub
3. Click on your preferred login method
4. Authorize the application to access your profile data
5. After successful authentication, you'll be redirected back to the application where your profile information will be displayed

## Project Structure

- `main.go`: The main application file that sets up the HTTP server and routes
- `service/google_auth.go`: Implementation of the Google OAuth2 authentication service
- `service/github_auth.go`: Implementation of the GitHub OAuth2 authentication service
- `run.sh`: Script to easily run the application with required environment variables

## Security Notes

- This application uses a fixed state string for OAuth2 flow. In a production environment, you should generate a random state string for each authorization request.
- Sensitive information like client secrets should not be committed to version control. Consider using a secure vault or environment variables for production deployments.
- This is a demonstration project and may require additional security measures for production use.