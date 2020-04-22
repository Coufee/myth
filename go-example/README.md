# go-example
http_server 处理前端进来的http请求  
middle      处理中间请求  
server      处理最终请求  

请求链路     client -> http_server -> middle -> server  
                                                |  
                                                v  
                                       server模块处理请求，并回复  
                                                |  
                                                v  
返回值链路   client <- http_server <- middle <- server  
