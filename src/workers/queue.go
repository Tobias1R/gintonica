package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	cnf "github.com/Tobias1R/gintonica/config"
	"github.com/Tobias1R/gintonica/pkg/helpers"
	"github.com/guntenbein/goconcurrency/errworker"
	"github.com/streadway/amqp"
)

const (
	n           = 2
	statusOK    = 0
	statusError = 1
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}

type RunningTask struct {
	id      string
	Order   int64
	Channel string
	Status  string
	Data    []byte
}

type worker struct {
	function              func(task *RunningTask) error
	queueName             string
	queue                 amqp.Queue
	pid                   int64
	status                string
	messageChannel        <-chan amqp.Delivery
	tasks                 map[string]*RunningTask
	isAlive               bool
	requestStop           bool
	connection            *amqp.Connection
	channel               *amqp.Channel
	receivedMessages      []string
	totalReceivedMessages int64
}

type queue struct {
	name    string
	workers []worker
}

type IRunningTask interface {
	GetStatus() error
	MarshalBinary() ([]byte, error)
	UnmarshalBinary(data []byte) error
	SetStatus(status string) error
	UpdateRedis() error
}

func (t *RunningTask) UpdateRedis() error {
	rdb := getRedisClient()
	jsonData, _ := t.MarshalBinary()
	err := rdb.Set(context.Background(), t.id, jsonData, 0).Err()
	if err != nil {
		log.Printf("Error connecting to Redis : %s", err)
	}
	return err
}

func (t *RunningTask) SetStatus(status string) error {
	t.Status = status
	return t.UpdateRedis()
}

func (t RunningTask) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *RunningTask) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	return nil
}

type Worker interface {
	connect() *amqp.Connection
	receiveMessage(t *RunningTask) error
	Start() // return pid
	Stop()
	Register(t *RunningTask) error
	NewWorker(q string) worker
	Run(t *RunningTask, taskChannel chan<- RunningTask, ew errworker.ErrWorkgroup) error
	Info() interface{}
}

func (w *worker) Info() interface{} {
	type structPayload struct {
		Status     bool  `json:"status"`
		TotalTasks int64 `json:"totalTasks"`
	}
	i := structPayload{}
	i.Status = w.isAlive
	i.TotalTasks = w.totalReceivedMessages

	return i
}

func (w *worker) receiveMessage(t *RunningTask) error {
	w.totalReceivedMessages++
	w.receivedMessages = append(w.receivedMessages, t.id)
	return nil
}

func (w *worker) Run(t *RunningTask, taskChannel chan<- RunningTask, ew errworker.ErrWorkgroup) error {
	ew = errworker.NewErrWorkgroup(2, true)
	ew.Go(func() error {
		w.function(t)
		task := *t
		taskChannel <- task
		println("AFTER CALL")
		return nil
	})
	println("RUN FLAG")
	return ew.Wait()
}

func NewWorker(q string, target func(task *RunningTask) error) worker {
	// return a new worker for given queue q
	return worker{
		function:       target,
		queueName:      q,
		queue:          amqp.Queue{},
		pid:            0,
		status:         "",
		messageChannel: make(<-chan amqp.Delivery),
		tasks:          map[string]*RunningTask{},
	}
}

