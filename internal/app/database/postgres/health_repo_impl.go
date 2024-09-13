package postgres

func (p *Postgres) PingDB() error {
	return p.db.Ping()
}
