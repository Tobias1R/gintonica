package workers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	cnf "github.com/Tobias1R/gintonica/config"
	"github.com/Tobias1R/gintonica/pkg/helpers"
	"github.com/google/uuid"
	"github.com/guntenbein/goconcurrency/errworker"

	"github.com/streadway/amqp"
)

const (
	n                   = 2
	statusOK            = 0
	statusError         = 1
	TASK_STATUS_PENDING = "PENDING"
	TASK_STATUS_RUNNING = "RUNNING"
	TASK_STATUS_DONE    = "DONE"
	TASK_STATUS_ACK     = "ACK"
	TASK_STATUS_FAIL    = "FAIL"
	MaxWorkers          = 5
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}

type RunningTask struct {
	id             string
	Order          int64
	Channel        string
	Status         string
	Data           []byte
	RunnningWorker string
}

type worker struct {
	function              func(task *RunningTask) error
	pid                   int
	status                string
	tasks                 map[string]*RunningTask
	isAlive               bool
	requestStop           bool
	receivedMessages      []string
	totalReceivedMessages int64
	runs                  int64
	workerId              string
	amqpQueue             amqp.Queue
	connection            *amqp.Connection
	channel               *amqp.Channel
	queueName             string
	messageChannel        <-chan amqp.Delivery
	currentTask           string
	currentTaskStatus     string
}

type queue struct {
	name        string
	workers     []worker
	requestStop bool
}

type IRunningTask interface {
	GetStatus() error
	MarshalBinary() ([]byte, error)
	UnmarshalBinary(data []byte) error
	SetStatus(status string) error
	UpdateRedis() error
	SetWorker(w *worker) error
}

func (t *RunningTask) UpdateRedis() error {
	rdb := getRedisClient()
	defer rdb.Close()

	key, err := rdb.Get(context.Background(), t.id).Result()
	if err != nil {
		log.Printf("UPDATE Message NOT stored in redis: %s", key)
		return err
	}
	jsonData, _ := t.MarshalBinary()
	errS := rdb.Set(context.Background(), t.id, jsonData, 0).Err()
	if errS != nil {
		log.Printf("Error connecting to Redis : %s", errS)
	}
	return err
}

func (t *RunningTask) SetWorker(w *worker) error {
	t.RunnningWorker = w.workerId
	return t.UpdateRedis()
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
	//Start()
	Stop()
	RegisterTask(t *RunningTask) error
	NewWorker(q string) worker
	Run(t *RunningTask, taskChannel chan<- RunningTask, ew errworker.ErrWorkgroup) error
	Info() interface{}
	removeMessage(message string) error
	RunFromRedis(key *string) error
}

func (w *worker) removeMessage(message string) error {
	w.receivedMessages = helpers.RemoveStringFromArray(w.receivedMessages, message)
	return nil
}

func (w *worker) Info() interface{} {
	type structPayload struct {
		Status        string `json:"status"`
		TotalMessages int64  `json:"totalReceivedMessages"`
		Task          string `json:"currentTask"`
		TaskStatus    string `json:"currentTaskStatus"`
		Runs          int64  `json:"runs"`
		WorkerId      string `json:"workerId"`
	}
	i := structPayload{}
	i.Status = w.status
	i.TotalMessages = w.totalReceivedMessages
	i.Task = w.currentTask
	i.TaskStatus = w.currentTaskStatus
	i.Runs = w.runs
	i.WorkerId = w.workerId

	return i
}

func (w *worker) receiveMessage(t RunningTask) error {
	w.totalReceivedMessages++
	t.SetWorker(w)
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
	w.runs++
	return ew.Wait()
}

