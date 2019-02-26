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