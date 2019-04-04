package comment

import (
	pb "github.com/andreymgn/RSOI-comment/pkg/comment/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	statusInvalidUUID  = status.Error(codes.InvalidArgument, "invalid UUID")
	statusNotFound     = status.Error(codes.NotFound, "comment not found")
	statusInvalidToken = status.Errorf(codes.Unauthenticated, "invalid token")
)

func internalError(err error) error {
	return status.Error(codes.Internal, err.Error())
}

// SingleComment converts Comment to SingleComment
func (c *Comment) SingleComment() (*pb.SingleComment, error) {
	createdAtProto, err := ptypes.TimestampProto(c.CreatedAt)
	if err != nil {
		return nil, internalError(err)
	}

	modifiedAtProto, err := ptypes.TimestampProto(c.ModifiedAt)
	if err != nil {
		return nil, internalError(err)
	}

	res := new(pb.SingleComment)
	res.Uid = c.UID.String()
	res.UserUid = c.UserUID.String()
	res.PostUid = c.PostUID.String()
	res.Body = c.Body
	res.ParentUid = c.ParentUID.String()
	res.CreatedAt = createdAtProto
	res.ModifiedAt = modifiedAtProto
	res.IsDeleted = c.IsDeleted

	return res, nil
}

// ListComments returns all comments of post
func (s *Server) ListComments(ctx context.Context, req *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	var pageSize int32
	if req.PageSize == 0 {
		pageSize = 10
	} else {
		pageSize = req.PageSize
	}

	postUID, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	var parentUID uuid.UUID
	if req.CommentUid == "" {
		parentUID = uuid.Nil
	} else {
		parentUID, err = uuid.Parse(req.CommentUid)
		if err != nil {
			return nil, statusInvalidUUID
		}
	}

	comments, err := s.db.getAll(postUID, parentUID, pageSize, req.PageNumber)
	if err != nil {
		return nil, internalError(err)
	}

	res := new(pb.ListCommentsResponse)
	for _, comment := range comments {
		singleComment, err := comment.SingleComment()
		if err != nil {
			return nil, err
		}
		res.Comments = append(res.Comments, singleComment)
	}

	res.PageSize = pageSize
	res.PageNumber = req.PageNumber

	return res, nil
}

// CreateComment creates a new comment
func (s *Server) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.SingleComment, error) {
	postUID, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	parentUID := uuid.Nil
	if req.ParentUid != "" {
		parentUID, err = uuid.Parse(req.ParentUid)
		if err != nil {
			return nil, statusInvalidUUID
		}
	}

	userUID, err := uuid.Parse(req.UserUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	comment, err := s.db.create(postUID, req.Body, parentUID, userUID)
	if err != nil {
		return nil, internalError(err)
	}

	return comment.SingleComment()
}

// UpdateComment updates comment by ID
func (s *Server) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	err = s.db.update(uid, req.Body)
	switch err {
	case nil:
		return new(pb.UpdateCommentResponse), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

func (s *Server) RemoveContent(ctx context.Context, req *pb.RemoveContentRequest) (*pb.RemoveContentResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	err = s.db.removeContent(uid)
	switch err {
	case nil:
		return new(pb.RemoveContentResponse), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

// DeleteComment deletes post by ID
func (s *Server) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, err
	}

	err = s.db.delete(uid)
	switch err {
	case nil:
		return new(pb.DeleteCommentResponse), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

// GetOwner returns comment owner
func (s *Server) GetOwner(ctx context.Context, req *pb.GetOwnerRequest) (*pb.GetOwnerResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	result, err := s.db.getOwner(uid)
	switch err {
	case nil:
		res := new(pb.GetOwnerResponse)
		res.OwnerUid = result
		return res, nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}
