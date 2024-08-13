package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/arya2004/xyfin/db/sqlc"
	"github.com/arya2004/xyfin/pb"
	"github.com/arya2004/xyfin/utils"
	"github.com/arya2004/xyfin/validators"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, validators.UnauthenticatedError(err)
	}

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, validators.InvalidArgumentError(violations)
	}

	if authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other users")
	}

	

	arg := db.UpdateUserParams{
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid: req.FullName != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid: req.Email != nil,
		},

		Username: req.GetUsername(),
		
	}

	if req.Password != nil {
		hashedPassword, err := utils.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failled to hash -password: %v", err)
		}
		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid: true,
		}

		arg.PasswordChangedAt = sql.NullTime{
			Time: time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {

		if err == sql.ErrNoRows{
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failled to update user: %v", err)

	}

	resp := &pb.UpdateUserResponse{
		User: ConvertUser(user),
	}
	
	return resp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validators.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, validators.FieldViolation("username", err))
	}

	if req.Password != nil {
		if err := validators.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, validators.FieldViolation("password", err))
		}
	}

	if req.FullName != nil{
		if err := validators.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, validators.FieldViolation("full_name", err))
		}
	}

	if req.Email != nil{
		if err := validators.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, validators.FieldViolation("email", err))
		}
	}

	

	return violations
}