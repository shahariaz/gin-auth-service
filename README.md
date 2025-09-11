# My Gin Boilerplate with RBAC

Production-ready Gin boilerplate with RBAC, JWT, GORM (MySQL), and scalable APIs.

## Setup
1. Start MySQL: `docker run -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:8`
2. Copy `.env.example` to `.env` and fill values.
3. Install deps: `go mod tidy`.
4. Generate TLS certs: `openssl req -new -newkey rsa:4096 -x509 -sha256 -days 365 -nodes -out cert.pem -keyout key.pem`.
5. Run: `go run ./cmd/server/main.go` (dev) or `docker build -t gin-app . && docker run -p 8080:8080 gin-app` (prod).

## Structure
- cmd/server: Entry point.
- internal/config: Env/config.
- internal/database: GORM/MySQL setup and migrations.
- internal/errs: Custom errors.
- internal/handler: Auth and user APIs.
- internal/lib: JWT and token store.
- internal/logger: Structured logging.
- internal/middleware: Auth, logging, timeout.
- internal/model: User and role models.
- internal/router: Route definitions.
- internal/service: Auth and user logic.
- internal/validation: Custom validators.
- static/: Static files.

## APIs
- POST /register: Register user.
- POST /login: Login, returns access/refresh tokens.
- POST /refresh: Refresh access token.
- POST /logout: Blacklist refresh token.
- GET /api/profile: Get user profile (JWT).
- PUT /api/profile: Update profile (JWT).
- DELETE /api/profile: Delete profile (JWT).
- GET /api/admin/users: List users (admin).
- POST /api/admin/users: Create user (admin).
- PUT /api/admin/users/:id: Update user (admin).
- DELETE /api/admin/users/:id: Delete user (admin).
- GET /health: Health check.
- GET /static/*: Static files.

## Best Practices
- HTTPS, secure headers (CSP, X-Frame-Options).
- JWT with RBAC (user/admin roles).
- Rate limiting (10 req/s), CORS, timeouts (5s).
- Password hashing (bcrypt).
- Structured logging (logrus).
- GORM with MySQL connection pooling.
- Token blacklisting (in-memory, Redis-ready).
- Graceful shutdown.
- Dockerized deployment.

## Next Steps
- Enable Redis for token blacklisting.
- Add tests (unit, integration).
- Setup CI/CD (GitHub Actions).
- Integrate Prometheus/Grafana.
- Scan with Trivy for security.