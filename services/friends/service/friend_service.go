package service

import (
	"context"

	"github.com/panagiotisptr/hermes-messenger/friends/model"
	"github.com/panagiotisptr/hermes-messenger/friends/repository"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils/filter"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/grpcserviceutils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/loggingutils"
	"go.uber.org/zap"
)

// FriendService represents a friend service
type FriendService struct {
	logger     *zap.Logger
	friendRepo repository.FriendRepository
}

// ProvideFriendService provides an instance of the friend
// service
func ProvideFriendService(
	logger *zap.Logger,
	friendRepo repository.FriendRepository,
) *FriendService {
	return &FriendService{
		logger: logger.With(
			zap.String("service", "FriendService"),
		),
		friendRepo: friendRepo,
	}
}

func (s *FriendService) Create(
	ctx context.Context,
	args *model.Friend,
) (*model.Friend, error) {
	l := loggingutils.GetLoggerWith(
		ctx,
		s.logger,
		loggingutils.LoggerWithRequestID,
		loggingutils.LoggerWithUserID,
	).With(
		zap.String("method", "AddFriend"),
	).Sugar()

	userId, err := grpcserviceutils.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	// User side
	f := filter.NewFilter()
	f.Add("userId", filter.Eq, userId)
	f.Add("friendId", filter.Eq, args.FriendId)
	ofr, err := s.friendRepo.FindOne(ctx, f)
	if err != nil {
		l.Error("failed to find friend", err)

		return nil, err
	}

	if ofr != nil {
		switch ofr.Status {
		case model.FriendStatusPendingFriend:
			return ofr, nil
		case model.FriendStatusDeclined:
			ofr.Status = model.FriendStatusPendingFriend
			_, err := s.friendRepo.Update(ctx, f, ofr)
			if err != nil {
				l.Error("failed to update friend", err)
			}
		case model.FriendStatusPendingUser:
			ofr.Status = model.FriendStatusAccepted
			_, err := s.friendRepo.Update(ctx, f, ofr)
			if err != nil {
				l.Error("failed to update friend", err)
			}
		}
	} else {
		_, err = s.friendRepo.Create(ctx, &model.Friend{
			UserId:   userId,
			FriendId: args.FriendId,
			Status:   model.FriendStatusPendingFriend,
		})
		if err != nil {
			l.Error("failed to create friend", err)
		}
	}

	// Friend side
	f = filter.NewFilter()
	f.Add("userId", filter.Eq, args.FriendId)
	f.Add("friendId", filter.Eq, userId)
	tfr, err := s.friendRepo.FindOne(ctx, f)
	if err != nil {
		l.Error("failed to find friend", err)

		return ofr, err
	}

	if tfr != nil {
		switch tfr.Status {
		case model.FriendStatusPendingFriend:
			ofr.Status = model.FriendStatusAccepted
			_, err := s.friendRepo.Update(ctx, f, ofr)
			if err != nil {
				l.Error("failed to update friend", err)

				return ofr, err
			}
		}
	} else {
		_, err = s.friendRepo.Create(ctx, &model.Friend{
			UserId:   userId,
			FriendId: args.FriendId,
			Status:   model.FriendStatusPendingUser,
		})
		if err != nil {
			l.Error("failed to create friend", err)
		}
	}

	return ofr, err
}

func (s *FriendService) Delete(
	ctx context.Context,
	args *model.Friend,
) error {
	l := loggingutils.GetLoggerWith(
		ctx,
		s.logger,
		loggingutils.LoggerWithRequestID,
		loggingutils.LoggerWithUserID,
	).With(
		zap.String("method", "DeleteFriend"),
	).Sugar()

	userId, err := grpcserviceutils.GetUserID(ctx)
	if err != nil {
		return err
	}

	// User side
	f := filter.NewFilter()
	f.Add("userId", filter.Eq, userId)
	f.Add("friendId", filter.Eq, args.FriendId)
	ofr, err := s.friendRepo.FindOne(ctx, f)
	if err != nil {
		l.Error("failed to find friend", err)

		return err
	}
	if ofr == nil {
		return nil
	}
	_, err = s.friendRepo.Delete(ctx, f)
	if err != nil {
		l.Error("failed to delete friend", err)

		return err
	}

	// Friend side
	f = filter.NewFilter()
	f.Add("userId", filter.Eq, args.FriendId)
	f.Add("friendId", filter.Eq, userId)
	tfr, err := s.friendRepo.FindOne(ctx, f)
	if err != nil {
		l.Error("failed to find friend", err)

		return err
	}
	if tfr == nil {
		return nil
	}
	_, err = s.friendRepo.Delete(ctx, f)
	if err != nil {
		l.Error("failed to delete friend", err)
	}

	return err
}

func (s *FriendService) Find(
	ctx context.Context,
) (<-chan *model.Friend, error) {
	ch := make(chan *model.Friend)
	l := loggingutils.GetLoggerWith(
		ctx,
		s.logger,
		loggingutils.LoggerWithRequestID,
		loggingutils.LoggerWithUserID,
	).With(
		zap.String("method", "GetFriends"),
	).Sugar()

	userId, err := grpcserviceutils.GetUserID(ctx)
	if err != nil {
		return ch, err
	}

	f := filter.NewFilter()
	f.Add("userId", filter.Eq, userId)
	frs, err := s.friendRepo.Find(ctx, f)
	if err != nil {
		l.Error("failed to find friends", err)
	}

	return frs, err
}
