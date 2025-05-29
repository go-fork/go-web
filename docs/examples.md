# Example Packages for Go Fork Registry

Đây là một số package mẫu bạn có thể thêm vào registry:

## Core Packages

### Authentication Package
```bash
./scripts/add-package.sh \
    "auth" \
    "go.fork.vn/auth" \
    "github.com/go-fork/auth" \
    "JWT authentication and authorization package"
```

### Database Package
```bash
./scripts/add-package.sh \
    "db" \
    "go.fork.vn/db" \
    "github.com/go-fork/db" \
    "Database abstraction layer with PostgreSQL and MySQL support"
```

### Logging Package
```bash
./scripts/add-package.sh \
    "logger" \
    "go.fork.vn/logger" \
    "github.com/go-fork/logger" \
    "Structured logging with multiple output formats"
```

## Provider Packages

### Redis Provider
```bash
./scripts/add-package.sh \
    "providers/redis" \
    "go.fork.vn/providers/redis" \
    "github.com/go-fork/providers/redis" \
    "Redis client with connection pooling"
```

### AWS Provider
```bash
./scripts/add-package.sh \
    "providers/aws" \
    "go.fork.vn/providers/aws" \
    "github.com/go-fork/providers/aws" \
    "AWS SDK wrapper with common utilities"
```

## Utility Packages

### HTTP Client
```bash
./scripts/add-package.sh \
    "http" \
    "go.fork.vn/http" \
    "github.com/go-fork/http" \
    "Enhanced HTTP client with retry and middleware support"
```

### Validation Package
```bash
./scripts/add-package.sh \
    "validator" \
    "go.fork.vn/validator" \
    "github.com/go-fork/validator" \
    "Struct validation with custom rules"
```

## Testing with curl

### Get all packages
```bash
curl http://localhost:8080/api/packages
```

### Get specific package
```bash
curl http://localhost:8080/api/packages/fork
```

### Add new package via API
```bash
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{
         "name": "test",
         "path": "go.fork.vn/test",
         "repository": "github.com/go-fork/test",
         "description": "Test package"
     }' \
     http://localhost:8080/api/packages
```

### Update package
```bash
curl -X PUT \
     -H "Content-Type: application/json" \
     -d '{
         "name": "test",
         "path": "go.fork.vn/test",
         "repository": "github.com/go-fork/test",
         "description": "Updated test package"
     }' \
     http://localhost:8080/api/packages/test
```

### Delete package
```bash
curl -X DELETE http://localhost:8080/api/packages/test
```
