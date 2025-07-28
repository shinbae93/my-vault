# vaultbox

A lightweight, self-hosted vault and secret manager for personal use. Securely store and manage passwords, tokens, and private data with AES-256-GCM encryption.

## Features

- 🔐 **Secure Encryption**: AES-256-GCM encryption with Argon2id key derivation
- 🏠 **Self-Hosted**: Complete control over your data
- 🔒 **Auto-Lock**: Automatic vault locking after 15 minutes of inactivity
- 📱 **RESTful API**: Clean API for integration
- 🐳 **Docker Ready**: Easy deployment with Docker Compose
- 🗄️ **PostgreSQL**: Reliable data storage

## Project Structure

```
my-vault/
├── backend/                 # Go API (Clean Architecture)
│   ├── cmd/                # Application entry point
│   ├── internal/
│   │   ├── handlers/       # HTTP handlers
│   │   ├── services/       # Business logic
│   │   ├── repository/     # Database access
│   │   ├── models/         # Data models
│   │   └── utils/          # Utilities (crypto, etc.)
│   ├── Dockerfile          # Backend container
│   ├── Makefile           # Build and development tasks
│   └── env.example        # Environment variables template
├── frontend/               # React 19 + Vite + shadcn/ui
├── docker-compose.yml      # Local development orchestration
└── README.md
```

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.22+ (for local development)
- Node.js 18+ (for frontend development)

### 1. Clone and Setup

```bash
git clone <your-repo-url>
cd my-vault
```

### 2. Backend Setup

```bash
cd backend

# Copy environment file
cp env.example .env

# Edit .env with your configuration
# (Optional: change master password, database settings)

# Install dependencies
make deps

# Run locally (requires PostgreSQL)
make run

# Or use Docker
make docker-build
make docker-run
```

### 3. Database Setup

The backend will automatically create the database schema on first run.

### 4. Frontend Setup (Optional)

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

### 5. Using Docker Compose (Recommended)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f backend

# Stop services
docker-compose down
```

## API Endpoints

### Vault Management

- `POST /api/unlock` - Unlock vault with master password
- `POST /api/lock` - Lock vault
- `GET /api/status` - Get vault status

### Secret Management (requires unlocked vault)

- `GET /api/secrets` - List all secrets
- `POST /api/secrets` - Create new secret
- `GET /api/secrets/:id` - Get specific secret
- `PUT /api/secrets/:id` - Update secret
- `DELETE /api/secrets/:id` - Delete secret

### Example Usage

```bash
# Unlock vault
curl -X POST http://localhost:3000/api/unlock \
  -H "Content-Type: application/json" \
  -d '{"master_password": "your-master-password"}'

# Create a secret
curl -X POST http://localhost:3000/api/secrets \
  -H "Content-Type: application/json" \
  -d '{
    "title": "GitHub Token",
    "type": "api_token",
    "value": "ghp_xxxxxxxxxxxxxxxxxxxx"
  }'

# List secrets
curl http://localhost:3000/api/secrets

# Lock vault
curl -X POST http://localhost:3000/api/lock
```

## Development

### Backend Development

```bash
cd backend

# Install dependencies
make deps

# Run with hot reload (requires air)
make dev

# Run tests
make test

# Build for production
make build-prod

# Format code
make fmt

# Lint code
make lint
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build
```

## Security Features

- **Argon2id Key Derivation**: Secure password-based key derivation
- **AES-256-GCM Encryption**: Military-grade encryption for secrets
- **In-Memory Keys**: Encryption keys never stored on disk
- **Auto-Lock**: Automatic vault locking after inactivity
- **CORS Protection**: Configured for local development

## Configuration

### Environment Variables

| Variable            | Description                 | Default       |
| ------------------- | --------------------------- | ------------- |
| `PORT`              | Server port                 | `3000`        |
| `DB_HOST`           | Database host               | `localhost`   |
| `DB_PORT`           | Database port               | `5432`        |
| `DB_USER`           | Database user               | `vaultbox`    |
| `DB_PASSWORD`       | Database password           | `supersecret` |
| `DB_NAME`           | Database name               | `vaultbox`    |
| `MASTER_PASSWORD`   | Master password             | `changeme`    |
| `AUTO_LOCK_TIMEOUT` | Auto-lock timeout (minutes) | `15`          |

## Production Deployment

1. **Change Default Passwords**: Update `MASTER_PASSWORD` and database credentials
2. **Enable HTTPS**: Set up SSL certificates and enable HTTPS
3. **Database Security**: Use strong database passwords and consider external database
4. **Network Security**: Configure firewall rules appropriately
5. **Backup Strategy**: Implement regular database backups

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues and questions, please open an issue on GitHub.
