package ha

import (
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
)

import (
	"fmt"
	"time"
)

func main() {
	var (
		config        clientv3.Config
		client        *clientv3.Client
		lease         clientv3.Lease
		leaseResp     *clientv3.LeaseGrantResponse
		leaseId       clientv3.LeaseID
		leaseRespChan <-chan *clientv3.LeaseKeepAliveResponse
		err           error
	)

	//客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:23790"},
		DialTimeout: 5 * time.Second,
	}
	//建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//上锁（创建租约，自动续租）
	lease = clientv3.NewLease(client)

	//设置一个ctx取消自动续租
	ctx, cancleFunc := context.WithCancel(context.TODO())

	//设置10秒租约（过期时间）
	if leaseResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}
	//拿到租约id
	leaseId = leaseResp.ID

	//自动续租（不停地往管道中扔租约信息）
	if leaseRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
	}
	//启动一个协程去监听
	go listenLeaseChan(leaseRespChan)

	get_lock := func() {
		//业务处理
		kv := clientv3.NewKV(client)
		//创建事务
		txn := kv.Txn(context.TODO())
		txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"), "=", 0)).
			Then(clientv3.OpPut("/cron/lock/job9", "xxx", clientv3.WithLease(leaseId))).
			Else(clientv3.OpGet("/cron/lock/job9")) //否则抢锁失败

			//提交事务
		if txtResp, err := txn.Commit(); err != nil {
			fmt.Println(err)
			return
		} else {
			//判断是否抢锁
			if !txtResp.Succeeded {
				fmt.Println("锁被占用：", string(txtResp.Responses[0].GetResponseRange().Kvs[0].Value))
				return
			} else {
				fmt.Println("锁被成功获取")
			}
		}
	}
	get_lock()

	time.Sleep(3 * time.Second)

	fmt.Println("处理任务")

	//释放锁（停止续租，终止租约）
	fmt.Println(time.Now())
	cancleFunc()                          //函数退出取消自动续租
	lease.Revoke(context.TODO(), leaseId) //终止租约（去掉过期时间）
	lease.Revoke(context.TODO(), -1)      //终止租约（去掉过期时间）

	fmt.Println("主函数休眠")
	time.Sleep(20 * time.Second)
	fmt.Println(time.Now())
}

func listenLeaseChan(leaseRespChan <-chan *clientv3.LeaseKeepAliveResponse) {
	var (
		leaseKeepResp *clientv3.LeaseKeepAliveResponse
	)
	for {
		select {
		case leaseKeepResp = <-leaseRespChan:
			if leaseKeepResp == nil {
				fmt.Println("心跳结束")
				goto END
			} else {
				fmt.Println("保活: ", leaseKeepResp.ID, " -> ", time.Now())
			}
		}
		fmt.Println("for loop")
	}
END:
	fmt.Println("心跳退出")
}
