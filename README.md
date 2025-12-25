# Senior Citizen Backend

Go API service for the senior-citizen platform. It uses:

- Hydra for OAuth2/OIDC
- Kratos for identity and self-service flows
- CockroachDB for application data
- Fiber + GORM for the API
- GIN indexes and cockroachdb changefeeds for fuzzy search

## Local development (Docker Compose)

1. Copy `.env.example` to `.env` and update secrets.
2. Start the stack:

```sh
docker compose up --build
```

If `docker compose` is not available, use:

```sh
docker-compose up --build
```

### Configuration

Key env vars (see `.env.example`):

- `DATABASE_URL` - CockroachDB connection string.
- `PUBLIC_URL` - public base URL of this API (used for login return_to).
- `HYDRA_ADMIN_URL` - Hydra admin API URL (internal).
- `HYDRA_INTROSPECTION_CLIENT_ID` / `HYDRA_INTROSPECTION_CLIENT_SECRET` - optional client creds for token introspection.
- `KRATOS_ADMIN_URL` - Kratos admin API URL (internal).
- `KRATOS_PUBLIC_URL` - Kratos public URL (internal).
- `KRATOS_BROWSER_URL` - Kratos public URL for the browser.
- `KRATOS_UI_URL` - Auth UI base URL (browser-facing).

Cockroach runs in secure mode in Docker; certs are auto-generated and mounted at `/certs`. The default `DATABASE_URL` in `.env.example` uses those certs. If you run the API outside Docker, you’ll need to export the certs and update the path accordingly.

Default ports:

- API: `http://127.0.0.1:3000`
- Hydra public: `http://127.0.0.1:4444`
- Hydra admin: `http://127.0.0.1:4445`
- Kratos public: `http://127.0.0.1:4433`
- Kratos admin: `http://127.0.0.1:4434`
- Auth UI: `http://127.0.0.1:3000/ui`

If you run the API locally outside Docker, set `HYDRA_ADMIN_URL` and `KRATOS_ADMIN_URL` to `http://127.0.0.1:4445` and `http://127.0.0.1:4434`.
For production or custom domains, update `config/hydra.yml` and `config/kratos.yml` URLs to match your public auth domain.

## Auth and SSO

The backend acts as a resource server. Clients should run the OAuth2 authorization code flow directly against Hydra. The API expects a `Bearer` token from Hydra and validates it via the Hydra admin introspection endpoint.

Hydra login is handled by the backend at:

```
GET /api/v1/auth/login?login_challenge=<challenge>
```

The login handler checks for an existing Kratos session and redirects to the auth UI if needed, then accepts the Hydra login request.

### Register a client

```sh
curl -X POST http://127.0.0.1:4445/admin/clients \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": "reverse-coding",
    "client_secret": "rc-secret",
    "grant_types": ["authorization_code", "refresh_token"],
    "response_types": ["code", "id_token", "token"],
    "scope": "openid email profile offline",
    "redirect_uris": ["http://localhost:3000/callback"],
    "token_endpoint_auth_method": "client_secret_post",
    "access_token_strategy": "jwt"
  }'
```

### Authorization flow

1. Authorization request (browser):

```http
GET http://127.0.0.1:4444/oauth2/auth?client_id=reverse-coding&response_type=code&scope=openid%20email%20profile%20offline&redirect_uri=http://localhost:3000/callback&state=<RANDOM_STATE>
```

2. Token exchange (server-side client):

```http
POST http://127.0.0.1:4444/oauth2/token
grant_type=authorization_code
code=<AUTH_CODE>
redirect_uri=http://localhost:3000/callback
client_id=reverse-coding
client_secret=rc-secret
```

### Consent handling

Hydra redirects consent requests to:

```
GET /api/v1/authconsent?consent_challenge=<challenge>
```

The backend accepts the consent request, pulls identity traits from Kratos, and returns Hydra’s redirect.

## Next.js apps (NextAuth/better-auth)

Use Hydra as your OIDC provider in each app. Users will have a single IdP session on the auth domain; each app keeps its own session cookie but SSO is seamless because Hydra already has a session.

For cross-subdomain SSO:

- Set Hydra and Kratos public URLs to `https://auth.example.com`.
- Set `PUBLIC_URL=https://auth.example.com` and `KRATOS_UI_URL=https://auth.example.com/ui` (or your UI domain).
- Configure Kratos cookie domain to `.example.com` (see Kratos config docs).
- Register each app redirect URI in Hydra.

