package integrationtests

import (
	"bytes"
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func connectToRabbitMQ() *RabbitMQClient {
	conn, err := amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
	Expect(err).NotTo(HaveOccurred())

	channel, err := conn.Channel()
	Expect(err).NotTo(HaveOccurred())

	return &RabbitMQClient{
		conn:    conn,
		channel: channel,
	}
}

func (c *RabbitMQClient) Close() {
	Expect(c.channel.Close()).NotTo(HaveOccurred())
	Expect(c.conn.Close()).NotTo(HaveOccurred())
}

func (c *RabbitMQClient) NotifyStatus(eventID, status string) error {
	message := map[string]string{"eventID": eventID, "status": status}
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = c.channel.QueueDeclare(
		"notification_status",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return c.channel.Publish(
		"",
		"notification_status",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// Получение сообщения из очереди.
func (c *RabbitMQClient) GetMessage(queueName string) (*amqp091.Delivery, error) {
	msgs, err := c.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	select {
	case msg := <-msgs:
		return &msg, nil
	}
}

var eventJSON = `{
	"event": {
		"title": "Test Event",
		"description": "This is a test event",
		"duration": "3600",
		"event_time": "2024-12-01T14:52:09Z",
		"notify_before": "60",
		"user_id": "51eb1ccc-1e1e-4132-9472-572c15c8c6fc"
	}
}`

var _ = Describe("Event Management", func() {
	baseURL := "http://calendar:8888/v1/events"

	Context("Adding events", func() {
		It("should successfully add an event", func() {
			resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer([]byte(eventJSON)))
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should fail with invalid date", func() {
			invalidEventJSON := `{
				"event": {
					"title": "Invalid Event",
					"description": "Invalid time format",
					"event_time": {
						"seconds": "invalid-date"
					}
				}
			}`
			resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer([]byte(invalidEventJSON)))
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})
})

var _ = Describe("Event Listing", func() {
	baseURL := "http://calendar:8888/v1/events"

	Context("Retrieving events", func() {
		BeforeEach(func() {
			addEvent(baseURL, eventJSON)
			addEvent(baseURL, eventJSON)
			addEvent(baseURL, eventJSON)
		})

		It("should return events for a specific day", func() {
			resp, err := http.Get(baseURL + "/day?date=2024-12-01T14:52:09Z")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should return events for a specific week", func() {
			resp, err := http.Get(baseURL + "/week?date=2024-12-01T14:52:09Z")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should return events for a specific month", func() {
			resp, err := http.Get(baseURL + "/month?date=2024-12-01T14:52:09Z")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})
})

var _ = Describe("Notification Sending", func() {
	var rabbitMQ *RabbitMQClient

	BeforeEach(func() {
		rabbitMQ = connectToRabbitMQ()
	})

	AfterEach(func() {
		rabbitMQ.Close()
	})

	Context("Sending notification statuses", func() {
		It("should send and receive a notification message", func() {
			err := rabbitMQ.NotifyStatus("test-event-id", "completed")
			Expect(err).NotTo(HaveOccurred())

			msg, err := rabbitMQ.GetMessage("notification_status")
			Expect(err).NotTo(HaveOccurred())
			var message map[string]string
			Expect(json.Unmarshal(msg.Body, &message)).To(Succeed())
			Expect(message["eventID"]).To(Equal("test-event-id"))
			Expect(message["status"]).To(Equal("completed"))
		})
	})
})

func addEvent(baseURL, eventJSON string) {
	resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer([]byte(eventJSON)))
	Expect(err).NotTo(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
}
