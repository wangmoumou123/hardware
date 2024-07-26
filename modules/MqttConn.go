package modules

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"sync"
	"time"
)

type MqttConnI interface {
	MsgHandle(client MQTT.Client, msg MQTT.Message)
	OnConnectHandle(client MQTT.Client)
	DisConnectHandle(client MQTT.Client, err error)
}
type msgCallback func(msg, tp string)

type MqttConn struct {
	Addr        string
	Topic       []string
	ClientId    string
	Conn        MQTT.Client
	MsgCallBack func(msg, tp string)
}

func (mq MqttConn) MsgHandle(client MQTT.Client, message MQTT.Message) {
	msg := string(message.Payload())
	tp := message.Topic()
	//log.Println(msg, "======", tp)
	if mq.MsgCallBack != nil {
		mq.MsgCallBack(msg, tp)
	}
}

func (mq MqttConn) OnConnectHandle(client MQTT.Client) {
	for _, tp := range mq.Topic {
		client.Subscribe(tp, 1, nil)
		log.Println("subject===>", tp)
	}
}

func (mq MqttConn) DisConnectHandle(client MQTT.Client, err error) {
	log.Println("==disconnected====")
}

func MqttConnInit(addr string, topics []string, msgCallback msgCallback, clientId ...string) *MqttConn {
	var cId string
	if len(clientId) == 0 {
		cId = "go"
	} else {
		cId = clientId[0]
	}
	mqttConn := &MqttConn{Addr: addr, Topic: topics, ClientId: cId, MsgCallBack: msgCallback}
	opts := MQTT.NewClientOptions()

	opts.AddBroker(mqttConn.Addr)
	opts.SetClientID(mqttConn.ClientId)
	opts.SetDefaultPublishHandler(mqttConn.MsgHandle)
	opts.SetOnConnectHandler(mqttConn.OnConnectHandle)
	opts.SetConnectionLostHandler(mqttConn.DisConnectHandle)

	//opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	// MQTT 3.1.1
	opts.SetProtocolVersion(4)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(time.Second * 5)
	// 创建 MQTT 客户端
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("Error connecting to broker: %v\n", token.Error())
		os.Exit(1)
	}
	mqttConn.Conn = client
	return mqttConn
}

func (mq *MqttConn) MqttPublish(msg string, topic ...string) {
	var tp string
	if len(topic) == 0 {
		tp = mq.Topic[0]
	} else {
		tp = topic[0]
	}
	mq.Conn.Publish(tp, 0, false, msg)

}

func (mq *MqttConn) RunAlways(wg *sync.WaitGroup) {
	defer func() {
		mq.Conn.Disconnect(200)
		wg.Done()
	}()
	for {
		mq.MqttPublish("xixihaha")
		time.Sleep(time.Second * 2)
	}
}
