PRAGMA foreign_keys=off;

CREATE TABLE event_members_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status TEXT CHECK(status IN ('going','notGoing')),
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO event_members_new (id, event_id, user_id, status)
SELECT id, event_id, user_id, 
       CASE status 
           WHEN 'not_going' THEN 'notGoing' 
           ELSE status 
       END
FROM event_members;

DROP TABLE event_members;
ALTER TABLE event_members_new RENAME TO event_members;

PRAGMA foreign_keys=on;
