package serviceimpl

import (
	"log"
	"time"

	messagebroker "github.com/NiflheimDevs/dyslexics-clock/internal/domain/message_broker"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"
	"github.com/go-co-op/gocron/v2"
)

type PublisherService struct {
	AlarmEventPublisher messagebroker.AlarmEventPublisher
	TimePublisher       messagebroker.TimePublisher
	RequestPublisher    messagebroker.RequestPublisher
}

func NewPublisherService(alarmEventPublisher messagebroker.AlarmEventPublisher, timePublisher messagebroker.TimePublisher, requestPublisher messagebroker.RequestPublisher) *PublisherService {
	p := &PublisherService{
		AlarmEventPublisher: alarmEventPublisher,
		TimePublisher:       timePublisher,
		RequestPublisher:    requestPublisher,
	}

	sched, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}
	sched.NewJob(gocron.DurationJob(time.Hour*24), gocron.NewTask(p.PublishTime()))
	sched.Start()
	return p
}

func (p *PublisherService) PublishTime() error {
	err := p.TimePublisher.PublishTime(time.Now().UTC().UnixMilli())
	if err != nil {
		log.Println("error publishing time", err)
		return err
	}
	return nil
}

func (p *PublisherService) PublishAlarm(deviceID string, alarm *model.Alarm) error {
	err := p.AlarmEventPublisher.PublishCreate(deviceID, alarm)
	if err != nil {
		log.Println("error creating alarm", err)
		return err
	}
	return nil
}

func (p *PublisherService) PublishAlarmUpdate(deviceID string, alarm *model.Alarm) error {
	err := p.AlarmEventPublisher.PublishUpdate(deviceID, alarm)
	if err != nil {
		log.Println("error updating alarm", err)
		return err
	}
	return nil
}

func (p *PublisherService) PublishAlarmDelete(deviceID string, alarmID string) error {
	err := p.AlarmEventPublisher.PublishDelete(deviceID, alarmID)
	if err != nil {
		log.Println("error deleting alarm", err)
		return err
	}
	return nil
}

func (p *PublisherService) PublishVolume(deviceID string, volume uint) error {
	err := p.RequestPublisher.PublishVolume(deviceID, volume)
	if err != nil {
		log.Println("error publishing volume", err)
		return err
	}
	return nil
}

func (p *PublisherService) PublishColor(deviceID string, color string) error {
	err := p.RequestPublisher.PublishColor(deviceID, color)
	if err != nil {
		log.Println("error publishing color", err)
		return err
	}
	return nil
}
