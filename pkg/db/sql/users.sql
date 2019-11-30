CREATE TABLE users (
   id serial PRIMARY KEY,
   uuid uuid UNIQUE NOT NULL,

   first_name varchar(60),
   last_name varchar(60),
   username varchar(60) UNIQUE,
   password varchar(60),
   email_address varchar(100),

   created_at date default CURRENT_DATE,
   last_updated date
);
