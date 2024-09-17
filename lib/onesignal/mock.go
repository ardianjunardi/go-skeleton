package onesignal

type (
	configMock struct{}
)

var (
	PushNotificationMock func(xPlayer, title, description string) error
)

func NewConfigMock() OneSignal {
	return &configMock{}
}

func (lib *configMock) PushNotification(xPlayer, title, description string) error {
	if PushNotificationMock != nil {
		return PushNotificationMock(xPlayer, title, description)
	}
	return nil
}
