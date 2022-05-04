package handlers

import (
	"chatRoom/models"
	"chatRoom/services"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

func StartRoom(ctx iris.Context, ws services.IWebsocketService) {
	ch, err := services.MqConn.Channel()

	log.Println("Start Room ws address %p", &ws)
	defer ch.Close()

	if err != nil {
		ctx.JSON(iris.Map{
			"Message": "Error",
		})
		return
	}
	var user models.UserClaim

	json.Unmarshal([]byte(ctx.Values().Get("user").(string)), &user)

	log.Println("User: " + user.Name + " Start a room")

	queue, err := ch.QueueDeclare(user.Name, false, true, false, false, nil)
	if err != nil {
		ctx.JSON(iris.Map{
			"Message": "Fail to declare queue",
		})
		return
	}

	var key = uuid.New().String()
	ch.QueueBind(queue.Name, "*"+"."+key, models.ChatRoomExchange, false, nil)

	messages, err := ch.Consume(queue.Name, uuid.New().String(), true, false, false, false, nil)
	if err != nil {
		ctx.JSON(iris.Map{
			"Message": "Fail to declare consumer",
		})
		return
	}

	go func(claim *models.UserClaim) {
		for d := range messages {
			var msg ChatMessage
			json.Unmarshal(d.Body, &msg)
			log.Println("Start Room Received User: " + msg.UserName + " message")
			if msg.UserName == claim.Name {
				continue
			}
			ws.WriteMessage(msg.UserName + ":" + msg.Message)
		}
	}(&user)

	ws.Start(ctx.ResponseWriter(), ctx.Request())

	ws.WriteMessage(key)

	wsCh := ws.StartReceive()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for body := range wsCh {
			msg := string(body)
			handlerWsMessage(ch, user, msg, key)
		}
	}()

	wg.Wait()

	log.Println("Exist start room")
}

func JoinRoom(ctx iris.Context, ws services.IWebsocketService) {
	roomId := ctx.Params().Get("roomId")

	log.Println("Join Room ws address %p", &ws)
	ch, err := services.MqConn.Channel()
	defer ch.Close()

	if err != nil {
		ctx.JSON(iris.Map{
			"Message": "Error",
		})
		return
	}

	ws.Start(ctx.ResponseWriter(), ctx.Request())

	var user models.UserClaim

	json.Unmarshal([]byte(ctx.Values().Get("user").(string)), &user)

	log.Println("User: " + user.Name + " Join a room")

	queue, err := ch.QueueDeclare(user.Name, false, true, false, false, nil)
	if err != nil {
		ctx.JSON(iris.Map{
			"Message": "Fail to declare queue",
		})
		return
	}

	ch.QueueBind(queue.Name, "*"+"."+roomId, models.ChatRoomExchange, false, nil)

	messages, err := ch.Consume(queue.Name, uuid.New().String(), true, false, false, false, nil)
	if err != nil {
		ctx.JSON(iris.Map{
			"Message": "Fail to declare consumer",
		})
		return
	}

	go func(claim *models.UserClaim) {
		for d := range messages {
			var msg ChatMessage
			json.Unmarshal(d.Body, &msg)
			log.Println("Join Room Received User: " + msg.UserName + " message")
			if msg.UserName == user.Name {
				continue
			}
			ws.WriteMessage(msg.UserName + ": " + msg.Message)
		}
	}(&user)

	wsCh := ws.StartReceive()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for body := range wsCh {
			msg := string(body)
			handlerWsMessage(ch, user, msg, roomId)
		}
	}()

	wg.Wait()

	log.Println("Exist join room")
}

func handlerWsMessage(ch *amqp.Channel, user models.UserClaim, msg string, key string) {
	data, _ := json.Marshal(&ChatMessage{
		UserId:   user.ID,
		UserName: user.Name,
		Message:  msg,
	})
	ch.Publish(
		models.ChatRoomExchange,
		user.ID+"."+key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
}

type ChatMessage struct {
	UserId   string
	UserName string
	Message  string
}
