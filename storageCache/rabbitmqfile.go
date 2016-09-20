package storageCache


import(
	"fmt"
	"g"
	"github.com/streadway/amqp"
)




 
func SendRabbitMq(message []byte,ch *amqp.Channel) error{
	
	//body:=proxy_msg
	err:=ch.Publish(
		g.Config().Rabbitmq.Exchange,
		g.Config().Rabbitmq.RouteingKey,
		false,
		false,
	    amqp.Publishing {
    		ContentType: "text/plain",
    		Body:        message,
  		},
	)

	if err != nil {
		fmt.Println("publish err:",err)
		return err
	}
	fmt.Println("write message to rabbitmq success!!!")
	return nil

}