Local dev tip: use `*.lvh.me` or `*.localtest.me` to test subdomain cookies.

## API Endpoints

All endpoints require `Authorization: Bearer <ACCESS_TOKEN>` unless noted.

### Users and Admin

- `GET /api/v1/user/get` - get current user (includes profile).
- `GET /api/v1/user/create?email=<EMAIL>` - get or create a user by email (admin only).
- `GET /api/v1/admin/get` - check if current user is admin.
- `POST /api/v1/admin/create` - create admin for a user (admin only). Body:
  ```json
  { "user_id": "<UUID>" }
  ```
  or
  ```json
  { "email": "user@example.com" }
  ```

### Auth Flow

- `GET /api/v1/auth/login?login_challenge=<challenge>` - Hydra login handler.
- `GET /api/v1/authconsent?consent_challenge=<challenge>` - Hydra consent handler.
- `GET /callback` - dev helper to display the auth code and error params.

### Profiles

- `POST /api/v1/profile/create` - create or update the current user profile.
- `PUT /api/v1/profile/update` - partial update of the current user profile.

### Members and Roles

- `GET /api/v1/members/get?role=&position=&page=&page_size=` - list memberships.
- `POST /api/v1/search/member` - search memberships.
- `POST /api/v1/roles/update` - update roles (admin only).
- `GET /api/v1/csv/export` - export member emails as CSV (admin only).

### Tenures

- `GET /api/v1/tenures?current=true&page=&page_size=`
- `GET /api/v1/tenures/:id`
- `POST /api/v1/tenures` (admin only)
- `PUT /api/v1/tenures/:id` (admin only)
- `DELETE /api/v1/tenures/:id` (admin only)

### Events

- `GET /api/v1/events?tenure_id=&event_type=&page=&page_size=`
- `GET /api/v1/events/:id`
- `POST /api/v1/events` (admin only)
- `PUT /api/v1/events/:id` (admin only)
- `DELETE /api/v1/events/:id` (admin only)

Create/update accepts `sponsor_ids` to attach sponsors.
Valid `event_type` values: `online`, `offline`, `hybrid`.

### Sponsors

- `GET /api/v1/sponsors?page=&page_size=`
- `GET /api/v1/sponsors/:id`
- `POST /api/v1/sponsors` (admin only)
- `PUT /api/v1/sponsors/:id` (admin only)
- `DELETE /api/v1/sponsors/:id` (admin only)

### Participations

- `GET /api/v1/participations?event_id=&page=&page_size=`
- `GET /api/v1/participations/:id`
- `POST /api/v1/participations`
- `PUT /api/v1/participations/:id`
- `DELETE /api/v1/participations/:id`

Non-admin users can only view/update/delete their own records. Admins can act on any user and set `user_id` in requests.

### Winners

- `GET /api/v1/winners?event_id=&page=&page_size=`
- `GET /api/v1/winners/:id`
- `POST /api/v1/winners` (admin only)
- `PUT /api/v1/winners/:id` (admin only)
- `DELETE /api/v1/winners/:id` (admin only)

### Email Preferences

- `GET /api/v1/email-preferences`
- `PUT /api/v1/email-preferences`

### Fuzzy Search

Since pg_trgm extension is not supported in CockroachDB, we implement fuzzy search using trigram generation, GIN indexes, and changefeeds.

1. Run the docker-compose file with `SHOULD_MIGRATE=true` to create the GIN indexes.
2. Parallely run go run ./cdc to start the changefeed listener.
3. Forward port 8081 to expose the changefeed listener.
4. Create a changefeed for the table you want to fuzzy search. Currently enabled on profiles table on name and company fields.
   `CREATE CHANGEFEED FOR TABLE profiles INTO 'webhook-https://cd6b20fcb1eb.ngrok-free.app/cdc' WITH OPTIONS (envelope = 'wrapped',format = 'json', initial_scan, updated)`
5. The above can be configured either in cockroachdb sql shell or with the help of cockroachdb cli.
6. The webhook is triggered whenver a row is inserted / updated in the table.
7. To search for a profile, send a POST request to `/api/v1/search/member`. Before you search ensure that the membership table is populated.
