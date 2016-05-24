package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/medcl/mysqlbeat/config"
)

type Mysqlbeat struct {
	beatConfig *config.Config
	done       chan struct{}
	period     time.Duration
	client     publisher.Client
}

// Creates beater
func New() *Mysqlbeat {
	return &Mysqlbeat{
		done: make(chan struct{}),
	}
}

/// *** Beater interface methods ***///

func (bt *Mysqlbeat) Config(b *beat.Beat) error {

	// Load beater beatConfig
	err := b.RawConfig.Unpack(&bt.beatConfig)
	if err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}

	return nil
}

func (bt *Mysqlbeat) Setup(b *beat.Beat) error {

	// Setting default period if not set
	if bt.beatConfig.Mysqlbeat.Period == "" {
		bt.beatConfig.Mysqlbeat.Period = "1s"
	}

	bt.client = b.Publisher.Connect()

	var err error
	bt.period, err = time.ParseDuration(bt.beatConfig.Mysqlbeat.Period)
	if err != nil {
		return err
	}

	return nil
}

func (bt *Mysqlbeat) Run(b *beat.Beat) error {
	logp.Info("mysqlbeat is running! Hit CTRL-C to stop it.")

	ticker := time.NewTicker(bt.period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
			"counter":    counter,
		}
		bt.client.PublishEvent(event)
		logp.Info("Event sent")
		counter++
	}
}

func (bt *Mysqlbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (bt *Mysqlbeat) Stop() {
	close(bt.done)
}
