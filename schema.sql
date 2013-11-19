DROP TABLE movie;
CREATE TABLE movie
(
  id integer primary key,
  name varchar(255),
  year integer,
  resolution varchar(255),
  filetype varchar(255),
  location varchar(255),
  picture varchar(255)
);
CREATE INDEX movie_name_idx ON movie (name);
CREATE UNIQUE INDEX movie_location_idx ON movie (location);