func NewWorker(q string, target func(task *RunningTask) error) worker {
	// return a new worker for given queue q
	return worker{
		function:              target,
		pid:                   0,
		status:                "INITIALIZING",
		tasks:                 map[string]*RunningTask{},
		isAlive:               false,
		requestStop:           false,
		receivedMessages:      []string{},
		totalReceivedMessages: 0,
		runs:                  0,
		workerId:              uuid.New().String(),
		amqpQueue:             amqp.Queue{},
		connection:            &amqp.Connection{},
		channel:               &amqp.Channel{},
		queueName:             q,
		messageChannel:        make(<-chan amqp.Delivery),
	}
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
	t.SetStatus(TASK_STATUS_DONE)
	return nil
}

func (w *worker) RegisterTask(t *RunningTask, channel string) error {
	// register task for worker
	w.tasks[channel] = t
	return nil
}

func (w *worker) Stop() {
	// init stoping process
	w.requestStop = true
}

func Start(w *worker) {
	w.status = "ON THE BEGNNING"

	// start worker
	w.connect()
	w.pid = os.Getpid()
	println("WORKER ", w.workerId, " PID: ", w.pid)

	stopChan := make(chan bool)
	errc := make(chan error)

	errGroup := sync.WaitGroup{}
	errGroup.Add(1)

	go func() {
		for err := range errc {
			fmt.Printf("error processing the code: %s\n", err)
		}
		errGroup.Done()
	}()

	taskGroup := sync.WaitGroup{}
	rdb := getRedisClient()
	for c, t := range w.tasks {
		taskGroup.Add(1)
		if w.requestStop {
			break
		}
		ew := errworker.NewErrWorkgroup(1, true)
		w.status = "RUNNING"
		// rabbit consumer
		go func() {
			// Ad infinitum
			for {
				println("heartbeat: ", w.workerId)
				if w.requestStop {
					break
				}
				w.isAlive = true

				for d := range w.messageChannel {
					if d.MessageId == "" {
						continue
					}
					log.Printf("worker %s Received message: %s", w.workerId, d.MessageId)
					t.Status = TASK_STATUS_PENDING
					t.id = d.MessageId
					t.Data = d.Body
					jsonData, _ := t.MarshalBinary()
					err := rdb.Set(context.Background(), d.MessageId, jsonData, 0).Err()
					if err != nil {
						log.Printf("Error connecting to Redis : %s", err)
						break
					}
					message := d.MessageId
					val, err := rdb.Get(context.Background(), message).Result()
					if err != nil {
						log.Printf("Message NOT stored in redis: %s", message)
					} else {
						var t RunningTask
						t.UnmarshalBinary([]byte(val))
						taskChannel := make(chan RunningTask, 1)
						w.receiveMessage(t)
						t.id = d.MessageId
						w.currentTask = d.MessageId
						w.currentTaskStatus = t.Status
						// task run
						if t.Status == TASK_STATUS_PENDING {
							w.currentTaskStatus = TASK_STATUS_RUNNING
							log.Printf("Message stored in redis: %s", message)
							if errRun := w.Run(&t, taskChannel, ew); errRun != nil {
								errc <- errRun
							}
						}

						if t.Status == TASK_STATUS_DONE {

							//w.receivedMessages = helpers.RemoveStringFromArray(w.receivedMessages, message)
							if err := d.Ack(false); err != nil {
								log.Printf("Error acknowledging message : %s", err)
							} else {
								log.Printf("Acknowledged message")

								// delete from redis
								errDel := rdb.Del(context.Background(), message).Err()
								if errDel != nil {
									log.Printf("Couldn't delete message stored in redis: %s", message)
								}
							}
						}

						if t.Status == TASK_STATUS_FAIL {
							if err := d.Reject(false); err != nil {
								log.Printf("Error rejecting message : %s", err)
							} else {
								log.Printf("Message rejected")

								// delete from redis
								errDel := rdb.Del(context.Background(), message).Err()
								if errDel != nil {
									log.Printf("Couldn't delete message stored in redis: %s", message)
								}
							}
						}

						//println("TASK STATUS AFTER RUN", t.Status)

						close(taskChannel)
						time.Sleep(1 * time.Second)
						w.currentTask = ""
						w.currentTaskStatus = ""

					}
				}

				time.Sleep(1 * time.Second)
			}
		}()

		log.Println("starting worker for queue: " + c)

	}

	taskGroup.Done()

	close(errc)
	errGroup.Wait()
	taskGroup.Wait()
	w.isAlive = false
	//w.connection.Close()
	// Stop for program termination
	<-stopChan

}

