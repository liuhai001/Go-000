学习笔记

Week03 作业题目：

1. 基于 errgroup 实现一个 http server 的启动和关闭;
2.linux signal 信号的注册和处理;
3.要保证能够一个退出，全部注销退出。


作业分析：
1、用errgroup 启动多个httpserver;
2、监听Linux signal信号，让http server 都优雅退出；
3、http server 一个退出，全部都要退出；

