## 同步转异步

### 如何判定服务商已经崩溃，不允许使用课程上的判断机制，你需要设计一个新的判断机制，并且解释这种判定机制的决策理由。

- 服务商做了限流，通过错误码告诉我们他那边限流了。
- 服务商承受不住大量请求，服务降级了，告诉我们服务目前不可用，没有超时。
- 服务商做的不好，或者种种原因没有成功限流和降级，导致我们的请求超时，超时意味着最后我们不知道服务商到底有没有成功给我们的用户发送短信，我们把这种情况视为失败，需要转异步发送，即使用户最后收到2条重复消息，也比收不到强。

以上失败的情况，我们转为异步处理，且我们需要将这次失败的服务商从服务中踢掉 x（配置）分钟，之后再加回来。


### 控制异步重试次数，转储到数据库之后，可以重试 N 次，重试间隔你可以自由决策。

老师提供的 lua 脚本里，验证码会在 redis 里存储 10 分钟，这其实是和服务商无关的。

其实，重试发送成功后，可以再次设置一下，把有效期重新设置为 10 分钟。 

否则如果在第 10 分钟重试，会导致用户即使收到验证码，也无法使用这个验证码来登录了，当然用户能不能等这么久也是一个问题，但是这能解决重试的情况有效期会比告诉用户的10分钟少的问题，比如在1分钟的时候重试，有效期变成9分钟的问题。

当然这时候 redis 即使设置失败，也是基本可以接受的，除了极端条件下，用户收到短信，验证码确已经过期了。

所以重试次数、重试间隔完全可以由配置决定，唯一需要注意的是配置的合理性，最坏情况下用户需要等 N * 间隔 的时间。


### 写明这种容错机制适合什么场景，并且有什么优缺点。

优点：
如果服务商都可用，我们的用户会分配给不同的服务商，即使我们的系统某个时间有大量请求，不会给一个服务商造成压力，因为大量请求被不同服务商消化了。
短信发送失败有各种各样的原因，比如发送数量超过限额，不一定是服务商挂了，这时候我们可以利用异步重试的机制换个服务商重试，重试间隔较短的情况下，用户端基本无感。

缺点：
如果服务商少，容易造成，没有服务商可用。


### 针对提出的缺点，写出后续的改进方案。

提供告警机制，服务商少于 x 个，给相关人员发送钉钉、企业微信、飞书消息等。