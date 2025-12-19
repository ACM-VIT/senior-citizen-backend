#!/bin/bash

# Generate secure random secrets
generate_secret_hex() {
    openssl rand -hex 32
}

generate_secret_short() {
    openssl rand -hex 16
}

generate_cipher_secret() {
    openssl rand -base64 32 | head -c 32
}

# Generate all secrets
HYDRA_DB_PASS=$(generate_secret_hex)
KRATOS_DB_PASS=$(generate_secret_hex)
COOKIE_SECRET=$(generate_secret_hex)
CSRF_COOKIE_SECRET=$(generate_secret_hex)
KRATOS_COOKIE_SECRET=$(generate_secret_hex)
KRATOS_CIPHER_SECRET=$(generate_cipher_secret)
HYDRA_SYSTEM_SECRET=$(generate_secret_hex)
HYDRA_OIDC_SALT=$(generate_secret_short)

cat > .env << EOF
# Database
DATABASE_URL=postgresql://root@cockroach:26257/senior_citizen?sslmode=verify-full&sslrootcert=/certs/ca.crt&sslcert=/certs/client.root.crt&sslkey=/certs/client.root.key
SHOULD_MIGRATE=true

# Backend
PORT=3000
PUBLIC_URL=https://hub.acmvit.in

# Hydra
HYDRA_ADMIN_URL=http://hydra:4445
HYDRA_INTROSPECTION_CLIENT_ID=
HYDRA_INTROSPECTION_CLIENT_SECRET=
HYDRA_DSN=postgres://hydra:${HYDRA_DB_PASS}@hydra-postgres:5432/hydra?sslmode=disable
HYDRA_DB_PASSWORD=${HYDRA_DB_PASS}
HYDRA_SYSTEM_SECRET=${HYDRA_SYSTEM_SECRET}
HYDRA_OIDC_SALT=${HYDRA_OIDC_SALT}

# Kratos
KRATOS_ADMIN_URL=http://kratos:4434
KRATOS_PUBLIC_URL=http://kratos:4433
KRATOS_BROWSER_URL=https://auth.acmvit.in
KRATOS_UI_URL=https://hub.acmvit.in/ui
KRATOS_DSN=postgres://kratos:${KRATOS_DB_PASS}@kratos-postgres:5432/kratos?sslmode=disable
KRATOS_DB_PASSWORD=${KRATOS_DB_PASS}
KRATOS_LOG_LEVEL=info
KRATOS_COOKIE_SECRET=${KRATOS_COOKIE_SECRET}
KRATOS_CIPHER_SECRET=${KRATOS_CIPHER_SECRET}

# App cookies
COOKIE_SECRET=${COOKIE_SECRET}
CSRF_COOKIE_NAME=sc_csrf
CSRF_COOKIE_SECRET=${CSRF_COOKIE_SECRET}
EOF

echo "✅ Generated new .env file with fresh secrets"
echo ""
echo "⚠️  Make sure to:"
echo "   1. Back up this file securely"
echo "   2. Never commit it to git"
echo "   3. Update HYDRA_INTROSPECTION_CLIENT_ID and HYDRA_INTROSPECTION_CLIENT_SECRET if needed"
