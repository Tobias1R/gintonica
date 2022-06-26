package workers

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	cnf "github.com/Tobias1R/gintonica/config"
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
	Function func(j []byte, task *RunningTask) error
	Order    int64
	Channel  string
	Status   string
}

type worker struct {
	queueName      string
	queue          amqp.Queue
	pid            int64
	status         string
	messageChannel <-chan amqp.Delivery
	tasks          map[string]RunningTask
	isAlive        bool
	requestStop    bool
	connection     *amqp.Connection
	channel        *amqp.Channel
}

type taskQueue struct {
	name    string
	workers []worker
	broker  *amqp.Connection
}

type Task interface {
	GetStatus(t *RunningTask) string
	Run(j []byte, taskChannel chan<- RunningTask, ew errworker.ErrWorkgroup) error
}

func (t *RunningTask) Run(j []byte, taskChannel chan<- RunningTask, ew errworker.ErrWorkgroup) error {
	ew = errworker.NewErrWorkgroup(2, true)
	ew.Go(func() error {

		t.Function(j, t)
		task := *t
		taskChannel <- task
		println("AFTER CALL")
		return nil
	})
	println("RUN FLAG")
	return ew.Wait()
}

type Worker interface {
	connect() *amqp.Connection
	Start() // return pid
	Stop()
	Register(t *RunningTask) error
	NewWorker(q string) worker
}

func NewWorker(q string) worker {
	// return a new worker for given queue q
	return worker{
		queueName:      q,
		queue:          amqp.Queue{},
		pid:            0,
		status:         "",
		messageChannel: make(<-chan amqp.Delivery),
		tasks:          map[string]RunningTask{},
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

func TestME(j []byte, t *RunningTask) error {
	println("RECEIVED BODY", string(j))
	println("TASK PID: ", os.Getpid())
	time.Sleep(5 * time.Second)
	t.Status = "INTERNAL DONE"
	return nil
}

func (w *worker) Register(t RunningTask, channel string) error {
	w.tasks[channel] = t
	return nil
}

func (w *worker) Stop() {
	w.requestStop = true
}

func (w *worker) Start() {
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
					taskChannel := make(chan RunningTask, 1)
					log.Printf("Received a message: %s", d.MessageId)
					t.Status = "RUNNING"

					if errRun := t.Run(d.Body, taskChannel, ew); errRun != nil {
						errc <- errRun
					}
					println("TASK STATUS AFTER RUN", t.Status)
					close(taskChannel)
					if err := d.Ack(false); err != nil {
						log.Printf("Error acknowledging message : %s", err)
						t.Status = "FAIL"
					} else {
						log.Printf("Acknowledged message")
						t.Status = "DONE"
					}
				}

				println("Worker alive", os.Getpid())
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
	// Stop for program termination
	<-stopChan
	w.connection.Close()
}

type Queue interface {
	Enqueue(task *interface{}) (*interface{}, error)
	StopWorker(w worker) (*interface{}, error)
	StartWorker(w worker) (*interface{}, error)
	Connect() (*amqp.Channel, error)
	AddWorker(w *worker) string
}
