# FlagForge Dashboard

This Next.js application provides a minimal dashboard for managing feature flags. The scaffold renders a placeholder list of flags and communicates with the FlagForge API when available.

## Scripts

- `npm run dev` — start the development server.
- `npm run build` — create a production build.
- `npm run start` — serve the production build.
- `npm run lint` — run ESLint.
- `npm run typecheck` — run TypeScript checks.

## Environment

Set `NEXT_PUBLIC_API_BASE_URL` to point to the API service. Docker Compose wires this automatically via the shared `.env` file.
