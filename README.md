# sw-services-go

Go (Gin/GORM) rewrite of DIGIT-OSS's Java Spring Boot `sw-services` (Sewerage
Connection Service). See `../../Resource/Conversion Guide sw-services.pdf` for
the mandated folder structure and `../../Resource/sw-services issues.pdf` for
the parity gaps this rewrite is built to close.

## Status: Phase 1 (foundation)

Implemented:
- Full folder-structure compliance (`cmd/`, `internal/{config,domain,transport,service,repository,validator,workflow,util}`, `migrations/`, `configs/`, `deploy/`, `docs/`).
- Normalized Postgres schema across 6 tables (`eg_sw_connection`, `eg_sw_service`, `eg_sw_plumberinfo`, `eg_sw_applicationdocument`, `eg_sw_connectionholder`, `eg_sw_roadcuttinginfo`) instead of collapsing child entities into JSONB — fixes the data-loss finding in the issues report. Also fixes two real Java schema bugs: the erroneous `UNIQUE(connectionid)` on connection holders, and `holdershippercentage` stored as numeric instead of varchar.
- `POST /sw-services/swc/_create`, `_update` (fetch-before-merge sparse update — a partial payload can't wipe existing columns), `_search`, `_plainsearch`, `_encryptOldData` (mirrors Java's own disabled endpoint), `GET /sw-services/actuator/health`.
- Config externalized via `configs/application.yaml` + env override (no hardcoded credentials).
- PII masking/unmasking convention (`internal/util/masking.go`) implemented for real; field-level encryption at rest is stubbed behind an `Encryptor` interface (`internal/util/encryptor.go`) rather than depending on the proprietary ABAC encryption microservice Java uses — see plan for rationale.

Not yet implemented (see roadmap in the plan file used to build this):
IDGen/Workflow/MDMS/Property client integrations, Kafka producers/consumers, User service integration, sw-calculator integration, PDF/Filestore/URL-shortener clients, and the richer workflow-driven validators (`ActionValidator`, `SewerageFieldValidator`, `MDMSValidator`).

## Running locally

```bash
# 1. Apply the schema
DATABASE_URL=postgres://postgres:postgres@localhost:5432/rainmaker_new?sslmode=disable ./deploy/migrate.sh

# 2. Run the service (reads configs/application.yaml from the working directory)
go run ./cmd/sw-services
```

The service listens on `:8091` under context path `/sw-services`. Import
`docs/postman/sw-services-postman.json` for ready-made requests.

## Testing

```bash
go build ./...
go vet ./...
```
