学习笔记
go语言error的理解：
1、简单；
2、考虑失败，面向错误编程，先判断错误，然后再继续业务逻辑；
3、没有隐藏的控制流，完全交给你来控制error；
4、error are values；

日志记录与错误无关且对调试没有帮助的信息应被视为噪音；记录日志是因为某些东西失败了：
错误要被日志记录；
应用程序处理错误，保证100%完整性；
之后不再报告当前错误；

Warp Error的一些总结：
github.com/pkg/errors
errors.New() 、 errors.Errorf()

errors.Wrap(err,"open failed!") //有堆栈信息
errors.Wrapf()

errors.WithMessage(err,"not found!") //无堆栈信息
errors.Cause(err) //获取根因
%+v 可以打印出堆栈信息

使用errors.Cause(err) 获取根因然后再与sentinel error判定。
1、wrap是应用代码可以选择的策略，具有最高可重用性的包直接返回根因（基础库，第三方库，标准库）
2、如果函数/方法调用返回处理了error，就不用往上抛error，应该return nil