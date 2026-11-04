# gorm

Using GORM (`gorm.io/gorm`) with PostgreSQL: auto-migration and seeding, has-one, has-many, and many-to-many associations, preloading, joins, scopes, built-in field serializers, and a custom serializer.

## What the code does

The module is named `learn`. Every program connects with the same hard-coded DSN: `host=localhost user=postgres password=wang900115 dbname=learn port=5432 sslmode=disable`.

### `model/` (package `model`, no `main`)
- `example1.go` - `User` has many `Notes` and has one `CreditCard` (foreign key `UserID`).
- `example2.go` - `Movie` and `Actor` (both embed `gorm.Model`), many-to-many through the join table `filmographs`.
- `example3.go` - `User1` has many `Order` (`User1ID`, `OrderTime`, `PaymentMode`, `Price`).
- `example4.go` - `Product` with built-in serializers: `Tags []string` as `serializer:json` stored as bytes, `Spec` as JSON, `SpecGob` as `serializer:gob`, and `CreatedTime int64` as `serializer:unixtime;type:time`.
- `example5.go` + `yaml_serializer.go` - `User2.UserSettings` is a `YamlMap`. Its `Scan`/`Value` methods implement GORM's `schema.SerializerInterface` and store the map as YAML text.

### `migration/` (one `main` per file; run each first to create and seed tables)
- `example1.go` - migrates `User`/`Note`/`CreditCard` and inserts John Doe, Jane Smith, and Alice Johnson, each with 2 notes and a card (John's card starts with `1234`).
- `example2.go` - migrates `Movie`/`Actor` and inserts 5 movies and 5 actors, **without linking them**.
- `example3.go` - migrates `User1`/`Order` and inserts Alice, Bob, and Charlie with CreditCard/Cash orders.
- `example4.go` - migrates `Product` and inserts one "Apple iPhone 13" product.
- `example5.go` - migrates `User2` and inserts one user with settings `{theme: dark, notifications: true}`.

### `example1/` - has-one / has-many queries
- `01.go`: `Find` all users, notes, and credit cards separately.
- `02.go`: `Preload("CreditCard").Preload("Notes")`.
- `03.go`: preload notes for `First` user named "John Doe".
- `04.go`: `Joins` on `credit_cards` + `Where number LIKE '1234%'`.
- `05.go`: notes where `user_id = 1`, ordered `title DESC`.

### `example2/` - many-to-many
- `01.go`: reads the join table manually through `model.Filmograph`.
- `02.go`: `Preload("Actors")` for "Titanic".
- `03.go`: `Preload("Movies")` for "Leonardo DiCaprio".

### `example3/` - scopes
- `01.go`: `Scopes(CreditCardOrders, PriceGreaterThan30)` on orders.
- `02.go`: a parameterized scope `UserFromDomain("example.com")` (`email LIKE '%example.com'`), plus `Preload("Orders", CreditCardOrders)` to filter the preloaded rows.

### `example4/` and `example5/`
- `example4/01.go`: loads product id 1 and prints the deserialized `Tags`, `Spec`, `SpecGob`, and `CreatedTime`.
- `example5/01.go`: loads `User2` id 1 and prints `UserSettings`.

## Concepts you can learn

- **AutoMigrate**: creates tables, columns, and join tables from struct definitions. It is convenient for experiments, but it is not a versioned migration tool.
- **Association conventions**: `UserID` / `User1ID` foreign keys and `many2many:filmographs` are inferred from names and tags. Nested `Create` inserts the children too.
- **Preload vs Joins**: `Preload` runs one extra query per association (avoiding N+1 loops like in `example2/01.go`). `Joins` filters the parent rows in SQL. `example1/04.go` combines both.
- **Scopes**: reusable `func(*gorm.DB) *gorm.DB` query fragments, including parameterized ones (closures), and as conditions on `Preload`.
- **Serializers**: store Go slices and maps as JSON, gob, or unix time in a single column, and define your own by implementing `Scan`/`Value` (`YamlMap`).

## How to run

Own `go.mod` (module `learn`). Each subdirectory holds several files that each have a `func main`, so run one file at a time from inside `gorm/`:

```bash
cd gorm
go run migration/example1.go   # create + seed tables first
go run example1/01.go
go run example1/04.go
go run migration/example3.go
go run example3/02.go
```

Prerequisite: PostgreSQL on `localhost:5432` with user `postgres`, password `wang900115`, and an existing database named `learn`. To use other credentials, edit the DSN in each file. For example, with Docker:

```bash
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=wang900115 -e POSTGRES_DB=learn postgres
```

## Notes

- **`example2/01.go` does not compile**: it uses `model.Filmograph`, which is not defined in `model/`.
- `migration/example2.go` never links movies and actors, so the `filmographs` table stays empty, and `example2/02.go` and `03.go` print no actors or movies. `02.go` also labels Titanic's actors "Inception actor:".
- Re-running a migration inserts duplicate rows (there are no unique constraints). Queries like `user_id = 1` and `id = 1` assume a fresh database.
- If the connection fails, `initializeDB` prints the error and returns `nil`, so the next call panics with a nil pointer. Query errors (`.Error`) are never checked in the `example*` programs.
- The DSN, including the password, is hard-coded and duplicated in every file.
- `YamlMap.Scan` only accepts `[]byte`, but pgx's `database/sql` driver returns `text` columns as `string`, so `Scan` fails with "source data not bytes". Because `example5/01.go` never checks `.Error`, it prints `map[]` instead of the settings. Accepting `string` too (or using `type:bytea`) would fix it.
- `type:time` on `CreatedTime` is Postgres's time-of-day type; whether it round-trips was not verified (it likely loses the date or fails to scan).
- Not run against a database here. `go vet` on each file passes except `example2/01.go`.
