<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [总结限流，熔断，降级的常用方式，重试的注意事项，负载均衡的常用方式。](#%E6%80%BB%E7%BB%93%E9%99%90%E6%B5%81%E7%86%94%E6%96%AD%E9%99%8D%E7%BA%A7%E7%9A%84%E5%B8%B8%E7%94%A8%E6%96%B9%E5%BC%8F%E9%87%8D%E8%AF%95%E7%9A%84%E6%B3%A8%E6%84%8F%E4%BA%8B%E9%A1%B9%E8%B4%9F%E8%BD%BD%E5%9D%87%E8%A1%A1%E7%9A%84%E5%B8%B8%E7%94%A8%E6%96%B9%E5%BC%8F)
      - [限流](#%E9%99%90%E6%B5%81)
        - [过载保护](#%E8%BF%87%E8%BD%BD%E4%BF%9D%E6%8A%A4)
        - [分布式限流](#%E5%88%86%E5%B8%83%E5%BC%8F%E9%99%90%E6%B5%81)
        - [熔断](#%E7%86%94%E6%96%AD)
        - [客户端流控](#%E5%AE%A2%E6%88%B7%E7%AB%AF%E6%B5%81%E6%8E%A7)
      - [降级](#%E9%99%8D%E7%BA%A7)
      - [重试](#%E9%87%8D%E8%AF%95)
      - [负载均衡](#%E8%B4%9F%E8%BD%BD%E5%9D%87%E8%A1%A1)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

总结限流，熔断，降级的常用方式，重试的注意事项，负载均衡的常用方式。
===============

#### 限流

在一段时间内，定义某个客户或应用可以接收或处理多少个请求的技术。

##### 过载保护

- 漏斗桶/令牌桶：
    - 需预设指标，实际效果取决于限流阈值是否合理，但因限流阈值的计算、集群机器增减等原因，总体来说太被动，不能快速适应变化。
    - 针对单个节点，无法分布式限流。
    - `不推荐`。

- 计算系统临近过载时的峰值吞吐无哦为限流的阈值来进行流量控制，达到系统保护的目的。
    - 服务器临近过载时，主动抛弃一定量的负载，目标是自保。
    - 在系统稳定的前提下，保持系统的吞吐。常见做法：利特尔法则。
    - CPU、内存作为信号量进行节流。
    - 队列管理：度列长度、LIFO。
    - 可控延迟算法：CoDel。

##### 分布式限流

控制某个应用全局的流量，而非针对单个节点
常见使用方式：从redis获取单个quota。

- 缺点：
    - 单个大流量的接口，使用redis容易产生热点。
    - pre-request 模式对性能有一定影响，高频的网络往返。
    
- 优化：

    - 每次心跳后，异步批量获取 quota，可以大大减少请求 redis 的频次，获取完以后本地消费，基于令牌桶拦截。
    - 每次申请的配额需要手动设定静态值略欠灵活，比如每次要20，还是50。
    - 基于单个节点按需申请，初次使用默认值，一旦有过去历史窗口的数据，可以基于历史窗口数据进行 quota 请求，避免出现不公平的现象。
    - 给每个用户设置限制，划分资源遵循最大最小公平分享（Max-Min Fairness）原则。直观上，公平分享分配给每个用户想要的可以满足的最小需求，然后将没有使用的资源均匀的分配给需要‘大资源’的用户。
    - 每个接口配置阈值，运营工作繁重，最简单的我们配置服务级别 quota，更细粒度的，我们可以根据不同重要性设定 quota，优先保证重要性高的请求可用性。

##### 熔断

某个用户超过资源配额时，后端任务会快速拒绝请求，返回“配额不足”的错误，但是拒绝回复仍然会消耗一定资源。有可能后端忙着不停发送拒绝请求，导致过载。

- Gutter
    - 基于熔断的 gutter kafka ，用于接管自动修复系统运行过程中的负载，这样只需要付出10%的资源就能解决部分系统可用性问题。
    - 我们经常使用 failover 的思路，但是完整的 failover 需要翻倍的机器资源，平常不接受流量时，资源浪费。高负载情况下接管流量又不一定完整能接住。所以这里核心利用熔断的思路，是把抛弃的流量转移到 gutter 集群，如果 gutter 也接受不住的流量，重新回抛到主集群，最大力度来接受。

##### 客户端流控

positive feedback: 用户总是积极重试，访问一个不可达的服务。

- 客户端需要限制请求频次，retry backoff 做一定的请求退让。

- 可以通过接口级别的 error_details，挂载到每个 API 返回的响应里。


#### 降级

通过降级回复来减少工作量，或者丢弃不重要的请求。而且需要了解哪些流量可以降级，并且有能力区分不同的请求。我们通常提供降低回复的质量来答复减少所需的计算量或者时间。

我们自动降级通常需要考虑几个点：
- 确定具体采用哪个指标作为流量评估和优雅降级的决定性指标（如，CPU、延迟、队列长度、线程数量、错误等）。
- 当服务进入降级模式时，需要执行什么动作？
- 流量抛弃或者优雅降级应该在服务的哪一层实现？是否需要在整个服务的每一层都实现，还是可以选择某个高层面的关键节点来实现？


注意点：
- 优雅降级不应该被经常触发 - 通常触发条件现实了容量规划的失误，或者是意外的负载。

- 演练，代码平时不会触发和使用，需要定期针对一小部分的流量进行演练，保证模式的正常。
- 应该足够简单。

常见处理：
- UI 模块化，非核心模块降级。
    - BFF 层聚合 API，模块降级。
- 页面上一次缓存副本。
- 默认值、热门推荐等。
- 流量拦截 + 定期数据缓存(过期副本策略)。

处理策略
- 页面降级、延迟服务、写/读降级、缓存降级
- 抛异常、返回约定协议、Mock 数据、Fallback 处理

#### 重试
当请求返回错误（例: 配额不足、超时、内部错误等），对于 backend 部分节点过载的情况下，倾向于立刻重试，但是需要留意重试带来的流量放大:
- 限制重试次数和基于重试分布的策略（重试比率: 10%）。
- 随机化、指数型递增的重试周期: exponential ackoff + jitter。
- client 测记录重试次数直方图，传递到 server，进行分布判定，交由 server 判定拒绝。
- 只应该在失败的这层进行重试，当重试仍然失败，全局约定错误码“过载，无须重试”，避免级联重试。

#### 负载均衡
目标：
- 均衡的流量分发。
- 可靠的识别异常节点。
- scale-out，增加同质节点扩容。
- 减少错误，提高可用性。

常见方法：JSQ（最闲轮训）负载均衡算法带来的问题，缺乏的是服务端全局视图，因目标需要综合考虑：负载+可用性。

使用 p2c 算法，随机选取的两个节点进行打分，选择更优的节点:
- 选择 backend：CPU，client：health、inflight、latency 作为指标，使用一个简单的线性方程进行打分。
- 对新启动的节点使用常量惩罚值（penalty），以及使用探针方式最小化放量，进行预热。
- 打分比较低的节点，避免进入“永久黑名单”而无法恢复，使用统计衰减的方式，让节点指标逐渐恢复到初始状态(即默认值)。
- 当前发出去的请求超过了 predict lagtency，就会加惩罚。

指标计算结合 moving average，使用时间衰减，计算vt = v(t-1) * β + at * (1-β) ，β 为若干次幂的倒数即: Math.Exp((-span) / 600ms)


