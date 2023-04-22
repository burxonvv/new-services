CREATE TABLE IF NOT EXISTS "comments" (
    "id" SERIAL PRIMARY KEY,
    "post_id" INTEGER,
    "text" TEXT,
    "created_at" TIME DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIME
)