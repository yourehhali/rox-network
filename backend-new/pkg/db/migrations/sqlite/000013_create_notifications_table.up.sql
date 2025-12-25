CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,   
    user_id INTEGER NOT NULL,               
    sender_id INTEGER,                      
    type TEXT NOT NULL,                    
    message TEXT NOT NULL,                 
    state TEXT DEFAULT 'unread',           
    group_id INTEGER,                       
    event_id INTEGER,                       
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP 
);
