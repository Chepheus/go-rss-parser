CREATE TABLE IF NOT EXISTS posts (
     id    integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
     title   varchar(255) NOT NULL CHECK (title <> ''),
     short_description varchar(255) NOT NULL CHECK (short_description <> ''),
     external_post_link varchar(255) NOT NULL CHECK (external_post_link <> ''),
     thumbnail varchar(255) NOT NULL CHECK (thumbnail <> ''),
     pub_date date
);