DELETE FROM group_members
WHERE id NOT IN (
    SELECT MIN(id)
    FROM group_members
    GROUP BY group_id, user_id
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_group_members_unique ON group_members(group_id, user_id);
