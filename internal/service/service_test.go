package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mirhijinam/pg-project/internal/model"
)

func TestCommandService_CreateCommand(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	type fields struct {
		repoBehave func(*MockRepo)
		Workers    *sync.Map
		poolBehave func(*MockPool)
	}
	type args struct {
		cmd    *model.Command
		isLong bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success creating command",
			fields: fields{
				repoBehave: func(mockRepo *MockRepo) {
					mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
					mockRepo.EXPECT().SetSuccess(gomock.Any(), gomock.Any()).Return(nil).Times(1)
					mockRepo.EXPECT().Writer(gomock.Any(), gomock.Any()).Return(func(p []byte) (n int, err error) {
						return len(p), nil
					}).Times(1)
				},
				Workers: new(sync.Map),
				poolBehave: func(mockPool *MockPool) {
					mockPool.EXPECT().Go(gomock.Any()).Do(func(f func()) {
						f()
					}).Times(1)
				},
			},
			args: args{
				cmd:    &model.Command{Name: "echo", Raw: "echo s1"},
				isLong: false,
			},
			wantErr: false,
		},
		{
			name: "Failure executing command with error",
			fields: fields{
				repoBehave: func(mockRepo *MockRepo) {
					mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
					expectedError := errors.New("exit status 127")
					mockRepo.EXPECT().SetError(gomock.Any(), gomock.Any(), gomock.Eq(expectedError.Error())).Return(nil).Times(1)
				},
				Workers: new(sync.Map),
				poolBehave: func(mockPool *MockPool) {
					mockPool.EXPECT().Go(gomock.Any()).Do(func(f func()) {
						f()
					}).Times(1)
				},
			},
			args: args{
				cmd:    &model.Command{Name: "test", Raw: "eho 'Hello, World!'"},
				isLong: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := NewMockRepo(ctrl)
			poolMock := NewMockPool(ctrl)
			wg := &sync.WaitGroup{}

			tt.fields.repoBehave(repoMock)
			poolMock.EXPECT().Go(gomock.Any()).Do(func(f func()) {
				wg.Add(1)
				go func() {
					defer wg.Done()
					f()
				}()
			}).AnyTimes()

			cs := &CommandService{
				repo:    repoMock,
				Workers: tt.fields.Workers,
				Pool:    poolMock,
			}

			if err := cs.CreateCommand(ctx, tt.args.cmd, tt.args.isLong); (err != nil) != tt.wantErr {
				t.Errorf("CommandService.CreateCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
			wg.Wait() // Дождаться завершения всех горутин перед завершением теста
		})
	}
}
