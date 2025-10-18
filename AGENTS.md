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

### Backend
- Golang standard library HTTP server
- PostgreSQL with Goose migrations
- Standard Go package structure

### Frontend
- Vue 3 Composition API with TypeScript
- TanStack Query for data fetching & caching
- Zod for schema validation
- Styling: Tailwind CSS with custom UI primitives
- Build: Vite
- Package manager: pnpm

**Frontend Folder Structure (Feature-Based):**
```
frontend/src/
├── features/          # Domain-driven feature modules
│   └── goals/
│       ├── queries/          # TanStack Query hooks with inline queryDataFns
│       ├── schemas/          # Zod schemas (entities & forms)
│       ├── components/       # Feature-specific components
│       ├── forms/           # Form components
│       └── index.ts         # Barrel exports
├── shared/            # Shared utilities (monorepo-ready)
│   ├── components/ui/       # Primitives (Box, Button, Text, InputField)
│   ├── api/                 # zodFetch, query client
│   └── schemas/             # Common schemas
├── pages/             # Route components ONLY
├── types/             # Global TypeScript types
└── router/
```

**Data Flow:**
- Component → TanStack Query hook → inline queryDataFn → zodFetch → Backend
- TanStack Query handles caching, loading states, errors, and refetching
- Zod schemas validate API responses at runtime
- Each query/mutation hook includes its own queryDataFn in the same file

**Import Patterns:**
```typescript
// Feature imports (use in pages/cross-feature)
import { useGoalCategories, GoalCard } from "@/features/goals";

// Shared imports
import { Box, Button } from "@/shared/components/ui";
import { zodFetch } from "@/shared/api";

// Within same feature (use relative imports)
import { useGoalCategories } from "../queries";
import type { Goal } from "../schemas";
```

**Key Files:**
- `shared/api/client.ts` - zodFetch with automatic auth & token refresh
- `shared/api/query-client.ts` - TanStack Query configuration
- `features/goals/queries/queryKeys.ts` - Hierarchical query key factory
- `features/goals/schemas/goal.schema.ts` - Zod entity schemas
- `features/goals/schemas/goal-form.schema.ts` - Form validation schemas
- All folders have `index.ts` barrel exports for clean imports

**zodFetch Automatic Authentication:**
- `zodFetch` automatically adds `Authorization: Bearer {token}` to all requests
- Gets token from `useAuth().getUser()?.access_token` internally
- No manual token passing needed in query hooks
- If user is not logged in (token undefined), header is omitted
- Server middleware ignores missing auth for `/login` and `/signup` endpoints
- On 401 response, automatically attempts token refresh and retries request

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

## Frontend Patterns & Best Practices

### TanStack Query Usage (One Hook Per File)
```typescript
// Query hook pattern with inline queryDataFn
// features/goals/queries/categories/useGoalCategories.ts

async function goalCategoriesQueryDataFn(): Promise<GoalCategory[]> {
  const result = await zodFetch(
    `${API_BASE}/goals/categories`,
    GoalCategoryResponseArraySchema,
  );

  if (isErrorResponse(result)) {
    throw new Error(result.message);
  }

  return result.data;
}

export function useGoalCategories() {
  return useQuery({
    queryKey: categoryKeys.lists(),
    queryFn: goalCategoriesQueryDataFn,
  });
}

// Mutation hook pattern with inline queryDataFn
// features/goals/queries/goals/useCreateGoal.ts
async function createGoalQueryDataFn(
  data: CreateGoalFormData
): Promise<Goal> {
  const result = await zodFetch(`${API_BASE}/goals`, GoalSchema, {
    method: http.MethodPost,
    body: JSON.stringify(data),
  });

  if (isErrorResponse(result)) {
    throw new Error(result.message);
  }

  return result;
}

export function useCreateGoal() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createGoalQueryDataFn,
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: categoryKeys.all });
      toast.success(`Goal created: ${data.title}`);
    },
    onError: (error: Error) => {
      toast.error(`Failed to create goal: ${error.message}`);
    },
  });
}

// Component usage
const { data, isLoading, error } = useGoalCategories();
const { mutate: createGoal, isPending } = useCreateGoal();
```

**Naming Convention:**
- Query data functions follow the pattern: `{hookName without "use"}QueryDataFn()`
- Example: `useGoalCategories()` hook → `goalCategoriesQueryDataFn()`
- Example: `useCreateGoal()` hook → `createGoalQueryDataFn()`
- Data functions are NOT exported - they're internal to the hook file
- Each hook lives in its own file for better organization and Git diffs