type Queue interface {
	GetInstance(name string) queue
	StopWorker(w *worker) error
	StartWorker(w *worker) error
	AddWorker(w *worker) error
	Status()
	StartConsumer(channelName string) error
	//connect() *amqp.Connection
	StartWorkers() error
	Start() error
	Info() interface{}
	Register(queueName string, totalWorkers int, target func(task *RunningTask) error)
}

func (q *queue) Info() interface{} {
	type structPayload struct {
		Total   int           `json:"totalWorkers"`
		Workers []interface{} `json:"workers"`
	}
	i := structPayload{}
	for _, w := range q.workers {
		i.Workers = append(i.Workers, w.Info())
	}
	i.Total = len(i.Workers)
	return i
}

func (q *queue) Start() error {
	// not checking errors here, i'm in lazy mode

	//go q.StartConsumer(q.name)
	go q.StartWorkers()
	SetQueueControl(q)
	return nil
}

func (q *worker) connect() {

	var err error
	q.connection, err = amqp.Dial(cnf.AMQPConnectionURL)
	handleError(err, "Can't connect to AMQP")

	q.channel, err = q.connection.Channel()
	handleError(err, "Can't create a amqpChannel")

	q.amqpQueue, err = q.channel.QueueDeclare(q.queueName, true, false, false, false, nil)
	handleError(err, "Could not declare `"+q.queueName+"` queue")

	err = q.channel.Qos(1, 0, false)
	handleError(err, "Could not configure QoS")

	q.messageChannel, err = q.channel.Consume(
		q.amqpQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	handleError(err, "Could not register consumer")

}

func GetQueueInstance(name string) queue {
	return queue{
		name:    name,
		workers: []worker{},

		requestStop: false,
	}
}

func (q *queue) Register(queueName string, totalWorkers int, target func(task *RunningTask) error) {
	var i int = 0
	if totalWorkers > MaxWorkers {
		totalWorkers = MaxWorkers
	}

	for i < totalWorkers {
		w := NewWorker(queueName, target)
		t := RunningTask{
			Order:          0,
			Channel:        q.name,
			Status:         "PENDING",
			Data:           []byte(""),
			RunnningWorker: w.workerId,
		}
		w.RegisterTask(&t, queueName)
		q.AddWorker(w)
		i++
	}
}

func (q *queue) AddWorker(w worker) error {
	if len(q.workers) > MaxWorkers {
		return errors.New("max workers reached")
	}

	q.workers = append(q.workers, w)
	return nil
}

func (q *queue) StartWorkers() error {
	stopChan := make(chan bool)
	errc := make(chan error)

	errGroup := sync.WaitGroup{}
	errGroup.Add(1)

	go func() {
		for err := range errc {
			fmt.Printf("error processing the code: %s\n", err)
		}
		errGroup.Done()
	}()
	//fmt.Println("worker status ", !w.q.connection.IsClosed(), string(status))
	taskGroup := sync.WaitGroup{}

	for i, w := range q.workers {
		w.status = string(statusOK)
		taskGroup.Add(1)
		go Start(&q.workers[i])
		//time.Sleep(5 * time.Second)
		log.Println("starting worker: " + w.workerId)

	}

	taskGroup.Done()

	close(errc)
	errGroup.Wait()
	taskGroup.Wait()

	// Stop for program termination
	<-stopChan
	return nil
}

var QueueControlObject *queue

func SetQueueControl(q *queue) {
	QueueControlObject = q
}
