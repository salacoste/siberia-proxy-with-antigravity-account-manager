# Configuration Guide

## Environment Variables

The backend uses environment variables for Cloud Service integration.

### Setup
1.  Copy `.env.example` to `.env` in `apps/backend/`.
    ```bash
    cp apps/backend/.env.example apps/backend/.env
    ```
2.  Fill in your credentials.

### Variables

| Variable | Description |
| :--- | :--- |
| `SUPABASE_URL` | URL of your Supabase instance. |
| `SUPABASE_KEY` | Public Anon key (or Service Role if needed). |
| `SUPABASE_DB_PASSWORD` | Database password (only used if directly connecting to DB, mostly unused by App currently). |
| `MINIO_*` | MinIO connection details (see `.env.example`). |

> **Warning**: Never commit your `.env` file to version control.
