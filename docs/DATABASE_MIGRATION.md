# Database Migration Flow
Adopting [RFC DB Migration Flow](https://kargox.atlassian.net/wiki/spaces/ENG/pages/2434269209/RFC+-+DB+schema+migration+flow), 
this project implement two kind of migration via `dbmate`, each with having their own migration directory and migration tracker:
1. `default`: `default` migration is migration that is automatically applied everytime application deploy. 
`default` migration is usually migration which are critical for application functionality, such as adding column in a table
`default` migration should be fast and generally finish under 1 minute, and use `lock_timeout` to ensure it finish under 1 minute.

2. `postdeployment`: `postdeployment` migration is migration which tend to take long time and less critical to application functionality.
`postdeployment` migration is be triggered manually, usually during business off-hour, as to minimize impact of migration locking.
example of migration that is usually under `postdeployment` is Adding index to a heavily populated table.

## Running manual migration
### Setting relevant secrets and environment variable
Manual DB Migration workflow require following repository/action secrets to be set:
1. `DBMATE_DATABASE_URL_${ENV}`, according to the ENV that it is applied, 
for example for applying to `prod` database, `DBMATE_DATABASE_URL_PROD` need to be set

Beside that, there's also optional environment variable for application, `AUTO_MIGRATE_POSTDEPLOYMENT_MIGRATIONS`, 
which if set to `true` would automatically run `postdeployment` migration every application startup. 
This is useful in case of `dev` and `integration` environment, which could be set either via Kubernetes Secret or Rhodes configuration. 
For `stg` and `prod` `AUTO_MIGRATE_POSTDEPLOYMENT_MIGRATIONS` should be false/not setted.

### Action
Manual migration is done via Github Action with Workflow name: `DBMate Operation`. 
After selecting `DBMate Operation` Workflow on Action page, engineer could select `Run workflow` and fill in necessary variable for it:
`environment`, `command`, `migration_type`.

`command` currently support three command: `status` for checking status of migration tracker, `migrate` to migrate all migration, and `rollback` to rollback single migration

By default only people who have write access to repository which are able to trigger manual `DBMate Operation` workflow. 
If more granular authorization needed, engineer could edit `.github/workflows/db-operation.yaml` and add filter on workflow to enable only select author (`github.actor` variable from github action).

## Guide on migration
This section give general recommendation which type of migration should be in `default`, and which should be in `postdeployment`

### Safeguard & Time consideration
Generally, a single `default` migration finish under 1 minute, while `postdeployment` migration finish under 20 minute.
As a safeguard, both `default` migration and `postdeployment` migration should set lock timeout at start of their migration.

```
-- For default migration
SET LOCAL lock_timeout = '60s';

-- For postdeployment migration
SET LOCAL lock_timeout = '1200s';
```

### Usual operation

#### Creating table / column / Foreign key
Creating table is safe to put under `default` migration.

Creating column without default value is safe to put under `default` migration.

Creating column with default value is safe to put under `default` migration, for PostgreSQL 12+, for lower version of PostgreSQL it might not be safe.

Adding a foreign key column is recommended to be separated into two migration, first creating foreign key constrain without validation, and second validation of foreign key.

### Altering / renaming column
Some changing type of column is safe, depending on database, for example on PostgreSQL 12+, following are safe:
1. increasing length on varchar or removing the limit
2. changing varchar to text
3. changing text to varchar with no length limit
4. Postgres 9.2+ - increasing precision (NOTE: not scale) of decimal or numeric columns. eg, increasing 8,2 to 10,2 is  safe. Increasing 8,2 to 8,4 is not safe.
5. Postgres 9.2+ - changing decimal or numeric to be  unconstrained
6. Postgres 12+ - changing timestamp to timestamptz when session TZ is UTC 

In case changing type of column is not safe, it is recommended to do multiple deployment, consisting of:
1. Create a new column (in `default` migration)
2. In application code, write to both columns
3. Backfill data from old column to new column (in `postdeployment` migration)
4. In application code, move reads from old column to the new column
5. In application code, remove reference old column.
6. Drop the old column (in `default` migration).

In renaming column, engineer could just change the reference in application code.
In case it's needed to actually change DB schema column name, it is recommended to do multiple deployment, consisting of:
1. Create a new column (in `default` migration)
2. In application code, write to both columns
3. Backfill data from old column to new column (in `postdeployment` migration)
4. In application code, move reads from old column to the new column
5. In application code, remove reference old column.
6. Drop the old column (in `default` migration).

### Removing column / table
Removing column or table should be separated into two step:
1. Remove any reference of column / table from application code, then deploy it into production
2. Remove column or table via migration, this can be in `default` migration

### Adding check / additional constrain
Adding check constraint should be done in two separate steps:
1. Create the constraint without validating constrain (this could be in `default` migration)
2. Validate added constraint (this could be in `default` migration)


### Data transformation
Small / fast data transformation could be put into `default` migration.

Big data transformation (< 20 minute) could be put into `postdeployment` migration

Really big data transformation (> 20 minute) should be made as batched migration job by application code.

### Adding index
Adding index on small amount of data (~1,000,000 row) could be put into `default` migration.

Adding index on large amount of data (> 1,000,000 rows) should be put into `postdeployment` migration, with concurrent index creation and disabling migration transaction and migration lock.

### Exception
It's possible for migration listed under `Usual Operation` finished faster than 1 minute, especially if it's operating with small amount of data.
Under such case, it's okay for operation which is under `postdeployment` under `Usual Operation` to be placed on `default`.

It's possible that migration which is critical to service / product functionality or have critical impact on service performance, 
but cannot be made into evolutionary via combination of Feature flag and `postdeployment` migration. In such case it's acceptable for the migration to be placed on `default` migration, with appropriate adjustment of `lock_timeout`.

### Reference
For better understanding of when to put migration in `default` and when to put in `postdeployment` and additional pattern, engineer could refer to following sources:
[Gitlab Database Guide](https://docs.gitlab.com/ee/development/database/)
[Safe Ecto Migration Guide](https://fly.io/phoenix-files/safe-ecto-migrations/)

## Folder Structure
Following relevant file & folder structure in go-testapp:
1. `db/postgres/dbmate/migrations` consist of migration for `default` migration
2. `db/postgres/dbmate/postdeployment_migrations` consist of migration for `postdeployment` migration
3. `files/deployment/db-operation.sh` consist of helper script for DBMate Operation workflow
4. `files/deployment/entrypoint.sh` consist of helper script for application entrypoint, where auto migration of `default` migration happen.
5. `.github/workflows/db-operation.yaml` workflow definition of DBMate Operation workflow.
6. `files/etc/dbmate_config` consist of default config for dbmate. for both `default` and `postdeployment` migration.

### Migrating from old go-testapp application
Engineer should copy/adjust following files:
1. `db/postgres/dbmate/migrations`
2. `db/postgres/dbmate/postdeployment_migrations`
3. `files/deployment/db-operation.sh`
4. `files/deployment/entrypoint.sh`
5. `.github/workflows/db-operation.yaml`
6. `files/etc/dbmate_config`

Also, for older version of `Dockerfile`, engineer should adjust it so it's using `entrypoint.sh` and it install `dbmate` and copy relevant migrations file (`db/`) into build artifact.
