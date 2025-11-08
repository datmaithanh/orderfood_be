SET TIME ZONE 'Asia/Ho_Chi_Minh';

ALTER DATABASE neondb SET timezone TO 'Asia/Ho_Chi_Minh';

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "hash_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "role" varchar NOT NULL DEFAULT 'staff',
  "email" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "customers" (
  "id" bigserial PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "phone_number" varchar(15) UNIQUE NOT NULL,
  "email" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "menus" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "price" numeric(10,2) NOT NULL,
  "category_id" bigint NOT NULL,
  "status" bool NOT NULL DEFAULT 'false',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "tables" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "qr_text" varchar(255) NOT NULL DEFAULT '',
  "qr_image_url" varchar(255) NOT NULL DEFAULT '',
  "status" varchar(20) NOT NULL DEFAULT 'available',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "customer_id" bigint NOT NULL,
  "table_id" bigint NOT NULL,
  "status" varchar(20) NOT NULL DEFAULT 'pending',
  "total_price" numeric(10,2) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "order_item" (
  "id" bigserial PRIMARY KEY,
  "order_id" bigint NOT NULL,
  "menu_id" bigint NOT NULL,
  "quantity" int NOT NULL,
  "price" numeric(10,2) NOT NULL,
  "note_item" varchar(255)NOT NULL DEFAULT '',
  "status" varchar(20) NOT NULL DEFAULT 'pending',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "payments" (
  "id" bigserial PRIMARY KEY,
  "order_id" bigint NOT NULL,
  "amount" numeric(10,2) NOT NULL,
  "payment_method" varchar(20) NOT NULL,
  "status" varchar(20) NOT NULL DEFAULT 'pending',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "users" ("role");

CREATE INDEX ON "customers" ("phone_number");

CREATE INDEX ON "categories" ("name");

CREATE INDEX ON "menus" ("category_id");

CREATE INDEX ON "menus" ("name");

CREATE INDEX ON "tables" ("name");

CREATE INDEX ON "orders" ("user_id");

CREATE INDEX ON "orders" ("customer_id");

CREATE INDEX ON "orders" ("table_id");

CREATE INDEX ON "order_item" ("order_id");

CREATE INDEX ON "order_item" ("menu_id");

CREATE INDEX ON "payments" ("order_id");

CREATE INDEX ON "payments" ("status");

ALTER TABLE "menus" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("customer_id") REFERENCES "customers" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("table_id") REFERENCES "tables" ("id");

ALTER TABLE "order_item" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_item" ADD FOREIGN KEY ("menu_id") REFERENCES "menus" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");
