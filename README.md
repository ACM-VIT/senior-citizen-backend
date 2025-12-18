# Senior Citizen Backend

Go API service for the senior-citizen platform. It uses:

- Ory Hydra for OAuth2/OIDC
- Ory Kratos for identity and self-service flows
- CockroachDB for application data
- Fiber + GORM for the API

## Local development (Docker Compose)

1. Copy `.env.example` to `.env` and update secrets.
2. Start the stack:

```sh
docker compose up --build
```

Default ports:

- API: `http://127.0.0.1:3000`
- Hydra public: `http://127.0.0.1:4444`
- Hydra admin: `http://127.0.0.1:4445`
- Kratos public: `http://127.0.0.1:4433`
- Kratos admin: `http://127.0.0.1:4434`
- Kratos UI: `http://127.0.0.1:4455`

If you run the API locally outside Docker, set `HYDRA_ADMIN_URL` and `KRATOS_ADMIN_URL` to `http://127.0.0.1:4445` and `http://127.0.0.1:4434`.

## Auth and SSO

The backend acts as a resource server. Clients should run the OAuth2 authorization code flow directly against Hydra. The API expects a `Bearer` token from Hydra and validates it via the Hydra admin introspection endpoint.

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

1) Authorization request (browser):

```http
GET http://127.0.0.1:4444/oauth2/auth?client_id=reverse-coding&response_type=code&scope=openid%20email%20profile%20offline&redirect_uri=http://localhost:3000/callback&state=<RANDOM_STATE>
```

2) Token exchange (server-side client):

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
