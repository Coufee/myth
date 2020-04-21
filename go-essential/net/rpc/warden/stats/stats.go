package stats

import (
	"context"
	"google.golang.org/grpc"
)

func Stats() grpc.UnaryServerInterceptor  {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		//var cpustat cpu.Stat
		//cpu.ReadStat(&cpustat)
		//if cpustat.Usage != 0 {
		//	trailer := gmd.Pairs([]string{nmd.CPUUsage, strconv.FormatInt(int64(cpustat.Usage), 10)}...)
		//	grpc.SetTrailer(ctx, trailer)
		//}
		return
	}
}
