package comment

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var (
	errNotCreated = errors.New("comment not created")
	errNotFound   = errors.New("comment not found")
)

// Comment describes comment to a post
type Comment struct {
	UID        uuid.UUID
	UserUID    uuid.UUID
	PostUID    uuid.UUID
	Body       string
	ParentUID  uuid.UUID
	CreatedAt  time.Time
	ModifiedAt time.Time
	IsDeleted  bool
}

type datastore interface {
	getAll(uuid.UUID, uuid.UUID, int32, int32) ([]*Comment, error)
	create(uuid.UUID, string, uuid.UUID, uuid.UUID) (*Comment, error)
	update(uuid.UUID, string) error
	removeContent(uuid.UUID) error
	delete(uuid.UUID) error
	getOwner(uuid.UUID) (string, error)
}

type db struct {
	*sql.DB
}

func newDB(connString string) (*db, error) {
	postgres, err := sql.Open("postgres", connString)
	return &db{postgres}, err
}

func (db *db) getAll(postUID uuid.UUID, parentUID uuid.UUID, pageSize, pageNumber int32) ([]*Comment, error) {
	query := "SELECT * FROM comments WHERE post_uid=$1 AND parent_uid=$2 ORDER BY created_at DESC LIMIT $3 OFFSET $4"

	lastRecord := pageNumber * pageSize
	rows, err := db.Query(query, postUID.String(), parentUID.String(), pageSize, lastRecord)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	result := make([]*Comment, 0)
	for rows.Next() {
		comment := new(Comment)
		var uid, userUID, pUID, parentUID string
		err := rows.Scan(&uid, &userUID, &pUID, &comment.Body, &parentUID, &comment.CreatedAt, &comment.ModifiedAt, &comment.IsDeleted)
		if err != nil {
			return nil, err
		}

		comment.UID, err = uuid.Parse(uid)
		if err != nil {
			return nil, err
		}

		comment.UserUID, err = uuid.Parse(userUID)
		if err != nil {
			return nil, err
		}

		comment.PostUID, err = uuid.Parse(pUID)
		if err != nil {
			return nil, err
		}

		if parentUID != "" {
			comment.ParentUID, err = uuid.Parse(parentUID)
			if err != nil {
				return nil, err
			}
		} else {
			comment.ParentUID = uuid.Nil
		}

		result = append(result, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *db) create(postUID uuid.UUID, body string, parentUID, userUID uuid.UUID) (*Comment, error) {
	comment := new(Comment)

	query := "INSERT INTO comments (uid, user_uid, post_uid, body, parent_uid, created_at, modified_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	uid := uuid.New()
	now := time.Now()

	comment.UID = uid
	comment.UserUID = userUID
	comment.PostUID = postUID
	comment.Body = body
	comment.ParentUID = parentUID
	comment.CreatedAt = now
	comment.ModifiedAt = now

	result, err := db.Exec(query, uid.String(), userUID.String(), postUID.String(), body, parentUID.String(), now, now)
	nRows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if nRows == 0 {
		return nil, errNotCreated
	}

	return comment, nil
}

func (db *db) update(uid uuid.UUID, body string) error {
	query := "UPDATE comments SET body=$1, modified_at=$2 WHERE uid=$3 AND is_deleted=false"
	result, err := db.Exec(query, body, time.Now(), uid.String())
	nRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if nRows == 0 {
		return errNotFound
	}

	return nil
}

func (db *db) removeContent(uid uuid.UUID) error {
	query := "UPDATE comments SET is_deleted=true, modified_at=$1 WHERE uid=$2 AND is_deleted=false"
	result, err := db.Exec(query, time.Now(), uid.String())
	nRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if nRows == 0 {
		return errNotFound
	}

	return nil
}

func (db *db) delete(uid uuid.UUID) error {
	query := "DELETE FROM comments WHERE uid=$1"
	result, err := db.Exec(query, uid.String())
	nRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if nRows == 0 {
		return errNotFound
	}

	return nil
}

func (db *db) getOwner(uid uuid.UUID) (string, error) {
	query := "SELECT user_uid FROM comments WHERE uid=$1"
	row := db.QueryRow(query, uid.String())
	var result string
	switch err := row.Scan(&result); err {
	case nil:
		return result, nil
	case sql.ErrNoRows:
		return "", errNotFound
	default:
		return "", err
	}
}
