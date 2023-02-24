# Go Testapp
Go testapp is a sample application which can be used as boilerplate for further Kargo projects in Golang. It features several key of application components that all Kargo golang application must implementes. In order to make a new Golang service, just copy this repository as a boilerplate to the new repository and just change the naming and entity.

## Features
- [x] Middleware
- [x] Config Initialization (yaml, json, env)
- [x] DB Migration (dbmate)
- [x] Customizeable Validator Package 
- [x] Gin (HTTP Handler)
- [ ] Custom Metrics
- [x] Clean Architecture
- [x] GORM
- [x] Error Structure
- [x] GQL Server
- [x] HTTP server
- [ ] Global Response Structure
- [x] Standardized Logging
- [ ] Unit Testing
- [ ] Swagger
- [x] CI/CD Workflow
- [x] Architecture Decision Records
- [x] Code Quality and Security Analysis
- [x] OpenTelemetry Tracing

## Setup

### Install Docker

We will use docker to install any dependencies that the service uses such as postgres, elasticsearch, redis, etc. To install docker, it is recommended that you follow the official [Get Docker](https://docs.docker.com/get-docker/) guideline.

### Install DB Migration Tools (+ Postgresql Driver)
dbmate
```
sudo curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64
sudo chmod +x /usr/local/bin/dbmate
```

### Install GQL Generator
For GQL in Golang, we use `gqlgen` library. 
Installation of `gqlgen` is on per project basis using `go mod`. Quick Start details are specified in this [link](https://github.com/99designs/gqlgen#quick-start) (for new project only) 

### Install Mockgen
For mock generation
```
go install github.com/golang/mock/mockgen@v1.6.0
```

### Install pre-commit
For setup github hooks
```
# Install pre-commit
pip install pre-commit

# Setup git pre-commit hooks
pre-commit install
```

### Install Horusec
For pre-commit security analysis
```
curl -fsSL https://raw.githubusercontent.com/ZupIT/horusec/master/deployments/scripts/install.sh | bash -s latest
```

## Getting Started with the Service Locally
### Making Database Migration
In general, there's two kind of migration in go-testapp, `default` and `postmigration`. Both of them have different migration directory and different migration tracker. (For more info on these two kind of migration, see `Database Migration Flow`):
1. `default`: this migration is migration that is automatically applied everytime there is application deployment.
2. `postmigration` this migration triggered manually independently of application deployment, it could be triggered via Github Action `db-operation`.

In order to do DDL operations on our Database, we use `dbmate` command which will run a database migration in sequential manner. The step 
to make a database migrations are as follows:
1. Create your DB migration folder as `db/<database-name>/dbmate/migrations`
2. Fill the `files/etc/env/env.sh` using `files/etc/env/env.sample` as starting point.
3. Run `source files/etc/env/env.sh` in terminal
4. Run `source ./files/etc/dbmate_config/dbmate.sh` in terminal (Alternatively, for `postdeployment` migration, `DBMATE_MIGRATION_TYPE=postdeployment source ./files/etc/dbmate_config/dbmate.sh`)
5. Create a DDL for `up` and `down` migration using CLI: `dbmate new <migration_file_name>`. This will give you the skeleton ONLY. You need to fill in the DDL query.
6. Run `dbmate up [N]` to run DDL command until the most recent one, if N is not filled. If N is filled, then it will apply `up` migration for the next 2 migrations sql from the current checkpoint.
7. To revert already applied migration, run `dbmate down [N]` to run DDL command until the beginning (if N is not filled). If N is filled, then it will apply `down` migration for the latest 2 migrations sql from the current checkpoint.

### Configuration
In this app, there are 2 different configurations:
- **Non-secret configuration**: stored in `files/etc/app_config/config.yaml`
- **Secret configuration**: stored in `files/etc/env/env.sh`

Secret is meant ONLY for configuration that people are not supposed to know. This includes, but not limited to:
- **Username**
- **Credentials**
- **Private Key**

Other config should go to `files/etc/app_config/config.yaml`. `files/etc/env/env.sh` can be constructed using skeleton from `files/etc/env/env.sample` in development. In staging and production, `env.sh` will be provided automatically via Kubernetes secret.

To run in local using the default postgres docker config, `files/etc/env/env.sh` should be filled with:
```
# Database Config
export DB_HOST=localhost
export DB_USERNAME=postgres
export DB_PASSWORD=postgres
export DB_PORT=5433
export DB_NAME=postgres
``` 

### Running the Service

Once you have all the required configuration run `cmd/start.sh` from go-testapp root directory to run the service.

## Project Structure

### Clean Architecture
![image](https://user-images.githubusercontent.com/102520846/172805794-7bc613ec-30d3-4898-8a5f-144ce3bb5b74.png)

We use clean architecture to separate between request receiver, business logic, and data layer logic. The main advantage of this separation of layer is to enable:
* Separation of concern
* Parallel works on each layer
* Adaptability to changing handler or data layer

As such, our core code will be in the form of:

```
|--- internal
|    |--- handler
|    |    |--- gql
|    |    |    |--- module A
|    |    |    |--- module B
|    |    |--- http
|    |    |    |--- module A
|    |    |    |--- module B
|    |
|    |--- usecase
|    |    |--- module A
|    |    |--- module B
|    |
|    |--- repo
|    |    |--- module A
|    |    |--- module B
|    |
|    |--- entity
|    |    |--- entity_a.go
|    |    |--- entity_b.go
|    |
```

Here are the details on what each component does:
* **Handler**: Presenter layer which does conversion of data structure from entities from and to a well-defined format of choosing like gRPC, HTTP, or GQL
* **Usecase**: Where business logic lives, this layer orchestrated entity and repository to achieve application specific needs
* **Repository**: Adapter for querying or manipulating data in the data layer. Change business entity into a data layer model and vice versa. Data layer to be accessed can include: in-memory, RDBMS, NoSQL, File System, External dependency to other service (internal or 3rd party) 
* **Entity**: Business object which has its own data structure and methods. Even though it can be the same as DB models but it does not necessarily have to be the same


For better visualization, current go-testapp application looks like the diagram below:
![image](https://user-images.githubusercontent.com/102520846/178428592-301e1626-f699-4d36-bb4d-269388cded07.png)


### Dependency Injection
Dependency injection states that any dependency should be provided as part of argument in initialization period. This allows us to mock and change implementation details later, as long as interface is getting implemented. Using dependency injection, the flow of initialization is reversed, starting from the leaf (no dependency)

Flow:
1. Initialize your data **layer**: External dependency, DB connection, redis connection, etc.
2. Use your initialized **data layer** to initialize **repository instance**
3. Use your initialized **repository instance** to initialize **usecase instance**
4. Use your initialized **usecase instance** to initialize **handler instance**
5. Use your initialized **handler instance** to initialize router

Therefore, if it happens that you want to change the data layer, you can do so by just changing one step (step 1) without changing anything else.

### Code Generation


Code generation is used for mundane code task that takes a lot of time. Using code generation reduce the amount of time to create additional boilerplate code. Currently, this project use 2 code generation library:
1. [GQL Server w/ gqlgen](https://github.com/99designs/gqlgen)
2. [Mocking w/ mockgen](https://github.com/golang/mock)

There are special sections for our generated code that SHOULD NOT be edited manually by developers. It should only be edited by another action from the code generator.

```
|--- gen
|    |--- mockgen
|    |    |--- handler
|    |    |    |--- rest
|    |    |    |--- grpc
|    |    |
|    |    |--- usecase
|    |    |    |--- module A
|    |    |    |--- module B
|    |
|    |--- graph
|    |    |--- generated
|    |    |    |--- *.generated.go
|    |    |
|    |    |--- model
|    |    |    |--- models_gen.go
|    |
```

### gqlgen
Graphql language by its definition is a dynamic language that is not statically typed. The transition from a dynamically typed language (GQL) to a statically typed language (GoLang) is an arduous task that is better solved by code generation.

gqlgen code generation configuration is defined by `gqlgen.yml`

There are 3 important kinds of generated files:
- `gen/graph/generated/*.generated.go`: Core logic of transforming data from GQL dynamic type to Golang static type. this file SHOULD NOT CHANGED
- `gen/graph/model/models_gen.go`: Golang struct for corresponding GQL data structure. This can be replaced with our own struct by toying with `gqlgen.yml` configuration, but generated model SHOULD NOT BE CHANGED.
- `internal/handler/gql/*resolvers.go`: Skeleton for GQL resolver, similar to handler for REST API. Developers need to implement the resolver (data parsing and calling the usecase layer)

### mockgen
Mock are used for unit testing. Mocking an interface is also a mundane task that can be delegated to a code generation tools.


From the code generation branch, there is no `repo` folder. That is because in Go, each data layer (which is the repo folder dependency) already have their own mocking mechanism.

In order to mock a new interface, first identify the file where the interface is written, then run from root:

```
mockgen -source=<internal/path/to/interface> -destination=<gen/path/to/mock_file>  
```

## Unit Test
**[Work In Progress]** For now please refer to the use of `mockgen` and existing boilerplate in `internal/usecase/order/order_test.go` on how to do the unit test.

### Making Unit Test
If you are using VSCode, you can easily create a new test case to a function by right clicking and choose to create unit test. It will generate a new file `<original_filename>_test.go` in the same directory. However, to mock you need to do some manual editing.

![image](https://user-images.githubusercontent.com/102520846/172808630-314c556a-9cdd-42fa-a543-5ed5b001a70f.png)

## CI/CD Workflow
CI/CD workflow is based on [Github Action](https://docs.github.com/en/actions). Workflow definition is on `.github/workflows`, below are summary of each workflow:
- PR workflow (`pr.yaml`): Doing unit testing & lint checking every pull request
- Push workflow (`push.yaml`): Do service docker image build, and deploy it to dev and staging environment every merge to master
- Tag workflow (`tag.yaml`): For deploying to staging / production environment, based on tag pushed to this repository, there are two kind of tag:
    1. `STG-X.Y.Z`, for deploying to staging environment, with `X.Y.Z` as semantic version
    1. `PRD-X.Y.Z`, for deploying to production environment, with `X.Y.Z` as semantic version

Following variable need to be adjusted on workflow file when deploying real service:
- `<Service Name>`: Service human readable name.
- `<service-image-registry-path>`: Docker image registry used by this repo for storing docker image build.
- `<service-rhodes-name>`: general name for service on Rhodes repository.

To implement CD workflow fully, you need to also make adjustment on `rhodes` repository of infrastructure.
`<service-image-registry-path>` and `<service-rhodes-name>` in particular need to be aligned with data on `rhodes` for this service.
Example of adjustment needed on rhodes for new service could be seen on [link](https://github.com/kargotech/rhodes/pull/482/files)

## Architecture Decision Records
Architecture Decision Records is lightweight documentation tools for recording any significant decision made in the project. 

Architecture Decision Records enable storing big decisions made in a log as a reference point for the team, 
help with onboarding new members and give context to others interested in the project.

For more details, please refer to `docs/architecture-decisions/index.md`

## Code Quality and Security Analysis
This project has code quality and security analysis performed every pull request and given as pull request feedback. 
Code quality and Security issue is tracked on SonarQube for long term view and details on strong and weak area of codebase.
Code quality is based on `golangci-lint` tools, and security analysis is based on `horusec`

For more details, please refer to `docs/CODE_QUALITY.md`

## Database Migration Flow
This project implement two kind of Database migration flow:
1. `default`: this migration is migration that is automatically applied everytime there is application deployment.
2. `postmigration` this migration triggered manually independently of application deployment, it could be triggered via Github Action `db-operation`.

For more details, please refer to `docs/DATABASE_MIGRATION.md`
## OpenTelemetry Tracing

### Initialized Tracing
Specification:
- Exporter: OTLP Collector (backend can be configured by infra via their config)
- Resource: service name via config (devel), in staging&production maybe injected via agent
- Provider: Batching

Configuration is set through environment variable with prefix OTEL_*.
Example:
```
OTEL_EXPORTER_OTLP_ENDPOINT=
```

### Gin
Implemented via gin middleware, to implement just need to call otelgin middleware with initialized tracer provider
```
router.Use(otelgin.Middleware(mainCfg.Service.Name, otelgin.WithTracerProvider(tp)))
```

### GORM
Implemented via gorm middleware, to implement just need to call otelgorm middleware with initialized tracer provider.
Option WithoutQueryVariables() is a must, especially for production to hide sensitive database value from tracing UI.
```
db.Use(otelgorm.NewPlugin(
		otelgorm.WithTracerProvider(tp),
		otelgorm.WithoutQueryVariables()))
```

To contribute to the existing Span, on any gorm call, dev need to add `WithContext(ctx)`, 
Example:
```
om.db.WithContext(ctx).Model(&order).Select("num_sales").Where("ksuid = ?", order.Ksuid).Updates(order).Error
```
