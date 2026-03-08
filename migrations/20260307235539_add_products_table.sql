-- Create "products" table
CREATE TABLE "products" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "price" numeric NOT NULL,
  "user_id" bigint NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_users" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
