package handlers

import (
	"chatRoom/models"
	"chatRoom/services"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/streadway/amqp"
)

func StartRoom(ctx iris.Context, ws services.IWebsocketService) {
	ch, err := services.MqConn.Channel()
	defer ch.Close()

	if err != nil {
		ctx.JSON(iris.Map{
			"Message": "Error",
		})
		return
	}
	var user models.UserClaim

	json.Unmarshal([]byte(ctx.Values().Get("user").(string)), &user)

	queue, err := ch.QueueDeclare(user.ID, false, true, false, false, nil)
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

	go func() {
		for d := range messages {
			var msg ChatMessage
			json.Unmarshal(d.Body, &msg)
			if msg.UserId == user.ID {
				continue
			}
			ws.WriteMessage(msg.Message)
		}
	}()

	ws.Start(ctx.ResponseWriter(), ctx.Request())

	ws.WriteMessage(key)

	ws.StartReceive(func(service services.IWebsocketService, bytes []byte) {
		msg := string(bytes)
		data, _ := json.Marshal(&ChatMessage{
			UserId:  user.ID,
			Message: msg,
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
	})
}

func JoinRoom(ctx iris.Context, ws services.IWebsocketService) {
	roomId := ctx.Params().Get("roomId")

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

	queue, err := ch.QueueDeclare(user.ID, false, true, false, false, nil)
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

	go func() {
		for d := range messages {
			var msg ChatMessage
			json.Unmarshal(d.Body, &msg)
			if msg.UserId == user.ID {
				continue
			}
			ws.WriteMessage(msg.Message)
		}
	}()

	ws.StartReceive(func(service services.IWebsocketService, bytes []byte) {
		msg := string(bytes)
		data, _ := json.Marshal(&ChatMessage{
			UserId:  user.ID,
			Message: msg,
		})
		ch.Publish(
			models.ChatRoomExchange,
			user.ID+"."+roomId,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        data,
			},
		)
	})

}

type ChatMessage struct {
	UserId  string
	Message string
}
