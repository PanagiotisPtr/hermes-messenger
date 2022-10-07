package memory

import (
	"context"
	"fmt"
	"testing"

	"github.com/panagiotisptr/hermes-messenger/libs/utils/testutils"
	"github.com/panagiotisptr/hermes-messenger/user/model"
	"github.com/panagiotisptr/hermes-messenger/user/repository"
	"go.uber.org/zap"
)

func getRepository() repository.Repository {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return ProvideUserRepository(logger)
}

func TestCreate(t *testing.T) {
	type input struct {
		ctx  context.Context
		args model.UserDetails
		repo repository.Repository
	}
	type output struct {
		u   *model.User
		err error
	}
	repo := getRepository()
	userDetailsMatch := func(d model.UserDetails, u *model.User) error {
		if err := testutils.Assert(
			"Email",
			d.Email,
			u.Email,
			testutils.StringsEqual,
		); err != nil {
			return err
		}
		if err := testutils.Assert(
			"FirstName",
			d.FirstName,
			u.FirstName,
			testutils.StringsEqual,
		); err != nil {
			return err
		}
		if err := testutils.Assert(
			"LastName",
			d.LastName,
			u.LastName,
			testutils.StringsEqual,
		); err != nil {
			return err
		}

		return nil
	}
	testcases := []testutils.Testcase[input, output]{
		{
			Name: "base case",
			Input: input{
				ctx: context.Background(),
				args: model.UserDetails{
					Email:     "email@domain.localhost",
					FirstName: "firstName",
					LastName:  "lastName",
				},
				repo: repo,
			},
			Process: func(i input) output {
				u, err := repo.Create(i.ctx, i.args)

				return output{
					u:   u,
					err: err,
				}
			},
			Check: func(i input, o output) error {
				if o.err != nil {
					return fmt.Errorf("Expected no errors. Got error %v", o.err)
				}
				if err := testutils.AssertNotNil(
					"user",
					o.u,
				); err != nil {
					return err
				}

				return userDetailsMatch(i.args, o.u)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, tc.RunFunc())
	}
}
