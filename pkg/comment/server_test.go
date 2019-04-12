package comment

import (
	"errors"
	"testing"
	"time"

	pb "github.com/andreymgn/RSOI-comment/pkg/comment/proto"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

var (
	errDummy     = errors.New("dummy")
	dummyUID     = uuid.New()
	nilUIDString = uuid.Nil.String()
)

type mockdb struct{}

func (mdb *mockdb) getAll(postUID uuid.UUID, parentUID uuid.UUID, pageNumber, pageSize int32) ([]*Comment, error) {
	result := make([]*Comment, 0)
	uid1 := uuid.New()
	uid2 := uuid.New()
	uid3 := uuid.New()
	pUID := uuid.New()

	result = append(result, &Comment{uid1, uid2, pUID, "first comment body", uuid.Nil, time.Now(), time.Now(), false})
	result = append(result, &Comment{uid2, uid3, pUID, "second comment body", uuid.Nil, time.Now(), time.Now(), false})
	result = append(result, &Comment{uid3, uid1, pUID, "third comment body", uid1, time.Now(), time.Now(), false})
	return result, nil
}

func (mdb *mockdb) getOne(uid uuid.UUID) (*Comment, error) {
	if uid == uuid.Nil {
		uid := uuid.New()

		return &Comment{uid, uid, uid, "first comment body", uuid.Nil, time.Now(), time.Now(), false}, nil
	}

	return nil, errDummy
}

func (mdb *mockdb) create(postUID uuid.UUID, body string, parentUID, userUID uuid.UUID) (*Comment, error) {
	if postUID == uuid.Nil {
		uid := uuid.New()
		return &Comment{uid, userUID, postUID, "first comment body", uuid.Nil, time.Now(), time.Now(), false}, nil
	}

	return nil, errDummy
}

func (mdb *mockdb) update(uid uuid.UUID, body string) error {
	if uid == uuid.Nil {
		return nil
	}

	return errDummy
}

func (mdb *mockdb) removeContent(uid uuid.UUID) error {
	if uid == uuid.Nil {
		return nil
	}

	return errDummy
}

func (mdb *mockdb) delete(uid uuid.UUID) error {
	if uid == uuid.Nil {
		return nil
	}

	return errDummy
}

func (mdb *mockdb) getOwner(uid uuid.UUID) (string, error) {
	return nilUIDString, nil
}

func TestListComments(t *testing.T) {
	s := &Server{&mockdb{}}
	var pageSize int32 = 3
	req := &pb.ListCommentsRequest{PostUid: nilUIDString, PageSize: pageSize}
	res, err := s.ListComments(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if len(res.Comments) != int(pageSize) {
		t.Errorf("unexpected number of comments: got %v want %v", len(res.Comments), pageSize)
	}
}

func TestGetComment(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.GetCommentRequest{Uid: nilUIDString}
	_, err := s.GetComment(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestGetCommentFail(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.GetCommentRequest{Uid: ""}
	_, err := s.GetComment(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestCreateComment(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.CreateCommentRequest{PostUid: nilUIDString, UserUid: nilUIDString}
	_, err := s.CreateComment(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestCreateCommentFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.CreateCommentRequest{}
	_, err := s.CreateComment(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestUpdateComment(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.UpdateCommentRequest{Uid: nilUIDString}
	_, err := s.UpdateComment(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestUpdateCommentFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.UpdateCommentRequest{}
	_, err := s.UpdateComment(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestDeleteComment(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.DeleteCommentRequest{Uid: nilUIDString}
	_, err := s.DeleteComment(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestDeleteCOmmentFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.DeleteCommentRequest{}
	_, err := s.DeleteComment(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}
