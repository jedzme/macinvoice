package notification

type Notification interface {
	SendEmail(serverLocation string, message string)
}

type service struct {
	Config
}

func NewService(config Config) (Notification, error) {

	return &service{config}, nil
}

func (s *service) SendEmail(serverLocation string, message string) {

}
