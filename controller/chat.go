package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"ginchat/utils"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Node struct {
	Conn          *websocket.Conn //连接
	Addr          string          //客户端地址
	FirstTime     uint64          //首次连接时间
	HeartbeatTime uint64          //心跳时间
	LoginTime     uint64          //登录时间
	DataQueue     chan []byte     //消息
	GroupSets     set.Interface   //好友 / 群
}

func (node *Node) Heartbeat(currentTime uint64) {
	node.HeartbeatTime = currentTime
	return
}

func (node *Node) IsHeartBeatTimeout(currenttime uint64) bool {
	if node.HeartbeatTime+viper.GetUint64("timeout.HeartbeatMaxTime") < currenttime {
		fmt.Println("心跳超时", node)
		node.Conn.Close()
		return true
	}
	return false
}

type Message struct {
	gorm.Model
	UserId     int64  //发送者
	TargetId   int64  //接受者
	Type       int    //发送类型  1私聊  2群聊  3心跳
	Media      int    //消息类型  1文字 2表情包 3语音 4图片 /表情包
	Content    string //消息内容
	CreateTime uint64 //创建时间
	ReadTime   uint64 //读取时间
	Pic        string
	Url        string
	Desc       string
	Amount     int //其他数字统计
}

var clientMap map[int64]*Node = make(map[int64]*Node)
var rwLocker sync.RWMutex

func Chat(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	Id := query.Get("userId")
	userId, err := strconv.ParseInt(Id, 10, 64)
	if err != nil {
		panic(err)
	}

	//token校验
	isvalid := true

	//升级为websocket连接
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalid
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentTime := uint64(time.Now().Unix())
	//初始化client对象
	node := &Node{
		Conn:          conn,
		Addr:          conn.RemoteAddr().String(),
		HeartbeatTime: currentTime,
		LoginTime:     currentTime,
		DataQueue:     make(chan []byte, 50),
		GroupSets:     set.New(set.ThreadSafe),
	}

	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	go sendProc(node)
	go recvProc(node)
	SetUserOnlineInfo("online_"+Id, []byte(node.Addr), time.Duration(viper.GetInt("timeout.RedisOnlineTime"))*time.Hour)
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("[ws]sendProc >>>> msg :", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		message := &Message{}
		err = json.Unmarshal(data, message)
		if err != nil {
			fmt.Println(err)
		}

		if message.Type == 3 {
			currentTime := uint64(time.Now().Unix())
			node.Heartbeat(currentTime)
		} else {
			dispatch(data)
			broadMsg(data)
			fmt.Println("[ws] recvProc <<<<< ", string(data))
		}
	}
}

func dispatch(data []byte) {
	message := &Message{}
	json.Unmarshal(data, message)
	switch message.Type {
	case 1: //单发
		fmt.Println("dispatch  data :", string(data))
		sendMsg(message.TargetId, data)
	case 2: //群发
		sendGroupMsg(message.TargetId, data)
	}
}

func sendMsg(TargetId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[TargetId]
	jsonMsg := &Message{}
	json.Unmarshal(msg, jsonMsg)
	ctx := context.Background()
	TargetIdstr := strconv.Itoa(int(TargetId))
	UserIdstr := strconv.Itoa(int(jsonMsg.UserId))
	jsonMsg.CreateTime = uint64(time.Now().Unix())
	r, err := utils.REDIS.Get(ctx, "online_"+UserIdstr).Result()
	if err != nil {
		fmt.Println(err)
	}
	if r != "" {
		if ok {
			fmt.Println("sendMsg >>> userID: ", TargetId, "  msg:", string(msg))
			node.DataQueue <- msg
		}
	}
	var key string
	if TargetId > jsonMsg.UserId {
		key = "msg_" + UserIdstr + "_" + TargetIdstr
	} else {
		key = "msg_" + TargetIdstr + "_" + UserIdstr
	}
	res, err := utils.REDIS.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println(err)
	}
	score := float64(cap(res)) + 1
	ress, err := utils.REDIS.ZAdd(ctx, key, &redis.Z{score, msg}).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ress)

	if err != nil {
		fmt.Println(err)
	}

}

func sendGroupMsg(targetId int64, msg []byte) {
	fmt.Println("开始群发消息")
	userIds := SearchUserByGroupId(uint(targetId))
	for i := 0; i < len(userIds); i++ {
		//排除给自己的
		if targetId != int64(userIds[i]) {
			sendMsg(int64(userIds[i]), msg)
		}
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}
func init() {
	go udpSendProc()
	go udpRecvProc()
	fmt.Println("init goroutine ")
}
func udpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: viper.GetInt("udp"),
	})
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendChan:
			fmt.Println("udpSendProc  data :", string(data))
			_, err := conn.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func udpRecvProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: viper.GetInt("port.udp"),
	})
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	for {
		buf := []byte{}
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("udpRecvProc  data :", string(buf[0:n]))
		dispatch(buf[0:n])
	}
}

func CleanTimeoutConnection(param interface{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("cleanConnection err", r)
		}
	}()

	currentTime := uint64(time.Now().Unix())
	for i := range clientMap {
		node := clientMap[i]
		if node.IsHeartBeatTimeout(currentTime) {
			fmt.Println("心跳超时..... 关闭连接：", node)
			node.Conn.Close()
		}
	}
	fmt.Println("clean completed")

}


