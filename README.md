# Senior Citizen Backend – API Documentation

This document provides the API reference for the **senior-citizen-backend** service.  
The backend uses:

- **Ory Hydra** → for OAuth2 / OIDC authentication and token management.
- **Ory Kratos** → for user identity, self-service flows (registration, login, recovery), and profile data.

The service exposes REST endpoints for:

- User management
- Profile creation and updates
- Admin management
- OAuth2 integration flows (Hydra + Kratos)

---

## 1. Register Client

**Endpoint:**  
`POST http://127.0.0.1:4445/clients`

**Description:**  
Register a new OAuth2 client with Hydra.

**Request Body:**

```json
{
  "client_id": "reverse-coding",
  "client_secret": "rc-secret",
  "grant_types": ["authorization_code", "refresh_token"],
  "response_types": ["code", "id_token", "token"],
  "scope": "openid email profile offline",
  "redirect_uris": ["http://localhost:3000/callback"],
  "token_endpoint_auth_method": "client_secret_post",
  "access_token_strategy": "jwt"
}
```

---

## 2. Get Tokens

**Step 1 – Authorization Request**

Open the following URL in a browser to start the OAuth2 flow:

```http
GET http://127.0.0.1:4444/oauth2/auth?client_id=reverse-coding&response_type=code&scope=openid%20email%20profile%20offline&redirect_uri=http://localhost:3000/callback&state=<RANDOM_STATE>
```

- Replace `<RANDOM_STATE>` with a random string for CSRF protection.
- After login/consent, you’ll be redirected to:

  ```
  http://localhost:3000/callback?code=<AUTH_CODE>&state=<RANDOM_STATE>
  ```

**Step 2 – Token Request**

**Endpoint:**
`POST http://127.0.0.1:4444/oauth2/token`

**Description:**
Exchange the authorization code for tokens.

**Request Body (x-www-form-urlencoded):**

```http
grant_type=authorization_code
code=<AUTH_CODE>
redirect_uri=http://localhost:3000/callback
client_id=reverse-coding
client_secret=rc-secret
```

---

## 3. Get Current User

**Endpoint:**
`POST /api/v1/user/get`

**Description:**
Get details of the authenticated user.

**Request Headers:**

```http
Authorization: Bearer <ACCESS_TOKEN>
```

---

## 4. Get All Members

**Endpoint:**
`POST /api/v1/members/get`

**Description:**
Fetch details of all members.

**Request Headers:**

```http
Authorization: Bearer <ACCESS_TOKEN>
```

---

## 5. Create or Update Profile

**Endpoint:**
`POST /api/v1/profile/create`

**Description:**
Create or update a user profile.

**Request Headers:**

```http
Authorization: Bearer <ACCESS_TOKEN>
```

**Request Body:**

```json
{
  "user_id": "user-id",
  "name": "John Doe",
  "profile_picture_url": "https://example.com/pic.png",
  "phone": "9876543210",
  "gender": "male",
  "bio": "I love programming",
  "the_sentence": "Learning never stops.",
  "linkedin": "https://linkedin.com/in/johndoe",
  "github": "https://github.com/johndoe",
  "twitter": "https://twitter.com/johndoe",
  "instagram": "https://instagram.com/johndoe",
  "job_title": "Software Engineer",
  "company": "Tech Corp",
  "website": "https://johndoe.dev"
}
```

---

## 6. Create Admin

**Endpoint:**
`POST /api/v1/admin/create`

**Description:**
Create or update an admin record.

**Request Headers:**

```http
Authorization: Bearer <ACCESS_TOKEN>
```

---

## 7. Export CSV

**Endpoint:**
`GET /api/v1/csv/export`

**Description:**
Export all users as CSV.

**Request Headers:**

```http
Authorization: Bearer <ACCESS_TOKEN>
```

---

## 8. Get or Create User

**Endpoint:**
`GET /api/v1/user/create?email=<EMAIL>`

