package grpcserver

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"git.ucloudadmin.com/bigdata/MAXIR/src/common/grace"
	"git.ucloudadmin.com/bigdata/MAXIR/src/common/log"
	pbModel "git.ucloudadmin.com/bigdata/MAXIR/src/common/proto/model"
	pbService "git.ucloudadmin.com/bigdata/MAXIR/src/common/proto/service"
	"git.ucloudadmin.com/bigdata/MAXIR/src/common/utils"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	// ignoreResponseMethods 针对返回数据较大请求不进行日志打印
	ignoreResponseMethods = []interface{}{}
)

var (
	ErrInvalidServiceState = errors.New("no access, invalid service state")
	ErrInvalidResponse     = errors.New("Invalid response")
)

func Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {

	defer func() {
		if !checkResponseIgnore(info.FullMethod) {
			log.Info(ctx, "gRPC outgoing response", log.Fields{
				"res":       resp,
				"rpcMethod": info.FullMethod,
			})
		}

		grace.Recover(ctx)
	}()

	// 打印请求信息
	log.Info(ctx, "gRPC incoming request", log.Fields{
		"req":       req,
		"rpcMethod": info.FullMethod,
	})

	// 参数 validation
	code := validateParams(ctx, req)
	if code != pbModel.RetCode_SUCCESS {
		return genResp(ctx, code)
	}

	// 将 SessionID 等重新写入到 grpc context Outgoing 中，否则再次调用第三方服务时参数无法传递过去
	// 因为 grpc 服务端需要从 IncomingContext 获取 Value，默认会将客户端 OutgoingContext Value 重置到 IncomingContext 中
	ctx = utils.DuplicateContext(ctx)

	// 继续处理请求
	resp, err = handler(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := handleRespMessage(ctx, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func genResp(ctx context.Context, code pbModel.RetCode) (*pbService.CommonRes, error) {
	resp := &pbService.CommonRes{
		RetCode: code,
	}
	if err := handleRespMessage(ctx, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func handleRespMessage(ctx context.Context, resp interface{}) error {
	messageField := reflect.ValueOf(resp).Elem().FieldByName("Message")
	if !messageField.IsValid() {
		return ErrInvalidResponse
	}
	if messageField.String() == "" {
		retcodeField := reflect.ValueOf(resp).Elem().FieldByName("RetCode")
		if !retcodeField.IsValid() {
			return ErrInvalidResponse
		}
		retcode := retcodeField.Int()
		messageField.SetString(pbModel.RetCode_name[int32(retcode)])
	}
	return nil
}

// checkResponseIgnore 检查是否要忽略返回结果
func checkResponseIgnore(fullMethod string) bool {
	return checkMethodContain(fullMethod, ignoreResponseMethods)
}

func checkMethodContain(fullMethod string, methods []interface{}) bool {
	var method string
	lastSeprationIndex := strings.LastIndex(fullMethod, "/")
	if lastSeprationIndex != -1 {
		method = fullMethod[lastSeprationIndex+1:]
	}

	for _, m := range methods {
		if strings.HasSuffix(getFunctionName(m), fmt.Sprintf("%s-fm", method)) {
			return true
		}
	}

	return false
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
