package mqsync

import (
	"testing"
	"time"
)

func TestT(t *testing.T) {

}
func TestPublish(t *testing.T) {
	// engine := newDataCenterConn()
	// engine.Exec("TRUNCATE TABLE mq_publish_info")
	// engine.Exec("TRUNCATE TABLE mq_consume_record")
	// go TestConsume(t)

	type args struct {
		model MqPublishInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "推送",
			args: args{
				model: MqPublishInfo{
					Exchange: "datacenter",
					RouteKey: "dc-sz-test-mqsync",
					Request:  "content1",
				},
			},
		},
	}
	for _, tt := range tests {
		for i := 0; i < 10; i++ {
			t.Run(tt.name, func(t *testing.T) {
				if err := Publish(tt.args.model); (err != nil) != tt.wantErr {
					t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
			// time.Sleep(time.Second)
		}
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
				fun: func(request string) (string, error) {
					return "测试成功", nil
				},
			},
		},
		{
			name: "消费2",
			args: args{
				name:     "dc-sz-test-mqsync2",
				key:      "dc-sz-test-mqsync",
				exchange: "datacenter",
				fun: func(request string) (string, error) {
					return "测试成功", nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go Consume(tt.args.name, tt.args.key, tt.args.exchange, tt.args.fun)
		})
	}
	time.Sleep(time.Second * 9999)
}
