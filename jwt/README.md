# jwt

A Gin + GORM (MySQL) blog app that does cookie-based authentication with HS256-signed JWTs. Based on Codeheim's [Episode #36: Mastering JWT Authentication in Golang](https://www.codeheim.io/courses/Episode-36-Mastering-JWT-Authentication-in-Golang-663a796d51be3a2c1c6d0d6b).

## What the code does

- `main.go`: starts Gin on the default port 8080. It loads the HTML templates from `templates/**/*`, connects to MySQL, auto-migrates the tables, and registers these routes:
  - `GET /signup`, `GET /login`: render the forms.
  - `POST /signup`, `POST /login`: read the `email` and `password` form fields.
  - `DELETE /logout`: clears the cookie.
  - `GET /blogs`, `GET /blogs/:id`: protected by `AuthMiddleware`.
- `models/setup.go`: connects to `gin_jwt:tmp_pwd@tcp(127.0.0.1:3306)/gin_jwt` and panics if the connection fails. `DBMigrate` creates the `blogs` and `users` tables.
- `models/user.go`: hashes passwords with bcrypt on signup and checks them with `bcrypt.CompareHashAndPassword` on login.
- `models/blog.go`: `BlogsAll` lists non-deleted blogs, newest `updated_at` first. `BlogsFind` loads one blog by ID.
- `controllers/sessions_controller.go`: on a successful signup or login, it signs a JWT with the claims `{"userID": <id>, "ttl": now+100 days}`. The token goes into an HttpOnly, SameSite=Lax cookie named `Auth` (max-age 100 days), and the user is redirected to `/blogs`. Logout overwrites the cookie with `"deleted"` and redirects to `/login`.
- `middlewares/auth.go`: reads the `Auth` cookie and parses the token. It rejects any signing method that isn't HMAC, checks the custom `ttl` claim against the current time, and loads the user by `userID`. The user is stored in the context under `"user"`, which `BlogsIndex` uses to show "Logged in as ...".
- `templates/`: Bootstrap-styled pages. The Logout link on the blog index sends a `DELETE /logout` with XHR.

The app has no route that creates blog posts. You have to insert rows into `blogs` yourself (for example, with SQL).

## Concepts you can learn

- **Stateless sessions with a JWT in a cookie**: the server keeps no session store. Each request is authenticated by checking the token's signature.
- **Validating the `alg` in the key function**: the code rejects any token whose method isn't `*jwt.SigningMethodHMAC`. This prevents algorithm-confusion attacks.
- **HttpOnly + SameSite cookies**: JavaScript can't read the token, and the browser won't send it on most cross-site requests.
- **bcrypt password hashing**: only salted, slow hashes are stored. Plaintext passwords never are.
- **Gin middleware chains**: `c.Set` / `c.Get` pass the authenticated user to the handlers, and `c.Next` / `c.Abort` control the request flow.

## How to run

Prerequisite: a MySQL server on `127.0.0.1:3306` with database `gin_jwt`, user `gin_jwt`, and password `tmp_pwd`. These values are hard-coded in `models/setup.go`.

```bash
cd jwt
go run .
# open http://localhost:8080/signup
```

Run it from inside `jwt/`, because the template glob and `./public` are relative paths.

## Notes

- **The middleware doesn't `return` after `c.Abort()`.** It keeps running after an auth failure. For example, with no cookie, `jwt.Parse("")` fails and `token.Claims` is then read from a nil token, which panics (Gin's recovery middleware catches it). Every failure branch needs a `return`.
- Expiry uses a custom `ttl` claim instead of the standard `exp` claim, so the jwt library's built-in expiry check doesn't apply. The token can last 100 days.
- The HMAC secret is hard-coded in both `sessions_controller.go` and `auth.go` (both have a TODO to move it to an env var). The cookie is set with `secure=false`.
- Signup shows "Email missing" when the email is already taken.
- `r.Static("/public", "./public")` points to a directory that doesn't exist. The `logout()` JS helper in `header.tpl` calls `/users/logout`, which isn't a registered route. Nothing uses that helper.
- `BlogsShow` logs an invalid `:id` and continues with ID 0.
- `go.mod` declares `go 1.22.0` and marks directly used packages such as `golang-jwt/jwt/v5` as `// indirect`.
- The directory has its own Apache-2.0 `LICENSE`.