**Description:**
Get or create a user by email.

**Request Headers:**

```http
Authorization: Bearer <ACCESS_TOKEN>
```

**Example:**

```http
GET /api/v1/user/create?email=johndoe@example.com
```

---

## 9. Protected Route

**Endpoint:**
`GET /api/v1/auth/protected`

**Description:**
A test route that requires authentication.

**Request Headers:**

```http
Authorization: Bearer <ACCESS_TOKEN>
```

---

## 10. OAuth2 Callback

**Endpoint:**
`GET /api/v1/auth/callback`

**Description:**
Handles the OAuth2 callback from the authorization server.
Exchanges the authorization code for tokens based on the `app` identifier.

**Query Parameters:**

| Name | Type   | Required | Description                                              |
| ---- | ------ | -------- | -------------------------------------------------------- |
| app  | string | Yes      | App identifier (`reverse-coding`, `cryptic-hunt`, `c2c`) |
| code | string | Yes      | Authorization code returned by `/oauth2/auth`            |

**Example Request:**

```http
GET http://127.0.0.1:3000/api/v1/auth/callback?app=reverse-coding&code=abc123
```

---

## 11. OAuth2 Consent Flow

**Endpoint:**  
`GET /api/v1/authconsent`

**Description:**  
Handles the OAuth2 **consent challenge** from Hydra.

- Fetches the consent request from **Hydra**
- Looks up the user identity in **Kratos**
- Accepts the consent request with granted scopes (`openid`, `profile`, `email`)
- Redirects back to Hydra’s authorization flow

**Query Parameters**

| Name              | Type   | Required | Description                                |
| ----------------- | ------ | -------- | ------------------------------------------ |
| consent_challenge | string | Yes      | Consent challenge string provided by Hydra |

**Example Request**

```http
GET http://127.0.0.1:3000/api/v1/authconsent?consent_challenge=abc123
```

---

## 11. Update Profile

**Endpoint:** `PUT /api/v1/profile/update`
**Description:** Update the profile information of the authenticated user.

**Request Headers**

```http
Authorization: Bearer <ACCESS_TOKEN>
Content-Type: application/json
```

**Request Body**

```json
{
  "name": "John Doe",
  "profile_picture_url": "https://example.com/pic.png",
  "phone": "9876543210",
  "gender": "male",
  "bio": "I love programming",
  "the_sentence": "Learning never stops.",
  "linkedin": "https://linkedin.com/in/johndoe",
  "github": "https://github.com/johndoe",
  "twitter": "https://twitter.com/johndoe",
  "instagram": "https://instagram.com/johndoe",
  "job_title": "Software Engineer",
  "company": "Tech Corp",
  "website": "https://johndoe.dev"
}
```

**Note:**

- Only include the fields you want to update.
- Fields not included in the request body will remain unchanged.

---

## 12. Update Roles

**Endpoint:** `POST /api/v1/roles/update`  
**Description:** Create or update roles and positions for users in a specific tenure. Only **admin users** can access this route.

**Request Headers**

```http
Authorization: Bearer <ACCESS_TOKEN>
Content-Type: application/json
```

**Request Body**

```json
[
  {
    "email": "johndoe@example.com",
    "role": "member",
    "position": "secretary",
    "tenure": "2025"
  },
  {
    "email": "janedoe@example.com",
    "role": "admin",
    "position": "chair",
    "tenure": "2025"
  }
]
```

**Notes:**

- `email` must correspond to an existing user.
- `role` and `position` must be valid (checked via helpers).
- `tenure` must exist.

---

## 13. Get Admin

**Endpoint:**
`POST /api/v1/admin/get`

**Description:**
Check if the authenticated user is an admin.

**Request Headers:**

```http
Authorization: Bearer <ACCESS_TOKEN>
```

### **Admin Routes**

1. Create admin
2. Export CSV
3. Get or create user
4. Update Roles
