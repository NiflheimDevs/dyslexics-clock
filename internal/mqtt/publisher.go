package mqtt

import "fmt"

func (c *Client) PublishAlarm(deviceID uint, action string, payload []byte) error {
	topic := fmt.Sprintf("server/%d/alarm/%s", deviceID, action)
	return c.Publish(topic, payload)
}

func (c *Client) PublishColor(deviceID uint, color string) error {
	topic := fmt.Sprintf("server/%d/color", deviceID)
	return c.Publish(topic, []byte(color))
}
