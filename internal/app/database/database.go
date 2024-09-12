package database

type Database interface {
	Connect() error

	BidRepository
	TenderRepository
	UserRepository
}
