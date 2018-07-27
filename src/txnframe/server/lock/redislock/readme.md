## redis  lock 

redis lock 实现提供Prelock及lock， 两种锁的区别在于， Prelock不可以共享， lock可以根据返回值中的SubTxnID和LocksubTxnID实现锁在同一TxnID中共享

每一种锁提供加锁和解锁方式， 额外提根据TxnID解锁下面所有SubTxnID锁的资源

#### 结构体介绍

锁的数据结构

``` golang
type Lock struct {
	//  当前锁的标记。 是开启事务得到的事务ID， 事务主ID，不可以为空
	TxnID string `json:"txnID"`

	// 锁资源的子事务ID， 可以为空， 为空将自动生成改项目
	SubTxnID string `json:"subTxnID"`

	// 被锁资源的标记或者名字
	LockName string `json:"lockName"`

    // 获取锁的最长等待时间
	Timeout time.Duration `json:"timeout"`

    // 锁的创建时间
	Createtime time.Time `json:"createTime"`
}
```

lock锁返回的数据结构

``` golang
type LockResult struct {
	// 锁资源的子事务ID， 
	SubTxnID string `json:"subTxnID"`

	// 获取lock 传入的TxnID事务中是否有子事务拥有锁，
	Locked bool `json:"locked"`

	// 拥有锁的子事务ID， 及时第一lock资源的子事务ID
	LockSubTxnID string `json:"lockSubTxnID"`
}
```

### 锁的正确性


使用redis SETNX 命令来保证锁的正确性，
SETNX 只有在key不存在的时候， 才可以设置成功。若key 已经存在，则 SETNX 不做任何动作。 
同时SETNX是原子操作，在并发操作上不会出现问题。



### 锁的逻辑正确性

 - redis key 的介绍 

Prelock, lock两种锁， 每一种有两个key， 一个key是string， 一个key是hash

string的key是锁， 这个key中的内容就是Lock 结构体
hash的key是一个Txn(主机事务)下锁与子事务的关系，   hash的key是锁名， 内容是这个锁对应子事务 Lock 结构体数据

 - redis key 的使用顺序

 1. 先用SETNX 产生锁的key， 如果出现错误返回， 上层逻辑会在Timeout重试
 2. 如果SETNX成功，然后新加锁与子事务的关系。如果新加失败，首先通知补偿方法，需要清除锁。并且返回错误
 3. 如果SETNX返回无任何操作，则表示锁已经存在，通过与key的内容比较返回是否拥有锁。 Prelock， lock 返回不一样， 比较内容方法也可以不同。

 - 补偿方法

 补偿方法有主机补偿和被动补偿。 

 1. 主机补偿收到通知后，会更具补偿的类型不同，做对应的补偿操作
 2. 被动补偿， 定时出发， 通过redis的Scan 来遍历所有是事务锁前缀的key，通过锁，锁的关系来判读是否需要做补偿

 