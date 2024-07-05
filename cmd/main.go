package main

import (
	"fmt"

	"github.com/ecol-master/sharing-wh-machines/internal/dbs/postgres"
	"github.com/ecol-master/sharing-wh-machines/internal/logger"
	"github.com/ecol-master/sharing-wh-machines/internal/repositories/machines"
)

func main() {
	err := logger.Setup()
	if err != nil {
		fmt.Println("can not initialize logger: ", err)
		return
	}

	cfg := postgres.Config{
		Addr:     "0.0.0.0",
		Port:     5400,
		User:     "postgres",
		Password: "postgres",
		DB:       "sharing_machines",
	}

	conn, err := postgres.New(cfg)
	if err != nil {
		panic(err)
	}

	machineId := "1FG5689"
	machinesRepo := machines.NewRepository(conn)
	err = machinesRepo.InsertMachine(machineId)
	fmt.Println(err)

	machine, err := machinesRepo.SelectMachine(machineId)
	fmt.Println(machine, err)
}
