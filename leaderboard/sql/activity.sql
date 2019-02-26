CREATE TABLE activity (
    id              INT GENERATED ALWAYS AS IDENTITY,
    user_id         integer NOT NULL references users(user_id),
    activity_type   text NOT NULL,
    activity_date   timestamp NOT NULL,
    owner           text NOT NULL,
    repo            text NOT NULL
);