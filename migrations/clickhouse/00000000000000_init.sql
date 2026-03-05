CREATE TABLE user_actions (
    dt DateTime,
    user_id Int8,
    action_id Int8
) ENGINE MergeTree()
ORDER BY (dt)
TTL dt + INTERVAL 5 MONTH;