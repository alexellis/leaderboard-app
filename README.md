# leaderboard-app

GitHub leaderboard for your organisation or repo

## Functions

* github-sub
* leaderboard

## Schema

```sql
drop table activity cascade;
drop table users;

CREATE TABLE users (
    user_id         integer PRIMARY KEY NOT NULL,
    user_login      text NOT NULL,
    track           BOOLEAN NOT NULL,
    created_at      timestamp not null
);
insert into users (user_id,user_login,track) values (653013,'alexellis',true);
insert into users (user_id,user_login,track) values (103022,'rgee0',true);

CREATE TABLE activity (
    id              INT GENERATED ALWAYS AS IDENTITY,
    user_id         integer NOT NULL references users(user_id),
    activity_type   text NOT NULL,
    activity_date   timestamp NOT NULL,
    owner           text NOT NULL,
    repo            text NOT NULL
);

insert into activity (id,user_id,activity_type,activity_date,owner,repo) values (DEFAULT,653013,'issue_created','2019-02-13 07:44:00','openfaas','org-tester');
insert into activity (id,user_id,activity_type,activity_date,owner,repo) values (DEFAULT,653013,'issue_comment','2019-02-13 07:44:05','openfaas','org-tester');
insert into activity (id,user_id,activity_type,activity_date,owner,repo) values (DEFAULT,653013,'issue_comment','2019-02-12 07:44:05','openfaas','org-tester');
insert into activity (id,user_id,activity_type,activity_date,owner,repo) values (DEFAULT,103022,'issue_comment','2019-02-12 07:44:05','openfaas','org-tester');

select * from activity order by activity_date  asc; 

select a.user_id, a.activity_type, count(a.id) from activity as a
where a.activity_date <= now()
group by (a.user_id, a.activity_type) ;

drop function get_leaderboard;

CREATE or REPLACE FUNCTION get_leaderboard()
    RETURNS TABLE(user_id integer, user_login text, issue_comments bigint, issues_created bigint)
  AS
$$
BEGIN
RETURN QUERY select
    a.user_id,
    u.user_login,
    count(ic.activity_type) as issue_comments,
    count(cm.activity_type) as issues_created
from activity as a
inner join users u
on a.user_id = u.user_id
left outer join activity as ic
on a.user_id= ic.user_id
    and a.activity_type = ic.activity_type
    and a.activity_date = ic.activity_date
    and a.owner = ic.owner
    and a.repo = ic.repo
    and ic.activity_type = 'issue_comment'
left outer join activity as cm
on a.user_id= cm.user_id
    and a.activity_type = cm.activity_type
    and a.activity_date = cm.activity_date
    and a.owner = cm.owner
    and a.repo = cm.repo
    and cm.activity_type = 'issue_created'
where u.track = true
group by a.user_id, u.user_login
order by issue_comments desc, issues_created desc;
END
$$  LANGUAGE 'plpgsql' VOLATILE;

select * from get_leaderboard();

```



