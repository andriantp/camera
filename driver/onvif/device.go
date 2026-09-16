package onvif

import usego "github.com/use-go/onvif"

func (c *Client) device(camera Camera) (*usego.Device, error) {
	return usego.NewDevice(usego.DeviceParams{
		Xaddr:      camera.Host,
		Username:   camera.Username,
		Password:   camera.Password,
		HttpClient: c.http,
	})
}
