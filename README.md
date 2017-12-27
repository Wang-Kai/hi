# hi

A implement of Resolver for gRPC , which achieve Register and Discovery base on etcd .

hi 实现在 etcd 中维护各微服务的地址注册，注册格式为：

key   ==> 微服务名/ip:port

value ==> ip:port

每个微服务在 startup 时注册其 微服务名&地址

最终的结果类似如下例子：

> /hi/order/1.2.3.4:8000  ===> 1.2.3.4:8000
> 
> /hi/order/1.2.3.4:8001  ===> 1.2.3.4:8001

当 gRPC 发请求的时候，会根据负载算法从 server list 中取出一个来发送 request