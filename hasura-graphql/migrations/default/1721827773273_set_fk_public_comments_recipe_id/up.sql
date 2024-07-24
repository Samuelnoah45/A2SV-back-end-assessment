alter table "public"."comments"
  add constraint "comments_recipe_id_fkey"
  foreign key ("recipe_id")
  references "public"."recipes"
  ("id") on update cascade on delete cascade;
