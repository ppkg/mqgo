package mqsync

import (
	"fmt"
	"testing"
)

func TestT(t *testing.T) {

}
func TestPublish(t *testing.T) {
	type args struct {
		model SyncMqInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "订阅",
			args: args{
				model: SyncMqInfo{
					Exchange: "datacenter",
					Queue:    "dc-sz-test-mqsync",
					RouteKey: "dc-sz-test-mqsync",
					Request:  "content1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Publish(tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConsume(t *testing.T) {
	/* go func() {
		for {
			go TestPublish(t)
			time.Sleep(time.Second * 2)
		}
	}() */

	type args struct {
		name     string
		key      string
		exchange string
		autoAck  bool
		fun      func(request string) (string, error)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "消费",
			args: args{
				name:     "dc-sz-test-mqsync",
				key:      "dc-sz-test-mqsync",
				exchange: "datacenter",
				autoAck:  true,
				fun: func(request string) (string, error) {
					fmt.Println(request)
					return "测试成功", nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Consume(tt.args.name, tt.args.key, tt.args.exchange, tt.args.fun)
		})
	}
}