**Authentication & Error Handling:**
- `zodFetch` automatically adds Authorization headers using `useAuth().getUser()?.access_token`
- No need to pass tokens manually - zodFetch handles it internally
- If token is undefined (user not logged in), it's simply omitted from headers
- Server returns 401 for auth failures, triggering automatic token refresh flow
- queryDataFns handle error checking internally and throw errors
- Return type is clean data (`Promise<Goal>`) not union (`Promise<Goal | ErrorResponse>`)
- Hooks simply call the queryDataFn and let errors propagate to TanStack Query
- onError callbacks in mutations handle toast notifications

### Query Keys (Hierarchical Pattern)
```typescript
// features/goals/queries/queryKeys.ts
export const categoryKeys = {
  all: ['categories'] as const,
  lists: () => [...categoryKeys.all, 'list'] as const,
  list: (filters: any) => [...categoryKeys.lists(), filters] as const,
  details: () => [...categoryKeys.all, 'detail'] as const,
  detail: (id: string) => [...categoryKeys.details(), id] as const,
};

// Invalidation
queryClient.invalidateQueries({ queryKey: categoryKeys.all }); // All categories
queryClient.invalidateQueries({ queryKey: categoryKeys.lists() }); // Just lists
```

### Query File Structure (One Hook Per File)
Each query/mutation hook lives in its own file organized by subdomain:

```
features/goals/queries/
├── queryKeys.ts              # Hierarchical query key factory
├── categories/               # Category-related queries
│   ├── useGoalCategories.ts
│   ├── useCreateGoalCategory.ts
│   ├── useUpdateGoalCategory.ts
│   ├── useDeleteGoalCategory.ts
│   └── index.ts             # Barrel export for category hooks
├── goals/                   # Goal-related queries
│   ├── useCreateGoal.ts
│   ├── useUpdateGoal.ts
│   ├── useDeleteGoal.ts
│   └── index.ts             # Barrel export for goal hooks
└── index.ts                 # Barrel export for all hooks
```

**Benefits:**
- **NO separate `api/` directory** - queryDataFns are colocated with hooks
- One hook per file (~30-40 lines each)
- Easy to find and navigate specific hooks
- Cleaner Git diffs - changes to one hook don't dirty large files
- Natural organization by subdomain (categories vs goals)
- queryDataFns are private to the file (not exported)
- Only hooks are exported

### TanStack Form (Create & Edit)
```typescript
// forms/EditGoalForm.vue
const form = useForm({
  defaultValues: { title: props.goal.title, ... },
  validators: { onChange: editGoalFormSchema },
});

// Subscribe to reactive form state for auto-save
const formValues = form.useStore((state) => state.values);
const isDirty = form.useStore((state) => state.isDirty);
const isValid = form.useStore((state) => state.isValid);

// Hybrid auto-save: debounced + on close
watchDebounced(formValues, async (values) => {
  if (!isDirty.value || !isValid.value) return;
  await updateGoal({ goalId, data: values });
}, { debounce: 500, deep: true });
```

**Key Patterns:**
- Use `form.useStore()` to get reactive refs for `watchDebounced`
- Only save when `isDirty && isValid`
- Hybrid auto-save: debounced (500ms) + explicit save on close
- Silent validation (no error display for inline editing)
- Separate schemas: `createSchema` (all required), `editSchema` (all required), `updateSchema` (all optional for PATCH)

### Zod Schemas (Single Source of Truth)
```typescript
// Define schema once
const GoalSchema = z.object({
  id: z.string().uuid(),
  title: z.string(),
  // ...
});

// Runtime validation
const goal = GoalSchema.parse(json);

// TypeScript types
export type Goal = z.infer<typeof GoalSchema>;
```

### Adding New Features
When adding a new feature (e.g., `features/auth/`):
1. Create folder structure: `queries/`, `schemas/`, `components/`, `forms/`
2. Define Zod schemas in `schemas/`
3. Create query keys in `queries/queryKeys.ts`
4. Create subdomain folders within `queries/` (e.g., `queries/sessions/`, `queries/users/`)
5. Create one file per hook with inline queryDataFn using `{hookNameWithoutUse}QueryDataFn()` naming
6. Add barrel exports:
   - Each subdomain folder gets an `index.ts` exporting all its hooks
   - `queries/index.ts` exports all subdomain folders + queryKeys
7. Build components/forms with their own barrel exports
8. Export from feature root `index.ts`

**Example:**
```
features/auth/
├── queries/
│   ├── queryKeys.ts
│   ├── sessions/
│   │   ├── useLogin.ts
│   │   ├── useLogout.ts
│   │   └── index.ts
│   └── index.ts
├── schemas/
├── components/
├── forms/
└── index.ts
```

### Migration Status
- ✅ Goals feature fully migrated to new structure
- ⏳ Auth feature (uses old hooks/api/useApi.ts)
- ⏳ Levels feature
- Old code in `components/goals/`, `hooks/goals/`, `components/primitives/` can be deleted after validation
