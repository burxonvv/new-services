create table "posts"(
    "id" uuid, 
    "user_id" uuid, 
    "title" text, 
    "description" text, 
    "likes" integer default 0, 
    "created_at" timestamp default current_timestamp, 
    "updated_at" timestamp default current_timestamp, 
    "deleted_at" timestamp
);