# leaderboard-app - a serverless application

[![OpenFaaS](https://img.shields.io/badge/openfaas-cloud-blue.svg)](https://www.openfaas.com)

This application is an example of how to write a Single Page App (SPA) with a Serverless approach. It provides a live leaderboard for your GitHub organisation or repos showing comments made and issues opened by your community and contributors.

* The front-end is written with Vue.js
* The backing data-store data is Postgres with a remote DBaaS or in-cluster deployment

See a live example tracking the openfaas/openfaas-incubator organizations: [https://alexellis.o6s.io/leaderboard-page](https://alexellis.o6s.io/leaderboard-page)

To test out the functionality comment on this issue: [Issue: Let's test the leaderboard!](https://github.com/openfaas/org-tester/issues/18)

Here's a preview of the app when the dark theme is enabled: 

![Dark Leaderboard example](docs/leaderboard-dark.png)

Subscribe to events by adding a webhook to the github-sub function:

![Subscribe](docs/subscribe.png)

## Functions

* github-sub

Receives webhooks from GitHub via an organization or repo subscription. Secured with HMAC by Alex Ellis

* leaderboard

Retrieves the current leaderboard in JSON by Alex Ellis

* leaderboard-page

Renders the leaderboard itself as a Vue.js app by Ken Fukuyama

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
insert into users (user_id,user_login,track, created_at) values (653013,'alexellisuk',true,now());
insert into users (user_id,user_login,track, created_at) values (103022,'rgee0',true,now());

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
    RETURNS TABLE(user_id integer, user_login text, issue_comments bigint, pr_comments bigint, issues_created bigint, pr_created bigint)
  AS
$$
BEGIN
RETURN QUERY select
    a.user_id,
    u.user_login,
    count(ic.activity_type) as issue_comments,
    count(prc.activity_type) as pr_comments,
    count(cm.activity_type) as issues_created,
    count(pr.activity_type) as pr_created
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
LEFT OUTER JOIN activity AS prc
ON a.user_id = prc.user_id
    and a.activity_type = prc.activity_type
    and a.activity_date = prc.activity_date
    and a.owner = prc.owner
    and a.repo = prc.repo
    and prc.activity_type = 'pr_review_comment'
left outer join activity as cm
on a.user_id= cm.user_id
    and a.activity_type = cm.activity_type
    and a.activity_date = cm.activity_date
    and a.owner = cm.owner
    and a.repo = cm.repo
    and cm.activity_type = 'issue_created'
LEFT OUTER JOIN activity AS pr
ON a.user_id = pr.user_id
    and a.activity_type = pr.activity_type
    and a.activity_date = pr.activity_date
    and a.owner = pr.owner
    and a.repo = pr.repo
    and pr.activity_type = 'pull_request_opened'
where u.track = true
group by a.user_id, u.user_login
order by issue_comments desc, issues_created desc;
END
$$  LANGUAGE 'plpgsql' VOLATILE;

select * from get_leaderboard();
```

## Contributing & license

Please feel free to fork and star this repo and use it as a template for your own applications. The license is MIT.

To contribute see [CONTRIBUTING.md](./CONTRIBUTING.md)


