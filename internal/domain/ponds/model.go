package ponds

type PondType struct {
	ID     int64  `db:"id"`
	FarmID int64  `db:"farm_id"`
	Name   string `db:"name"`
}

type PondFarmType struct {
	PondType
	FarmName string `db:"farm_name"`
}
