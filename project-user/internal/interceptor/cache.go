package interceptor

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"test.com/project-common/encrypts"
	"test.com/project-grpc/user/login"
	"test.com/project-user/internal/dao"
	"test.com/project-user/internal/repo"
	"time"
)

type CacheInterceptor struct {
	cache    repo.Cache
	cacheMap map[string]any
}

type CacheRespOption struct {
	path   string
	typ    any
	expire time.Duration
}

func New() *CacheInterceptor {
	cacheMap := make(map[string]any)
	cacheMap["/login.service.v1.LoginService/MyOrgList"] = &login.OrgListResponse{}
	cacheMap["/login.service.v1.LoginService/FindMemInfoById"] = &login.MemberMessage{}
	return &CacheInterceptor{cache: dao.Rc, cacheMap: cacheMap}
}

func (c *CacheInterceptor) Cache() grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		respType := c.cacheMap[info.FullMethod]
		if respType == nil {
			return handler(ctx, req)
		}
		//先查询是否有缓存 有的话 直接返回 无 先请求 然后存入缓存
		con, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		marshal, _ := json.Marshal(req)
		cacheKey := encrypts.Md5(string(marshal))
		respJson, _ := c.cache.Get(con, info.FullMethod+"::"+cacheKey)
		if respJson != "" {
			json.Unmarshal([]byte(respJson), &respType)
			zap.L().Info(info.FullMethod + " 走了缓存")
			return respType, nil
		}
		resp, err = handler(ctx, req)
		bytes, _ := json.Marshal(resp)
		c.cache.Put(con, info.FullMethod+"::"+cacheKey, string(bytes), 5*time.Minute)
		zap.L().Info(info.FullMethod + " 放入缓存")
		return
	})
}
