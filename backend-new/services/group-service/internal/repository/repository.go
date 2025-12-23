package repository

import (
	"database/sql"
	"social-network/pkg/database"
	"time"
)

var (
	CHECK_GROUP_MEMBERSHIP     = "SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ? AND status = 'member'"
	INSERT_GROUP               = "INSERT INTO groups (title, description, creator_id) VALUES (?, ?, ?)"
	INSERT_GROUP_MEMBER        = "INSERT INTO group_members (group_id, user_id, status) VALUES (?, ?, ?)"
	UPDATE_GROUP_MEMBER_STATUS = "UPDATE group_members SET status = ? WHERE group_id = ? AND user_id = ?"
	REMOVE_GROUP_MEMBER        = "DELETE FROM group_members WHERE group_id = ? AND user_id = ?"
	SELECT_GROUPS_WITH_STATUS  = `
		SELECT g.id, g.title, g.description, g.creator_id, gm.status
		FROM groups g
		LEFT JOIN group_members gm
		  ON gm.group_id = g.id AND gm.user_id = ?
	`
	SELECT_GROUP_BY_ID           = "SELECT id, title, description, creator_id FROM groups WHERE id = ?"
	SELECT_GROUP_MEMBERS_IDS     = "SELECT user_id FROM group_members WHERE group_id = ? AND status = 'member'"
	SELECT_GROUP_MEMBERS_DETAILS = `
		SELECT u.id, u.username, u.first_name, u.last_name, COALESCE(u.avatar, ''), gm.status
		FROM group_members gm
		JOIN users u ON gm.user_id = u.id
		WHERE gm.group_id = ? AND gm.status = 'member'
	`
	INSERT_GROUP_POST = `
		INSERT INTO group_posts (title, content, image, creator_id, group_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	SELECT_GROUP_POST_BY_ID = `
		SELECT gp.id, gp.group_id, gp.creator_id, gp.title, gp.content, gp.image, gp.created_at, gp.created_at,
		       u.username, u.first_name || ' ' || u.last_name
		FROM group_posts gp
		JOIN users u ON gp.creator_id = u.id
		WHERE gp.id = ?
	`
	SELECT_GROUP_POSTS = `
		SELECT gp.id, gp.group_id, gp.creator_id, gp.title, gp.content, gp.image, gp.created_at, gp.created_at,
		       u.username, u.first_name || ' ' || u.last_name
		FROM group_posts gp
		JOIN users u ON gp.creator_id = u.id
		WHERE gp.group_id = ?
		ORDER BY gp.created_at DESC
	`
	INSERT_EVENT = `
		INSERT INTO events (title, description, creator_id, group_id, event_time)
		VALUES (?, ?, ?, ?, ?)`
	INSERT_EVENT_MEMBER = "INSERT INTO event_members (event_id, user_id, status) VALUES (?, ?, ?)"
	CHECK_EVENT_MEMBER  = "SELECT COUNT(*) FROM event_members WHERE event_id = ? AND user_id = ?"
	UPDATE_EVENT_MEMBER = "UPDATE event_members SET status = ? WHERE event_id = ? AND user_id = ?"
	SELECT_GROUP_EVENTS = `
		SELECT e.id, e.title, e.description, e.event_time, e.creator_id,
		       COALESCE(em.status, 'not_going') as user_status
		FROM events e
		LEFT JOIN event_members em ON e.id = em.event_id AND em.user_id = ?
		WHERE e.group_id = ?
	`
	SELECT_GROUPS_FOR_INVITE = `
    SELECT g.id, g.title, g.description, g.creator_id
    FROM groups g
    JOIN group_members current_member 
        ON current_member.group_id = g.id 
       AND current_member.user_id = ? 
       AND current_member.status = 'member'
    LEFT JOIN group_members target_member 
        ON target_member.group_id = g.id 
       AND target_member.user_id = ? 
       AND target_member.status = 'member'
    WHERE target_member.user_id IS NULL;
    `
)

type Group struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatorID   int    `json:"creator_id"`
}

type GroupWithStatus struct {
	Group
	UserStatus string `json:"user_status"`
}

type GroupPost struct {
	ID        int    `json:"id"`
	GroupID   int    `json:"group_id"`
	CreatorID int    `json:"creator_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	UserName  string `json:"username"`
	FullName  string `json:"full_name"`
}

