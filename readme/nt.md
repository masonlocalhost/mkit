# ORM, database and migrations

## Database

Use Postgresql v14, code-first approach.

## ORM and Query customization

Use GORM (https://gorm.io/).

- Version v1.30.0.
- Use traditional API (Not generic API).
- Models defined in `nt-common` project.

For complex query, use Bob query builder (https://bob.stephenafamo.com/docs/query-builder/intro) and execute with Gorm Raw (with Find(), Scan(), Exec()...).

- Version v0.42.0.
- Postgresql dialect.

Tips:

- All db interaction code should be located in `internal/repository/`.
- For gorm has-one/ belong-to relation, use Joins() instead of Preload() for eager load without secondary query.
- For custom query, use gorm Raw() with Find() or First(), and the query is preloadable but cannot use join preload.

## Database migrations

### Using Atlas Migration Tool

Use [atlas v1.1.0](https://github.com/ariga/atlas/releases/tag/v1.1.0) to generate schema migration files from GORM models
using the Program Mode strategy (view [docs](https://atlasgo.io/guides/orms/gorm/program)).

To install atlas, run command below:

```bash
 go install ariga.io/atlas/cmd/atlas@v1.1.0
```

Project structure:

- `atlas.hcl` for configuration
- `db/loader` for loading GORM models
- `db/migrations` for storing SQL migration files

Atlas detects schema changes when:

- New models are added to the loader
- Existing models are modified

It will then generate a new migration file reflecting those changes. You can also add custom migration files,
but only DML statements are allowed (e.g., inserting seed data or migrating existing data).

### Steps to Add a New Migration

1. Configure a local database connection in `atlas.hcl` so Atlas can inspect the schema.
   ```hcl
   ...
   env "gorm" {
     src = data.external_schema.gorm.url
     dev = "postgres://postgres:password@localhost:5432/nt_db?sslmode=disable"
   ...
   ```
2. Add new models or update existing ones.  
   If adding a new model, register it in `db/loader/loader.go`.

   ```go
   func main() {
       stmts, err := gormschema.New("postgres").Load(
       // ínert models to generate DDL sql
       &domain.ntUserHistory{},
       &domain.BlacklistedEntity{},
       &domain.ntUserConfig{},
       &domain.ntUserInvokeToken{},

       // Add new models here!
       &domain.YourNewModel{},
      )
   }
   ```

3. Run Atlas diff to generate a migration file.  
   The file will be created in `db/migrations`.

   ```bash
   # Run diff, migration file will be created, for example: 20260325035225_add_new_model.up.sql
   atlas migrate diff add_new_model.up --env gorm
   ```

4. If needed, add custom DML statements to the generated file or create a new migration file.  
   After modifying or adding migrations, re-hash the migration directory.

   ```bash
   # Create custom migration
   atlas migrate new my_custom_migration.up --dir "file://db/migrations"
   ```

   After update custom migration, you must `re hash` by running:

   ```bash
   atlas migrate hash --dir "file://db/migrations"
   ```

5. Apply migrations manually or execute them when the service starts.

### Tips

- For the migration run when service start, migration file must be post-fixed with `.up`, eg: `add_abc.up.sql`
- Define all schema elements, including constraints, in GORM models. Do not add them manually in custom migration files.
- For partial indexes, use PostgreSQL’s normalized form rather than the common shorthand form.

  ```sql
  -- Common form
  CREATE INDEX record_partial_idx
  ON public.records USING btree (type)
  WHERE type IN ('type1', 'type2') AND completed = false;

  -- Normalized form
  CREATE INDEX record_partial_idx
  ON public.records USING btree (type)
  WHERE (((type)::text = ANY ((ARRAY['type1'::character varying, 'type2'::character varying])::text[]))
  AND (completed = false));
  ```

  ```go
   // usage in gorm model
   type Record struct {
       ID        uint   `gorm:"primaryKey" json:"id"`

       // Partial index in Gorm
       Type      string `gorm:"type:varchar(100);index:record_partial_idx,where:((type::text = ANY ((ARRAY['type1'::character varying\\, 'type2'::character varying])::text[])) AND (completed = false))" json:"type"`

       Completed bool   `gorm:"default:false" json:"completed"`
       Data      string `json:"data"`
   }
  ```

### Use gorm automigrate (DEPRECATED, DO NOT USE)

Use golang code for migration, library: github.com/go-gormigrate/gormigrate/v2 (version v2.1.5)

- Migration code can be found in `db/migration.go`, do not edit existed migration code.
- Use gorm auto migration for models, remember to call CustomMigrate if there are custom schema code (like custom unique index, attached to model). For example:
  ```golang
  func (ie *IssueEntity) CustomMigrate(db *gorm.DB) error {
      return db.Exec(fmt.Sprintf(`
          CREATE UNIQUE INDEX IF NOT EXISTS uk_issue_entity_entity_id_issue_type_issue_sub_type
          ON %s(%s, %s)
          WHERE %s = '%s' AND %s = '%s'
      `,
          IssueEntityTableName,
          IssueEntity_ENTITY_ID_COLUMN,
          IssueEntity_ISSUE_SUB_TYPE_COLUMN,
          IssueEntity_ISSUE_TYPE_COLUMN,
          IssueTypeExposure,
          IssueEntity_ISSUE_SUB_TYPE_COLUMN,
          IssueSubTypeLeakedCredentials,
      )).Error
  }
  // Implements CustomMigrator interface
  type CustomMigrator interface {
      CustomMigrate(db *gorm.DB) error
  }
  ```
- Warning: Code-first migrations (e.g., using GORM AutoMigrate) may not be fully repeatable or deterministic across
  environments, unlike `.sql`-based migrations. It is recommended to generate and review `.sql` migration scripts from
  GORM models using tool like [atlas](https://gorm.io/docs/migration.html#Atlas-Integration), especially for production use.

## Data cleanup queries

Some reserved queries.

```sql
delete from technologies where (
    exists(
        select 1 from technology_entities where technology_id = id
    )
)

delete from issues where (
    not exists(
        select 1 from issue_entities where issue_id = id
    )
)
```

# Error Handling

## Custom Error Type

This project uses custom error types defined in:
`nt-common/pkg/error/serviceerror`

### Purpose

Custom error types are used for:

- Explicit error handling across layers (`repo -> service -> handler`)
- Automatic HTTP response mapping
- Automatic logging via middleware (only logs internal server errors)
- Attaching structured metadata for responses (e.g., status code, message, custom payloads — TODO)
- Differentiating expected errors from unexpected ones

### Behavior

- Errors **not** using the custom error type are treated as **internal errors**
- Internal errors are:
  - Logged by middleware
  - Returned as generic server errors to clients

### Repository Errors

- Repository-specific errors (e.g., constraint violations, duplicate keys)
  should be:
  - Defined within the repository package
  - Mapped to appropriate custom error types before propagating upward

# RESTful & gRPC APIs

Use the reusable server in:
`nt-common/pkg/server`

### Features

- Supports **RESTful** and **gRPC** APIs
- Built-in **middlewares**
- Manages **server and dependency lifecycle**
- Provides **logging, metrics, and tracing** (OpenTelemetry)
- Centralized **configuration management**

---

## RESTful

- Framework: [go-gin](https://gin-gonic.com/)
- Controllers: `/internal/controller`
- Swagger UI (dev only): `/api-docs/`
- API spec: `/pkg/restful/docs/spec.yaml`

### Setup

1. **Define DTOs**
   - Objects: `/pkg/dto`
   - Requests/Responses: `/pkg/api/restful/rpc`
   - Use [validator v10](https://github.com/go-playground/validator) for validation  
     (automatically triggered during Gin request binding)

2. **Implement Controllers**
   - Add handlers in `/internal/controller`
   - Register routes to the server

3. **Update API Spec**
   - Add endpoints to `/pkg/restful/docs/spec.yaml`

---

## gRPC

- Library: [grpc v1.77.0](https://google.golang.org/grpc)
- Proto files: located in `nt-common`
- Code generation: via `buf CLI` (see `nt-common` README)

### Setup

1. **Define Protos**
   - Location: `/api/proto/{service}/...`
   - Define messages and RPCs
   - Optional: Add REST support using `google.api.http` annotations and  
     [vanguard-go](https://github.com/connectrpc/vanguard-go)

   ```protobuf
   import "google/api/annotations.proto";

   service ntCore {
     rpc Transform(ncrpc.TransformRequest) returns (ncrpc.TransformReply) {
       option (google.api.http) = {
         post: "/nt-core/transform"
         body: "*"
       };
     }
   }
   ```

   - This allows REST clients to call: POST /nt-core/transform

   ```go
    encoding.RegisterCodec(vanguardgrpc.NewCodec(&vanguard.JSONCodec{
    MarshalOptions:   protojson.MarshalOptions{EmitUnpopulated: true},
    UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
    }))

    transcoder, err := vanguardgrpc.NewTranscoder(s.Deps.GRPCServer)
    if err != nil {
    logger.Fatalf("cannot init gRPC transcoder: %v", err)
    }

    restHandler := h2c.NewHandler(transcoder, &http2.Server{})
   ```

2. **Generate Code**

   Generate protobuf models and gRPC stubs using `buf`

3. **Implement handlers**

   Location: /internal/server/{service}/handler.go

# Logging, Tracing, and Metrics

All logs, traces, and metrics are collected using **OpenTelemetry**.  
The server in `nt-common` supports exporting this data to an OpenTelemetry collector.

---

## Logging

- Library: [logrus](https://github.com/sirupsen/logrus)
- Integrated with middleware to send logs to the collector

### Guidelines

- Log only when necessary
  - Use appropriate log levels to avoid noise and excessive storage usage
- The server middleware:
  - Detects errors from `nt-common/pkg/error/serviceerror`
  - Automatically logs **internal errors**
  - For non-internal errors, log manually if needed
- Always include **context** in logs
  - Enables correlation and aggregation in the collector
- For non-API flows (e.g., background jobs, workers):
  - Logging must be handled manually for all levels
    Always propagate the returned ctx to downstream calls to maintain trace continuity

---

## Tracing & Metrics

- Most `nt-common` integrations already include tracing and metrics:
  - `gin`, `grpc`, `gorm`, `redis`, etc.
- For custom flows (e.g., cronjobs, background tasks), you must create traces manually
- Always propagate the returned ctx to downstream calls to maintain trace continuity

### Example

```go
// Start custom trace
tr := otel.Tracer("manage_collection_autostart")

ctx, span := tr.Start(
    s.ctx,
    "CronJob: Check and start a collection scanning regularly",
)

// End trace
defer span.End()
```
