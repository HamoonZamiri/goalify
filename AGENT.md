# Goalify

Goalify is a full-stack productivity web application that tracks goals and goal categories, giving users XP when they complete tasks. TypeScript Vue 3 Composition API frontend with Golang standard library backend.

## Project Structure
- `frontend/` - Vue 3 + TypeScript + Tailwind CSS frontend
- `backend/` - Golang HTTP server with PostgreSQL
- Frontend dev server: http://localhost:5173  
- Backend dev server: http://localhost:8080
- Database: PostgreSQL on port 5432 running in a Docker container

## Build & Commands

### Frontend (run from `frontend/` directory)
- Fix linting/formatting: `pnpm format:fix`
- Run tests: `pnpm test`
- Run single test: `pnpm test src/file.test.ts`
- Start development server: `pnpm dev`
- Build for production: `pnpm build`
- Preview production build: `pnpm preview`

### Backend (run from `backend/` directory)
- Build: `make build`
- Development with hot reload: `make dev`
- Development with pretty JSON logs: `make jqdev`
- Run tests: `make test`
- Unit tests only: `make unit`
- Integration tests only: `make inte`
- Start with Docker services: `make start`
- Docker services up: `make up`
- Docker services down: `make down`

## Code Style
- TypeScript: Strict mode with Composition API
- Tabs for indentation (2 spaces for YAML/JSON/MD)
- Use JSDoc docstrings for TypeScript definitions, not `//` comments
- 80 character line limit
- Imports: Use consistent-type-imports
- Use descriptive variable/function names
- Golang: Follow standard Go conventions, use `gofmt`

## Testing
- Frontend: Vitest for unit testing, @vue/test-utils for component rendering
- Backend: Go standard testing package
- Use `expect(VALUE).toXyz(...)` instead of storing in variables
- Omit "should" from test names (e.g., `it("validates input")` not `it("should validate input")`)
- Test files: `*.test.ts` or `*.spec.ts` (frontend), `*_test.go` (backend)

## Architecture
- Frontend: Vue 3 Composition API with TypeScript
- Backend: Golang standard library HTTP server
- Database: PostgreSQL with Goose migrations  
- Styling: Tailwind CSS with custom component primitives (Box, Text, InputField)
- Build: Vite (frontend), Go toolchain (backend)
- Package manager: pnpm (frontend)
- Containerization: Docker Compose for development services

## Security
- Use appropriate data types that limit exposure of sensitive information
- Never commit secrets or API keys to repository
- Use environment variables for sensitive data (.env.example provided)
- Validate all user inputs on both client and server
- Use HTTPS in production
- Regular dependency updates
- Follow principle of least privilege

## Git Workflow
- ALWAYS run `pnpm format:fix` before committing frontend changes
- Run `pnpm build` to verify typecheck passes
- Run `make test` for backend changes
- NEVER use `git push --force` on the main branch
- Use `git push --force-with-lease` for feature branches if needed
- Always verify current branch before force operations

## Configuration
When adding new configuration options:
1. Environment variables in `.env.example`
2. Update both frontend and backend configurations as needed
3. Document in this AGENT.md file
4. Backend: Use consistent environment variable naming
5. Frontend: Use Vite's env variable conventions

All configuration keys use consistent naming and MUST be documented.
