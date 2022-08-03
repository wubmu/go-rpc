
## 服务器注册(service register)
- 通过反射实现服务注册功能
- 在服务端实现服务调用，代码约 150 行

### 1. 结构体映射为服务

RPC 框架的一个基础能力是：像调用本地程序一样调用远程服务。那如何将程序映射为服务呢？那么对 Go 来说，这个问题就变成了如何将结构体的方法映射为服务。

对 `net/rpc` 而言，一个函数需要能够被远程调用，需要满足如下五个条件：

- the method’s type is exported. – 方法所属类型是导出的。
- the method is exported. – 方式是导出的。
- the method has two arguments, both exported (or builtin) types. – 两个入参，均为导出或内置类型。
- the method’s second argument is a pointer. – 第二个入参必须是一个指针。
- the method has return type error. – 返回值为 error 类型。


```go
func (t *T) MethodName(argType T1, replyType *T2) error
```

假设客户端发过来一个请求，包含 ServiceMethod 和 Argv。

```
{
    "ServiceMethod"： "T.MethodName"
    "Argv"："0101110101..." // 序列化之后的字节流
}
```

通过 “T.MethodName” 可以确定调用的是类型 T 的 MethodName，如果硬编码实现这个功能，很可能是这样：

```go
switch req.ServiceMethod {
    case "T.MethodName":
        t := new(t)
        reply := new(T2)
        var argv T1
        gob.NewDecoder(conn).Decode(&argv)
        err := t.MethodName(argv, reply)
        server.sendMessage(reply, err)
    case "Foo.Sum":
        f := new(Foo)
        ...
}
```

也就是说，如果使用硬编码的方式来实现结构体与服务的映射，那么每暴露一个方法，就需要编写等量的代码。那有没有什么方式，能够将这个映射过程自动化呢？可以借助反射。

通过反射，我们能够非常容易地获取某个结构体的所有方法，并且能够通过方法，获取到该方法所有的参数类型与返回值。例如：

```go
func main() {
var wg sync.WaitGroup
typ := reflect.TypeOf(&wg)
for i := 0; i < typ.NumMethod(); i++ {
method := typ.Method(i)
argv := make([]string, 0, method.Type.NumIn())
returns := make([]string, 0, method.Type.NumOut())
// j 从 1 开始，第 0 个入参是 wg 自己。
for j := 1; j < method.Type.NumIn(); j++ {
argv = append(argv, method.Type.In(j).Name())
}
for j := 0; j < method.Type.NumOut(); j++ {
returns = append(returns, method.Type.Out(j).Name())
}
log.Printf("func (w *%s) %s(%s) %s",
typ.Elem().Name(),
method.Name,
strings.Join(argv, ","),
strings.Join(returns, ","))
}
}
```


