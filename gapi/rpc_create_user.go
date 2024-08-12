package gapi

import (
	"context"

	db "github.com/arya2004/xyfin/db/sqlc"
	"github.com/arya2004/xyfin/pb"
	"github.com/arya2004/xyfin/utils"
	"github.com/arya2004/xyfin/validators"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, validators.InvalidArgumentError(violations)
	}

	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failled to hash -password: %v", err)
	}

	arg := db.CreateUserParams{
		FullName: req.GetFullName(),
		Email: req.GetEmail(),
		Username: req.GetUsername(),
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exist %v", err)

			}
		}
		return nil, status.Errorf(codes.Internal, "failled to create user: %v", err)

	}

	resp := &pb.CreateUserResponse{
		User: ConvertUser(user),
	}
	
	return resp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validators.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, validators.FieldViolation("username", err))
	}

	if err := validators.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, validators.FieldViolation("password", err))
	}

	if err := validators.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, validators.FieldViolation("full_name", err))
	}

	if err := validators.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, validators.FieldViolation("email", err))
	}

	return violations
}