package zookeeper

import (
	"errors"
	"testing"

	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/sync"
	"github.com/ContinuumLLC/platform-common-lib/src/sync/zookeeper/mock"
	"github.com/golang/mock/gomock"
	"github.com/samuel/go-zookeeper/zk"
)

func TestInstance(t *testing.T) {
	srv := Instance(sync.Config{})
	_, ok := srv.(zookeeper)
	if !ok {
		t.Error("Zookeeper Service Type is not IMPL")
	}
}

func Test_zookeeper_send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	con := mock.NewMockConnection(ctrl)
	con.EXPECT().Exists("1").Return(false, &zk.Stat{}, errors.New("Error"))

	con.EXPECT().Exists("2").Return(false, &zk.Stat{}, nil)
	con.EXPECT().Create("2", gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New("Error"))

	con.EXPECT().Exists("3").Return(true, &zk.Stat{}, nil)
	con.EXPECT().Set("3", gomock.Any(), gomock.Any()).Return(nil, errors.New("Error"))

	con.EXPECT().Exists("4").Return(true, &zk.Stat{}, nil)
	con.EXPECT().Set("4", gomock.Any(), gomock.Any()).Return(&zk.Stat{}, nil)

	type fields struct {
		config sync.Config
		log    logging.Logger
	}
	type args struct {
		path string
		data string
		conn Connection
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "1",
			fields:  fields{config: sync.Config{}, log: logging.GetLoggerFactory().Get()},
			args:    args{path: "1", data: "", conn: con},
			wantErr: true,
		},
		{
			name:    "2",
			fields:  fields{config: sync.Config{}, log: logging.GetLoggerFactory().Get()},
			args:    args{path: "2", data: "", conn: con},
			wantErr: true,
		},
		{
			name:    "3",
			fields:  fields{config: sync.Config{}, log: logging.GetLoggerFactory().Get()},
			args:    args{path: "3", data: "", conn: con},
			wantErr: true,
		},
		{
			name:    "4",
			fields:  fields{config: sync.Config{}, log: logging.GetLoggerFactory().Get()},
			args:    args{path: "4", data: "", conn: con},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := zookeeper{
				config: tt.fields.config,
				log:    tt.fields.log,
			}
			if err := z.send(tt.args.path, tt.args.data, tt.args.conn); (err != nil) != tt.wantErr {
				t.Errorf("zookeeper.send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_zookeeper_listen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := make(chan zk.Event, 1)
	c <- zk.Event{}
	con := mock.NewMockConnection(ctrl)
	con.EXPECT().GetW("1").Return([]byte("1"), &zk.Stat{}, c, errors.New("Error"))

	con.EXPECT().GetW("2").Return([]byte("2"), &zk.Stat{}, c, zk.ErrNoNode)
	con.EXPECT().Exists("2").Return(false, &zk.Stat{}, errors.New("Error"))

	con.EXPECT().GetW("3").Return([]byte("3"), &zk.Stat{}, c, nil)

	type fields struct {
		config sync.Config
		log    logging.Logger
	}
	type args struct {
		path string
		conn Connection
		c    chan sync.Response
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    string
	}{
		{
			name:    "1",
			fields:  fields{config: sync.Config{}, log: logging.GetLoggerFactory().Get()},
			args:    args{path: "1", conn: con, c: make(chan sync.Response, 1)},
			wantErr: true,
		},
		{
			name:    "2",
			fields:  fields{config: sync.Config{}, log: logging.GetLoggerFactory().Get()},
			args:    args{path: "2", conn: con, c: make(chan sync.Response, 1)},
			wantErr: true,
		},
		{
			name:    "3",
			fields:  fields{config: sync.Config{}, log: logging.GetLoggerFactory().Get()},
			args:    args{path: "3", conn: con, c: make(chan sync.Response, 1)},
			wantErr: false,
			want:    "3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := zookeeper{
				config: tt.fields.config,
				log:    tt.fields.log,
			}
			z.listen(tt.args.path, tt.args.conn, tt.args.c)

			r := <-tt.args.c

			if (r.Error != nil) != tt.wantErr {
				t.Errorf("zookeeper.listen() error = %v, wantErr %v", r.Error, tt.wantErr)
				return
			}

			if r.Data != tt.want {
				t.Errorf("zookeeper.listen() data = %v, want %v", r.Data, tt.want)
			}

		})
	}
}
