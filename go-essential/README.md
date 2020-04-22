# go-essential
app         用于生成项目结构  
base        基于interface实现可插拔的组件功能  
container   容器实现，包括队列  
            池(slice 切片实现,list 列表实现)  
            group 减少对象创建 懒加载容器  
log         日志  
net         网络链接模块，包括broker(http)  
            warden(grpc)   
utils       简单函数工具  
                assert  类型错误  
                file    文件判断和获取  
                list    数据和切片类型转换和判断  
                number  转换成值类型，和数值截取  
                rand    随机数相关  
                string  字符串拼接，合并，转换，补全  
                verify  校验信息，如邮箱，qq，ip，后缀名，手机号等等  
helper      复杂函数工具   
                time    时间  
metadata    元数据 对象，资源等等信息索引
