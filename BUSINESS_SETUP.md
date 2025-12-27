# Business Setup (Multi-Tenant, EC2 + RDS Postgres)

This backend now supports email/password auth, multi-tenant users, and
cross-account access via STS AssumeRole.

## Required Environment Variables

Backend:
- `DB_URL` (Postgres connection string)
- `JWT_SECRET` (random 32+ chars)
- `APP_BASE_URL` (public URL for verification links)
- `SMTP_HOST`
- `SMTP_PORT` (default 587 if omitted)
- `SMTP_USER`
- `SMTP_PASS`
- `SMTP_FROM`

Optional:
- `AWS_REGION` (default region if user settings empty)

## Initial Flow

1) User signs up via `POST /auth/signup` and receives a verification email.
2) User verifies via `GET /auth/verify?token=...`.
3) User logs in via `POST /auth/login` and receives access + refresh tokens.
4) User creates AWS accounts via `POST /aws/accounts`.
5) Frontend sends `Authorization: Bearer <token>` and `X-AWS-ACCOUNT-ID`.

## Client AWS Role (per account)

Each client creates a role in their AWS account and shares:
- `role_arn`
- `external_id`
- `account_id`

Your backend assumes this role for each request.

## Notes

- Tokens: access 15 minutes, refresh 30 days (rotated).
- Email verification required before login.
