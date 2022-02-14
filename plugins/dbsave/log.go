package dbsave

import (
	"context"
	"golog-db/config"
	"strings"
	"sync"
	"time"

	"github.com/VenomPCPL/golog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _opts = options.InsertMany().SetOrdered(false)

type LogService struct {
	logger  *golog.Logger
	db      *mongo.Client
	options *config.Options
	buff    []interface{}
	m       sync.RWMutex
}

func (v *LogService) loop() {
	for {
		time.Sleep(v.options.FlushEvery)
		v.m.Lock()
		if len(v.buff) > 0 {
			v.logger.Debug().Send("Flushing logs...")
			_, err := v.db.Database(v.options.Database).Collection(v.options.AppName).InsertMany(context.Background(), v.buff, _opts)
			if err != nil {
				v.logger.Error().Send("Error while flushing logs: %v", err.Error())
			} else {
				if v.options.Verbose {
					v.logger.Debug().Send("Logs flushed successfully (%v logs)", len(v.buff))
				}
			}
			v.buff = nil
		}
		v.m.Unlock()
	}
}

func (v *LogService) hook(msg *golog.Message, data []byte, _ []byte) {
	v.m.Lock()
	defer v.m.Unlock()
	v.buff = append(v.buff, Message{
		Message: data,
		Level:   msg.Level(),
		Time:    time.Now().UnixMilli(),
		Tags:    v.options.Tags,
	})
check:
	if int64(len(v.buff)) > v.options.MaxBufferSize {
		v.buff = v.buff[1:]
		goto check
	}
}

func NewLoggingService(log *golog.Logger, op ...config.Option) (service *LogService, err error) {
	opts := &config.Options{Database: "golog-db", AppName: "app", Verbose: true, MaxBufferSize: 5000, FlushEvery: time.Second * 5, MaxLogs: 20_000, MaxSize: 1024 * 1024 * 10}
	for i := range op {
		op[i](opts)
	}
	service = &LogService{}
	log2 := log.Module("golog-db")
	service.logger = log2
	service.options = opts
	service.db, err = mongo.NewClient(opts.DatabaseOptions)
	if err != nil {
		return
	}
	if err = service.db.Connect(context.Background()); err != nil {
		return
	}
	if opts.Verbose {
		log2.Debug().Send("Connected to database")
		log2.Info().Send("Listening for logs!")
	}
	_opts := options.CreateCollection().SetCapped(true).SetMaxDocuments(opts.MaxLogs).SetSizeInBytes(opts.MaxSize)
	err = service.db.Database(opts.Database).CreateCollection(context.Background(), opts.AppName, _opts)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return
		}
		err = nil
	}

	go service.loop()
	log.WriteHook(service.hook)
	return
}
