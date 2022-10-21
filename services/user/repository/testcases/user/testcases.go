package testcases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/testutils"
	"github.com/panagiotisptr/hermes-messenger/user/model"
	"github.com/panagiotisptr/hermes-messenger/user/repository"
)

type CreateTestcaseInput struct {
	ctx  context.Context
	args *model.User
	repo repository.UserRepository
}

type CreateTestcaseOutput struct {
	u   *model.User
	err error
}

func UserDetailsMatch(
	a *model.User,
	b *model.User,
) error {
	if err := testutils.Assert(
		"Email",
		a.Email,
		b.Email,
		testutils.StringsEqual,
	); err != nil {
		return err
	}
	if err := testutils.Assert(
		"FirstName",
		a.FirstName,
		b.FirstName,
		testutils.StringsEqual,
	); err != nil {
		return err
	}
	if err := testutils.Assert(
		"LastName",
		a.LastName,
		b.LastName,
		testutils.StringsEqual,
	); err != nil {
		return err
	}

	return nil
}

func CreateTestcases(
	r repository.UserRepository,
) []testutils.Testcase[CreateTestcaseInput, CreateTestcaseOutput] {
	ctx := context.WithValue(
		context.Background(),
		"user-id",
		uuid.New(),
	)
	ctx = context.WithValue(
		ctx,
		"request-id",
		uuid.New(),
	)
	testcases := []testutils.Testcase[
		CreateTestcaseInput,
		CreateTestcaseOutput,
	]{
		{
			Name: "base case",
			Input: CreateTestcaseInput{
				ctx: ctx,
				args: &model.User{
					Email:     "email@domain.localhost",
					FirstName: "firstName",
					LastName:  "lastName",
				},
				repo: r,
			},
			Process: func(i CreateTestcaseInput) CreateTestcaseOutput {
				u, err := r.Create(i.ctx, i.args)

				return CreateTestcaseOutput{
					u:   u,
					err: err,
				}
			},
			Check: func(i CreateTestcaseInput, o CreateTestcaseOutput) error {
				if o.err != nil {
					return fmt.Errorf("Expected no errors. Got error %v", o.err)
				}
				if err := testutils.AssertNotNil(
					"user",
					o.u,
				); err != nil {
					return err
				}
				if err := UserDetailsMatch(
					i.args,
					o.u,
				); err != nil {
					return err
				}
				userId, ok := ctx.Value("user-id").(uuid.UUID)
				if err := testutils.AssertEqual(
					"userId exists",
					ok,
					true,
				); err != nil {
					return err
				}
				if err := testutils.AssertEqual(
					"userId",
					userId,
					*o.u.CreatedBy,
				); err != nil {
					return err
				}
				if err := testutils.AssertEqual(
					"userId",
					userId,
					*o.u.CreatedBy,
				); err != nil {
					return err
				}

				return nil
			},
		},
	}

	return testcases
}