type Event struct {
	ID          int    `json:"id"`
	GroupID     int    `json:"group_id"`
	CreatorID   int    `json:"creator_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EventTime   string `json:"event_time"`
	UserStatus  string `json:"user_status,omitempty"`
}

func CheckGroupMembership(groupID, userID int) (bool, error) {
	var count int
	err := database.DB.QueryRow(CHECK_GROUP_MEMBERSHIP, groupID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func InsertGroup(title, description string, creatorID int) (int64, error) {
	result, err := database.DB.Exec(INSERT_GROUP,
		title, description, creatorID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func InsertGroupMember(groupID, userID int, status string) error {
	_, err := database.DB.Exec(INSERT_GROUP_MEMBER,
		groupID, userID, status)
	return err
}

func UpdateGroupMemberStatus(groupID, userID int, status string) error {
	_, err := database.DB.Exec(UPDATE_GROUP_MEMBER_STATUS,
		status, groupID, userID)
	return err
}

func RemoveGroupMember(groupID, userID int) error {
	_, err := database.DB.Exec(REMOVE_GROUP_MEMBER,
		groupID, userID)
	return err
}

func GetGroups(userID int) ([]GroupWithStatus, error) {
	rows, err := database.DB.Query(SELECT_GROUPS_WITH_STATUS, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []GroupWithStatus
	for rows.Next() {
		var g GroupWithStatus
		var status sql.NullString
		if err := rows.Scan(&g.ID, &g.Title, &g.Description, &g.CreatorID, &status); err != nil {
			return nil, err
		}
		if status.Valid {
			g.UserStatus = status.String
		} else {
			g.UserStatus = "not_member"
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func GetGroupByID(groupID int) (*Group, error) {
	var g Group
	err := database.DB.QueryRow(SELECT_GROUP_BY_ID, groupID).
		Scan(&g.ID, &g.Title, &g.Description, &g.CreatorID)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func GetGroupMembers(groupID int) ([]int, error) {
	rows, err := database.DB.Query(SELECT_GROUP_MEMBERS_IDS, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		members = append(members, id)
	}
	return members, nil
}

type GroupMember struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	Status    string `json:"status"`
}

func GetGroupMembersDetails(groupID int) ([]GroupMember, error) {
	rows, err := database.DB.Query(SELECT_GROUP_MEMBERS_DETAILS, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []GroupMember
	for rows.Next() {
		var m GroupMember
		if err := rows.Scan(&m.ID, &m.Username, &m.FirstName, &m.LastName, &m.Avatar, &m.Status); err != nil {
			continue
		}
		members = append(members, m)
	}
	return members, nil
}

func InsertGroupPost(groupID, creatorID int, title, content, image string) (int64, error) {
	result, err := database.DB.Exec(INSERT_GROUP_POST,
		title, content, image, creatorID, groupID, time.Now())
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetGroupPostByID(postID int) (*GroupPost, error) {
	var p GroupPost
	var createdAt, updatedAt time.Time
	err := database.DB.QueryRow(SELECT_GROUP_POST_BY_ID, postID).Scan(&p.ID, &p.GroupID, &p.CreatorID, &p.Title, &p.Content, &p.Image, &createdAt, &updatedAt, &p.UserName, &p.FullName)

	if err != nil {
		return nil, err
	}
	p.CreatedAt = createdAt.Format(time.RFC3339)
	p.UpdatedAt = updatedAt.Format(time.RFC3339)
	return &p, nil
}

func GetGroupPosts(groupID int) ([]GroupPost, error) {
	rows, err := database.DB.Query(SELECT_GROUP_POSTS, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []GroupPost
	for rows.Next() {
		var p GroupPost
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&p.ID, &p.GroupID, &p.CreatorID, &p.Title, &p.Content, &p.Image, &createdAt, &updatedAt, &p.UserName, &p.FullName); err != nil {
			continue
		}
		p.CreatedAt = createdAt.Format(time.RFC3339)
		p.UpdatedAt = updatedAt.Format(time.RFC3339)
		posts = append(posts, p)
	}
	return posts, nil
}

func InsertEvent(groupID, creatorID int, title, description, eventTime string) (int64, error) {
	result, err := database.DB.Exec(INSERT_EVENT,
		title, description, creatorID, groupID, eventTime)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func InsertEventMember(eventID, userID int, status string) error {
	_, err := database.DB.Exec(INSERT_EVENT_MEMBER,
		eventID, userID, status)
	return err
}

func UpdateEventMemberStatus(eventID, userID int, status string) error {
	var count int
	err := database.DB.QueryRow(CHECK_EVENT_MEMBER, eventID, userID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		_, err = database.DB.Exec(UPDATE_EVENT_MEMBER, status, eventID, userID)
	} else {
		_, err = database.DB.Exec(INSERT_EVENT_MEMBER, eventID, userID, status)
	}
	return err
}

func GetGroupEvents(groupID, userID int) ([]Event, error) {
	rows, err := database.DB.Query(SELECT_GROUP_EVENTS, userID, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		e.GroupID = groupID
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.EventTime, &e.CreatorID, &e.UserStatus); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func GetGroupsForInvite(currentUserID, targetUserID int) ([]Group, error) {
	rows, err := database.DB.Query(SELECT_GROUPS_FOR_INVITE, currentUserID, targetUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Title, &g.Description, &g.CreatorID); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}
