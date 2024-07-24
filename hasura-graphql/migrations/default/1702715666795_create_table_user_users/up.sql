CREATE TABLE "user"."users" ("id" uuid NOT NULL, "first_name" text NOT NULL, "last_name" text NOT NULL, "email" text NOT NULL, "phone_number" text NOT NULL, PRIMARY KEY ("id") , UNIQUE ("email"), UNIQUE ("phone_number"));