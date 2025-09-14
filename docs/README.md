# API Documentation

This directory contains the auto-generated OpenAPI/Swagger documentation for the Gin Authentication API.

## Files Generated

- `docs.go` - Go package containing the embedded documentation
- `swagger.json` - OpenAPI specification in JSON format
- `swagger.yaml` - OpenAPI specification in YAML format

## Accessing the Documentation

When running the application in development mode (`GIN_MODE=debug`), you can access the interactive Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

## Authentication Testing

To test protected endpoints in Swagger UI:

1. First, register a new user or login with existing credentials via the `/register` or `/login` endpoints
2. Copy the `access_token` from the response
3. Click the "Authorize" button in Swagger UI
4. Enter `Bearer <your_access_token>` in the authorization field
5. Now you can test protected endpoints

## API Endpoints

### Public Endpoints
- `POST /register` - Register new user
- `POST /login` - User login
- `POST /refresh` - Refresh access token  
- `POST /logout` - User logout
- `GET /health` - Health check

### Protected Endpoints (require JWT token)
- `GET /api/profile` - Get user profile
- `PUT /api/profile` - Update user profile
- `DELETE /api/profile` - Delete user profile

### Admin Endpoints (require admin role)
- `GET /api/admin/users` - List all users
- `POST /api/admin/users` - Create new user
- `PUT /api/admin/users/{id}` - Update user
- `DELETE /api/admin/users/{id}` - Delete user

## Regenerating Documentation

To regenerate the documentation after making changes to the API:

```bash
swag init -g cmd/server/main.go
```

This will update all files in the `docs/` directory.
