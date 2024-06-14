package postgres

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func connectToPostgres() {
	conn, err := pgx.
}
