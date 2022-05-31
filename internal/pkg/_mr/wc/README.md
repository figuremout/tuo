# MapReduce
参考: [6.824 lab1](https://pdos.csail.mit.edu/6.824/labs/lab-mr.html)，
[别人写的lab1](https://github.com/Anarion-zuo/MIT-6.824/blob/119a48fffd31b973c6d1cc8474ad78d72abbfa0e/src/mr/worker.go#L190)，
[mapreduce整体设计 split shuffle](https://blog.csdn.net/qq_35283816/article/details/79719468)

- map任务将文本分割为`nReduce`个intermediate文件，喂给reduce任务
- `nReduce`是reduce任务的个数
- worker应将 X'th reduce 任务输出写入文件`mr-out-X`
- `mr-out-X`文件一行为每个reduce任务输出，格式为"%v %v"
- worker应把intermediate map 结果写入当前目录下的文件中，以让worker之后可以读取作为reduce任务输入（分布式应改进为不写入文件直接传输流式字节）
- worker一直向coodinator轮询请求任务
- map和reduce函数是通过go plugin运行时加载的
- 6.824 lab1是在所有worker共享同一个文件系统的基础上，但若需要worker运行在不同机器上需要一个像GFS的全局文件系统
- intermediate命名规则是`mr-X-Y`，X是map任务id，Y是reduce任务id
- 用json库将intermediate的k/v对存进文件（或者转化成字节流？）
- map中用ihash(key)函数来根据key(即为分割的word)选择reduce任务id，将intermediate
- coordinator必须是并发的，共享数据结构必须安全
- 用go的race detector
- worker有时需要wait，比如reduce任务只有在map任务全部结束才能开始，可以让worker周期轮询coordinator，或者grpc handler用sync.Cond进行wait
- 每个任务有timeout，超时未返回结果视为失败并重新分配给其他worker
- 要有crash recovery能力，如map reduce时exit了
- 为避免文件写时被其他人读，先写入tmp文件，完全写入后再重命名