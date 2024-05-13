package service

import (
	"context"
	"errors"
	"reflect"
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

func TestCommandService_GetCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepo(ctrl)

	type fields struct {
		repo    Repo
		Workers *sync.Map
		Pool    Pool
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Command
		wantErr bool
		setup   func()
	}{
		{
			name: "Success case",
			fields: fields{
				repo: mockRepo,
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			want: model.Command{
				Id:   1,
				Name: "echo",
				Raw:  "echo 'hello world'",
			},
			wantErr: false,
			setup: func() {
				mockRepo.EXPECT().GetOne(gomock.Any(), 1).Return(model.Command{
					Id:   1,
					Name: "echo",
					Raw:  "echo 'hello world'",
				}, nil).Times(1)
			},
		},
		{
			name: "Failure case",
			fields: fields{
				repo: mockRepo,
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			want:    model.Command{},
			wantErr: true,
			setup: func() {
				mockRepo.EXPECT().GetOne(gomock.Any(), 1).Return(model.Command{}, errors.New("not found")).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			cs := &CommandService{
				repo:    tt.fields.repo,
				Workers: tt.fields.Workers,
				Pool:    tt.fields.Pool,
			}
			got, err := cs.GetCommand(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandService.GetCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommandService.GetCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestCommandService_GetCommandList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepo(ctrl)

	type fields struct {
		repo    Repo
		Workers *sync.Map
		Pool    Pool
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Command
		wantErr bool
		setup   func()
	}{
		{
			name: "Success case",
			fields: fields{
				repo: mockRepo,
			},
			args: args{
				ctx: context.TODO(),
			},
			want: []model.Command{
				{Id: 1, Name: "Command One", Raw: "echo 'one'"},
				{Id: 2, Name: "Command Two", Raw: "echo 'two'"},
			},
			wantErr: false,
			setup: func() {
				mockRepo.EXPECT().GetList(gomock.Any()).Return([]model.Command{
					{Id: 1, Name: "Command One", Raw: "echo 'one'"},
					{Id: 2, Name: "Command Two", Raw: "echo 'two'"},
				}, nil).Times(1)
			},
		},
		{
			name: "Failure case",
			fields: fields{
				repo: mockRepo,
			},
			args: args{
				ctx: context.TODO(),
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockRepo.EXPECT().GetList(gomock.Any()).Return(nil, errors.New("database error")).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			cs := &CommandService{
				repo:    tt.fields.repo,
				Workers: tt.fields.Workers,
				Pool:    tt.fields.Pool,
			}
			got, err := cs.GetCommandList(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandService.GetCommandList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommandService.GetCommandList() = %v, want %v", got, tt.want)
			}
		})
	}
}
