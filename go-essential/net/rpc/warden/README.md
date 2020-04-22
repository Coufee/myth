# warden
rpc 服务结构

# 中间件
auth            用户验证  
backoff         指数退避算法（exponential backoff algorithm）来持续增加重试之间的延迟时间，直到达到最大限制  
balancer        负载均衡    
codec           请求加密解密    
credentials     证书  
hystrix         熔断器 
logger          日志  
metrics         服务信息收集  
prepare         预处理 
retelimit       限流器 
registry        服务注册和服务发现 
recovery        中间件调用恢复 
retry           重入  
stats           cpu和内存等状态 
tracing         链路追踪  
validator       rpc参数验证 
