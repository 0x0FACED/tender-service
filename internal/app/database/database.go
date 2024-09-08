package database

type Database interface {
	Connect() error

	BidRepository
	TenderRepository
}

type BidRepository interface {
}

type TenderRepository interface {
}
