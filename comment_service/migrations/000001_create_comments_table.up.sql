create table "comments"(
    "id" uuid, 
    "post_id" text, 
    "user_id" text, 
    "text" varchar, 
    "created_at" timestamp default current_timestamp, 
    "deleted_at" timestamp
);