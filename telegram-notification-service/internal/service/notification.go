package service

type notificationService struct {
	api TelegramAPI
}

func New(tgAPI TelegramAPI) *notificationService {
	return &notificationService{}
}
