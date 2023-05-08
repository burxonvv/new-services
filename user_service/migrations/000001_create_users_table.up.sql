CREATE TABLE IF NOT EXISTS "users" (
    "id" uuid,
    "first_name" TEXT,
    "last_name" TEXT,
    "email" TEXT,
    "created_at" TIME DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIME DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIME,
    "password" text
)