
我们在做数据库操作的时候，假设在 dao 层中遇到一个 sql.ErrNoRows ，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码?

我认为应该wrap，因为这样处理的时候可以获取到调用堆栈。
数据层出现`sql.ErrNoRows` 在DAO层被捕获封装成NoSuchUser 错误。其他错误使用OtherError 表示。
方便应用进行判断，因为有可能业务层并不知道底层是SQL还是NoSQL。

业务层可以使用errors.As(err, &NoSuchUser{}) 做进一步判断，比如某些接口是允许查询不到用户，某些接口把查询不到用户作为一种异常。