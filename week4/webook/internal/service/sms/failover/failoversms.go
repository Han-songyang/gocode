package failover

import (
	"context"
	"errors"
	"log"
	"webook/internal/service/sms"
)

type FailOverSMSService struct {
	svcs []sms.Service
}

func NewFailOverSMSService(svcs []sms.Service) *FailOverSMSService {
	return &FailOverSMSService{
		svcs: svcs,
	}
}

// Send 全部轮询完了都没成功，那就说明所有的服务商挂了
// 出现这种情况的更大可能性是自己网络崩了
// 缺点：
//
//	每次都从头开始轮询，绝大多数请求会在svcs[0] 就成功，负载不均衡。
//	如果 svcs 有几十个，轮询都很慢
func (f *FailOverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplId, args, numbers...)
		if err == nil {
			// 发送成功
			return nil
		}
		log.Println(err)
	}
	return errors.New("轮询了所有服务商，但是都发送失败了")
}
