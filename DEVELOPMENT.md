# BondiHub Development Guide

This guide covers setting up the development environment and contributing to the BondiHub project.

## 🛠️ Development Setup

### Prerequisites
- Go 1.21+ ([Download](https://golang.org/dl/))
- Node.js 18+ ([Download](https://nodejs.org/))
- PostgreSQL 14+ ([Download](https://www.postgresql.org/download/))
- Git ([Download](https://git-scm.com/downloads))
- Docker (optional) ([Download](https://www.docker.com/get-started))

### Backend Setup (Go)

#### 1. Clone Repository
```bash
git clone <repository-url>
cd BondiHub/backend
```

#### 2. Install Dependencies
```bash
go mod download
```

#### 3. Environment Configuration
```bash
# Copy environment file
cp env.example .env

# Edit configuration
nano .env
```

#### 4. Database Setup
```bash
# Start PostgreSQL
sudo systemctl start postgresql

# Create database
sudo -u postgres psql
CREATE DATABASE bondihub;
CREATE USER bondihub_dev WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE bondihub TO bondihub_dev;
\q
```

#### 5. Run Application
```bash
# Development mode
go run main.go

# Or with air for hot reload
go install github.com/cosmtrek/air@latest
air
```

### Frontend Setup (Angular)

#### 1. Navigate to Frontend
```bash
cd ../frontend
```

#### 2. Install Dependencies
```bash
npm install
```

#### 3. Environment Configuration
```bash
# Copy environment file
cp src/environments/environment.ts src/environments/environment.dev.ts

# Edit configuration
nano src/environments/environment.dev.ts
```

#### 4. Run Development Server
```bash
# Start development server
ng serve

# Or with specific port
ng serve --port 4200
```

### Database Migrations

#### Using GORM AutoMigrate
```bash
# Run migrations
go run main.go
# Migrations run automatically on startup
```

#### Manual SQL Migrations
```bash
# Create migration
psql -h localhost -U bondihub_dev -d bondihub -f migrations/001_initial_schema.sql
```

## 🏗️ Project Structure

```
BondiHub/
├── backend/                 # Go backend API
│   ├── config/             # Configuration files
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # Middleware functions
│   ├── models/             # Data models
│   ├── routes/             # Route definitions
│   ├── services/           # Business logic
│   ├── utils/              # Utility functions
│   ├── migrations/         # Database migrations
│   ├── main.go            # Application entry point
│   ├── go.mod             # Go dependencies
│   └── Dockerfile         # Docker configuration
├── frontend/               # Angular frontend
│   ├── src/
│   │   ├── app/
│   │   │   ├── core/       # Core services and models
│   │   │   ├── features/   # Feature modules
│   │   │   ├── shared/     # Shared components
│   │   │   └── app.component.ts
│   │   ├── assets/         # Static assets
│   │   ├── environments/   # Environment configs
│   │   └── styles.scss     # Global styles
│   ├── angular.json        # Angular configuration
│   ├── package.json        # Node dependencies
│   └── tailwind.config.js  # Tailwind configuration
├── nginx/                  # Nginx configuration
├── docker-compose.yml      # Docker services
├── README.md              # Project overview
├── API.md                 # API documentation
├── DEPLOYMENT.md          # Deployment guide
└── DEVELOPMENT.md         # This file
```

## 🔧 Development Tools

### Backend Tools

#### Code Formatting
```bash
# Format Go code
go fmt ./...

# Or use goimports
go install golang.org/x/tools/cmd/goimports@latest
goimports -w .
```

#### Linting
```bash
# Install golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2

# Run linter
golangci-lint run
```

#### Testing
```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./handlers -v
```

#### Database Tools
```bash
# Install sqlc for type-safe SQL
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate code from SQL
sqlc generate
```

### Frontend Tools

#### Code Formatting
```bash
# Format TypeScript/HTML
npx prettier --write "src/**/*.{ts,html,scss}"

# Or use Angular CLI
ng format
```

#### Linting
```bash
# Run ESLint
ng lint

# Fix linting issues
ng lint --fix
```

#### Testing
```bash
# Run unit tests
ng test

# Run e2e tests
ng e2e

# Run tests with coverage
ng test --code-coverage
```

#### Build
```bash
# Development build
ng build

# Production build
ng build --prod

# Build with specific environment
ng build --configuration=production
```

## 🧪 Testing

### Backend Testing

#### Unit Tests
```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./handlers -v

# Run tests with race detection
go test -race ./...
```

#### Integration Tests
```bash
# Run integration tests
go test -tags=integration ./...

# Run with test database
TEST_DB_URL="postgres://user:pass@localhost/bondihub_test" go test ./...
```

#### API Testing
```bash
# Install hey for load testing
go install github.com/rakyll/hey@latest

# Test API endpoints
hey -n 1000 -c 10 http://localhost:8080/api/v1/health
```

### Frontend Testing

#### Unit Tests
```bash
# Run tests in watch mode
ng test

# Run tests once
ng test --watch=false

# Run specific test
ng test --include="**/auth.service.spec.ts"
```

#### E2E Tests
```bash
# Run e2e tests
ng e2e

# Run with specific browser
ng e2e --browsers=ChromeHeadless
```

## 🔍 Debugging

### Backend Debugging

#### Using Delve
```bash
# Install delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug application
dlv debug main.go

# Debug with arguments
dlv debug main.go -- --config=config.yaml
```

#### Logging
```bash
# Enable debug logging
GIN_MODE=debug go run main.go

# Structured logging
LOG_LEVEL=debug go run main.go
```

### Frontend Debugging

#### Angular DevTools
```bash
# Install Angular DevTools
npm install -g @angular/devtools

# Enable in browser
# Add to main.ts:
import { enableProdMode } from '@angular/core';
import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';
import { AppModule } from './app/app.module';

if (environment.production) {
  enableProdMode();
}

platformBrowserDynamic().bootstrapModule(AppModule)
  .catch(err => console.error(err));
```

#### Browser DevTools
- Use Chrome DevTools for debugging
- Enable source maps in development
- Use Angular Augury extension

## 📝 Code Style Guidelines

### Go Code Style

#### Naming Conventions
- Use camelCase for variables and functions
- Use PascalCase for exported types and functions
- Use snake_case for database fields
- Use descriptive names

#### Error Handling
```go
// Good
if err != nil {
    return fmt.Errorf("failed to process user: %w", err)
}

// Bad
if err != nil {
    return err
}
```

#### Comments
```go
// Package handlers provides HTTP handlers for the API
package handlers

// UserHandler handles user-related HTTP requests
type UserHandler struct {
    // userService handles user business logic
    userService *services.UserService
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
    // Implementation here
}
```

### TypeScript/Angular Code Style

#### Naming Conventions
- Use camelCase for variables and functions
- Use PascalCase for classes and interfaces
- Use kebab-case for file names
- Use descriptive names

#### Component Structure
```typescript
@Component({
  selector: 'app-user-list',
  standalone: true,
  imports: [CommonModule],
  template: `...`,
  styles: [`...`]
})
export class UserListComponent implements OnInit, OnDestroy {
  // Properties
  users: User[] = [];
  
  // Constructor
  constructor(private userService: UserService) {}
  
  // Lifecycle hooks
  ngOnInit(): void {}
  ngOnDestroy(): void {}
  
  // Public methods
  loadUsers(): void {}
  
  // Private methods
  private handleError(error: any): void {}
}
```

## 🚀 Performance Optimization

### Backend Optimization

#### Database Optimization
```sql
-- Create indexes for frequently queried fields
CREATE INDEX idx_houses_status ON houses(status);
CREATE INDEX idx_houses_landlord_id ON houses(landlord_id);
CREATE INDEX idx_payments_agreement_id ON payments(agreement_id);

-- Use prepared statements
-- Use connection pooling
-- Optimize queries
```

#### Memory Optimization
```go
// Use sync.Pool for frequently allocated objects
var userPool = sync.Pool{
    New: func() interface{} {
        return &User{}
    },
}

// Reuse objects
user := userPool.Get().(*User)
defer userPool.Put(user)
```

### Frontend Optimization

#### Bundle Optimization
```bash
# Analyze bundle size
ng build --stats-json
npx webpack-bundle-analyzer dist/bondihub-frontend/stats.json

# Enable tree shaking
ng build --prod --aot --build-optimizer
```

#### Lazy Loading
```typescript
// Lazy load feature modules
const routes: Routes = [
  {
    path: 'admin',
    loadChildren: () => import('./admin/admin.module').then(m => m.AdminModule)
  }
];
```

## 🔒 Security Best Practices

### Backend Security

#### Input Validation
```go
// Validate input
if err := c.ShouldBindJSON(&req); err != nil {
    utils.ValidationErrorResponse(c, "Invalid input", err)
    return
}
```

#### SQL Injection Prevention
```go
// Use parameterized queries
query := "SELECT * FROM users WHERE email = $1"
rows, err := db.Query(query, email)
```

#### CORS Configuration
```go
// Configure CORS
config := cors.Config{
    AllowOrigins:     []string{"http://localhost:4200"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}
```

### Frontend Security

#### XSS Prevention
```typescript
// Sanitize user input
import { DomSanitizer } from '@angular/platform-browser';

constructor(private sanitizer: DomSanitizer) {}

sanitizeHtml(html: string) {
  return this.sanitizer.sanitize(SecurityContext.HTML, html);
}
```

#### CSRF Protection
```typescript
// Include CSRF token in requests
import { HttpClientXsrfModule } from '@angular/common/http';

@NgModule({
  imports: [
    HttpClientXsrfModule.withOptions({
      cookieName: 'XSRF-TOKEN',
      headerName: 'X-XSRF-TOKEN'
    })
  ]
})
```

## 📊 Monitoring and Logging

### Backend Monitoring

#### Structured Logging
```go
import "github.com/sirupsen/logrus"

log := logrus.WithFields(logrus.Fields{
    "user_id": userID,
    "action": "create_house",
    "ip": c.ClientIP(),
})
log.Info("House created successfully")
```

#### Metrics Collection
```go
import "github.com/prometheus/client_golang/prometheus"

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
)
```

### Frontend Monitoring

#### Error Tracking
```typescript
// Track errors
import { ErrorHandler } from '@angular/core';

@Injectable()
export class GlobalErrorHandler implements ErrorHandler {
  handleError(error: any): void {
    // Send error to monitoring service
    console.error('Global error:', error);
  }
}
```

#### Performance Monitoring
```typescript
// Track performance metrics
import { NgZone } from '@angular/core';

constructor(private ngZone: NgZone) {}

trackPageLoad() {
  this.ngZone.runOutsideAngular(() => {
    const loadTime = performance.now();
    // Send to analytics
  });
}
```

## 🤝 Contributing

### Git Workflow

#### Branch Naming
- `feature/user-authentication` - New features
- `bugfix/payment-validation` - Bug fixes
- `hotfix/security-patch` - Critical fixes
- `refactor/database-models` - Code refactoring

#### Commit Messages
```
feat: add user authentication
fix: resolve payment validation issue
docs: update API documentation
style: format code with prettier
refactor: restructure database models
test: add unit tests for auth service
chore: update dependencies
```

#### Pull Request Process
1. Create feature branch
2. Make changes
3. Write tests
4. Update documentation
5. Create pull request
6. Code review
7. Merge to main

### Code Review Guidelines

#### What to Look For
- Code quality and readability
- Test coverage
- Security vulnerabilities
- Performance implications
- Documentation updates
- Breaking changes

#### Review Checklist
- [ ] Code follows style guidelines
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] No security issues
- [ ] Performance is acceptable
- [ ] No breaking changes

## 📚 Additional Resources

### Documentation
- [Go Documentation](https://golang.org/doc/)
- [Angular Documentation](https://angular.io/docs)
- [Gin Framework](https://gin-gonic.com/docs/)
- [PrimeNG Components](https://primeng.org/)
- [Tailwind CSS](https://tailwindcss.com/docs)

### Tools
- [Postman](https://www.postman.com/) - API testing
- [Insomnia](https://insomnia.rest/) - API testing
- [DBeaver](https://dbeaver.io/) - Database management
- [VS Code](https://code.visualstudio.com/) - Code editor
- [Chrome DevTools](https://developers.google.com/web/tools/chrome-devtools)

### Learning Resources
- [Go by Example](https://gobyexample.com/)
- [Angular University](https://angular-university.io/)
- [Database Design](https://www.lucidchart.com/pages/database-diagram/database-design)
- [REST API Design](https://restfulapi.net/)
