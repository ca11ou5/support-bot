package http

type Handler struct {
	AdminService
	TelegramService
}

func NewHandler() *Handler {
	return &Handler{}
}

type AdminService interface {
}

type TelegramService interface {
}
