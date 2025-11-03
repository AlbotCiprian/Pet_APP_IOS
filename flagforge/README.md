# FlagForge

FlagForge is a production-ready scaffold for a feature flag and remote configuration SaaS platform. It includes a Go API, background worker, Next.js dashboard, and SDK placeholders designed to give teams a head start on building a full-featured experimentation platform.

## Quickstart

1. Copy `.env.example` to `.env` and adjust values as needed.
2. Run the stack:

```bash
docker-compose up --build
```

3. Access the services:
   - API health check: [http://localhost:8080/healthz](http://localhost:8080/healthz)
   - Worker health check: [http://localhost:8090/healthz](http://localhost:8090/healthz)
   - Dashboard: [http://localhost:3000](http://localhost:3000)

## Services

```
flagforge/
├─ api/          # Go HTTP API service with clean architecture
├─ worker/       # Go worker handling background jobs and cache invalidation
├─ dashboard/    # Next.js dashboard for managing feature flags
├─ sdk/ios/      # Swift SDK placeholder
└─ deploy/       # Database migrations
```

## Environments & API Keys

FlagForge ships with support for multiple environments (`dev`, `stage`, `prod`) and a concept of API keys. Server-side keys authenticate management traffic, while client-side keys access evaluated flag payloads. Keys are scoped to environments so you can independently roll out features across environments.

## License

FlagForge is released under the [MIT License](LICENSE).
