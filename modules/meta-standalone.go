package modules

import (
	"log"
	"maunium.net/go/mautrix"
	"time"
)

type InitFunction func() error
type RunFunction func(client *mautrix.Client) error

type Module struct {
	Name  string
	Class PermissionClass
	Init  InitFunction
	Run   RunFunction
}

var modules []Module

func RegisterModule(module Module) {
	modules = append(modules, module)
}

func Init() error {
	log.Println("Loading modules...")
	for _, module := range modules {
		if module.Init == nil {
			continue
		}

		log.Printf("Init module %s...", module.Name)
		err := module.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

func Run(client *mautrix.Client) {
	log.Printf("Starting threads for modules")

	for _, module := range modules {
		if module.Run == nil {
			continue
		}

		log.Printf("Starting thread for %s", module.Name)
		go func(client *mautrix.Client) {
			for {
				err := module.Run(client)
				if err != nil {
					log.Printf("Thread for %s got error: %v", module.Name, err)
				} else {
					log.Printf("Thread for %s terminated without error", module.Name)
				}

				log.Printf("Wating 1s before restarting...")

				time.Sleep(time.Second * 1)
			}
		}(client)
	}
}