func (w *worker) connect() {

	var err error
	w.connection, err = amqp.Dial(cnf.AMQPConnectionURL)
	handleError(err, "Can't connect to AMQP")

	w.channel, err = w.connection.Channel()
	handleError(err, "Can't create a amqpChannel")

	w.queue, err = w.channel.QueueDeclare(w.queueName, true, false, false, false, nil)
	handleError(err, "Could not declare `"+w.queueName+"` queue")

	err = w.channel.Qos(1, 0, false)
	handleError(err, "Could not configure QoS")

	w.messageChannel, err = w.channel.Consume(
		w.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	handleError(err, "Could not register consumer")

}

func TestME(t *RunningTask) error {
	println("RECEIVED BODY", string(t.Data))
	println("TASK PID: ", os.Getpid())
	time.Sleep(10 * time.Second)
	t.SetStatus("stage 1")
	time.Sleep(10 * time.Second)
	t.SetStatus("stage 2")
	time.Sleep(10 * time.Second)
	t.SetStatus("stage 3")
	time.Sleep(10 * time.Second)
	t.SetStatus("INTERNAL DONE")
	return nil
}

func (w *worker) Register(t *RunningTask, channel string) error {
	// register task for worker
	w.tasks[channel] = t
	return nil
}

func (w *worker) Stop() {
	// init stoping process
	w.requestStop = true
}

func (w *worker) Start() {
	SetQueueControl(w)
	// start worker
	w.connect()
	println("WORKER PID: ", os.Getpid())

	stopChan := make(chan bool)
	errc := make(chan error)
	status := statusOK

	errGroup := sync.WaitGroup{}
	errGroup.Add(1)

	go func() {
		for err := range errc {
			status = statusError
			fmt.Printf("error processing the code: %s\n", err)
		}
		errGroup.Done()
	}()
	fmt.Println("worker status ", !w.connection.IsClosed(), string(status))
	taskGroup := sync.WaitGroup{}
	rdb := getRedisClient()
	for c, t := range w.tasks {
		taskGroup.Add(1)
		if w.requestStop {
			break
		}
		ew := errworker.NewErrWorkgroup(1, true)

		// rabbit consumer
		go func() {
			for {
				if w.requestStop {
					break
				}
				w.isAlive = true

				for d := range w.messageChannel {

					log.Printf("Received a message: %s", d.MessageId)
					t.Status = "PENDING"
					t.id = d.MessageId
					t.Data = d.Body
					jsonData, _ := t.MarshalBinary()
					err := rdb.Set(context.Background(), d.MessageId, jsonData, 0).Err()
					if err != nil {
						log.Printf("Error connecting to Redis : %s", err)
						break
					}
					// rabbitMQ ackowledegment
					if err := d.Ack(false); err != nil {
						log.Printf("Error acknowledging message : %s", err)
					} else {
						log.Printf("Acknowledged message")
						w.receiveMessage(t)
					}
				}

				time.Sleep(1 * time.Second)
			}
		}()
		// redis time
		go func() {
			for {
				if w.requestStop {
					break
				}
				for _, message := range w.receivedMessages {
					val, err := rdb.Get(context.Background(), message).Result()
					if err != nil {
						log.Printf("Message NOT stored in redis: %s", message)
					}
					var t RunningTask
					t.UnmarshalBinary([]byte(val))
					taskChannel := make(chan RunningTask, 1)

					// task run
					log.Printf("Message stored in redis: %s", message)
					if errRun := w.Run(&t, taskChannel, ew); errRun != nil {
						errc <- errRun
					}

					//println("TASK STATUS AFTER RUN", t.Status)

					time.Sleep(1 * time.Second)

					// delete from redis
					errDel := rdb.Del(context.Background(), message).Err()
					if errDel != nil {
						log.Printf("Couldn't delete message stored in redis: %s", message)
					}

					w.receivedMessages = helpers.RemoveStringFromArray(w.receivedMessages, message)
					close(taskChannel)
				}
			}
		}()
		log.Println("starting worker for queue: " + c)

	}
	taskGroup.Done()

	close(errc)
	errGroup.Wait()
	taskGroup.Wait()
	w.isAlive = false
	// Stop for program termination
	<-stopChan
	w.connection.Close()
}

type Queue interface {
	GetInstance(name string) queue
	StopWorker(w worker) (*interface{}, error)
	StartWorker(w worker) (*interface{}, error)
	AddWorker(w *worker) error
	Status()
}

func GetQueueInstance(name string) queue {
	return queue{
		name:    name,
		workers: []worker{},
	}
}

func (q *queue) AddWorker(w *worker) error {
	q.workers = append(q.workers, *w)
	return nil
}

var QueueControlObject worker

func SetQueueControl(w *worker) {
	QueueControlObject = *w
}
