package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"sync"
	"time"
	"transok/backend/consts"
	"transok/backend/utils/logger"

	"transok/backend/utils/mdns"

	"github.com/betamos/zeroconf"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DiscoverService struct {
	ctx    context.Context
	cancel context.CancelFunc
	server *zeroconf.Client
	ticker *time.Ticker
	done   chan bool
	id     string
}

var (
	discoverService *DiscoverService
	once            sync.Once
	discoveryType   = zeroconf.NewType("_transok._tcp")
)

func init() {
	discoveryType.Domain = "local."
}

func GetDiscoverService() *DiscoverService {
	once.Do(func() {
		discoverService = &DiscoverService{
			id: uuid.New().String(),
		}
	})
	return discoverService
}

/* Start starts listening for mDNS broadcasts */
func (s *DiscoverService) Start() error {
	logger.Info("Starting discovery...")

	discovery := zeroconf.New().Browse(func(e zeroconf.Event) {
		logger.Debug("Found device", zap.Any("entry", e))
		var jsonData string
		for _, txt := range e.Text {
			if len(txt) > 5 && txt[:5] == "data=" {
				jsonData = txt[5:]
				break
			}
		}

		if jsonData == "" {
			return
		}

		// Remove extra escape characters
		jsonData = strings.ReplaceAll(jsonData, `\"`, `"`)

		var data consts.DiscoverPayload

		// Parse JSON string into DiscoverPayload struct
		if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
			logger.Error("Failed to parse JSON data", zap.Error(err))
			return
		}

		// Skip if sender is self
		// if data.Sender == s.id {
		// 	continue
		// }

		logger.Info("Service discovered", zap.Any("data", data))
		mdns.GetDispatcher().Dispatch(data)
	},
		discoveryType)

	_, err := discovery.Open() //Start discovering
	if err != nil {
		logger.Error("Error discovering devices", zap.Error(err))
		return err
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	//TODO: Browsing ain't stopping
	/*
		entries := make(chan *zeroconf.ServiceEntry)
		//err = resolver.Browse(s.ctx, "_transok._tcp", "local.", entries)

		go func() {
			logger.Info("Waiting for service discovery...")

			for entry := range entries {
				logger.Debug("Received raw service entry", zap.Any("entry", entry))

				// Look for TXT record containing data
				var jsonData string
				for _, txt := range entry.Text {
					if len(txt) > 5 && txt[:5] == "data=" {
						jsonData = txt[5:]
						break
					}
				}

				if jsonData == "" {
					continue
				}

				// Remove extra escape characters
				jsonData = strings.ReplaceAll(jsonData, `\"`, `"`)

				var data consts.DiscoverPayload

				// Parse JSON string into DiscoverPayload struct
				if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
					logger.Error("Failed to parse JSON data", zap.Error(err))
					continue
				}

				// Skip if sender is self
				// if data.Sender == s.id {
				// 	continue
				// }

				logger.Info("Service discovered", zap.Any("data", data))

				// Use dispatcher to handle the message
				mdns.GetDispatcher().Dispatch(data)
			}
		}()
	*/
	return nil
}

/* Broadcast starts broadcasting the service */
func (s *DiscoverService) Broadcast(port int, payload consts.DiscoverPayload) error {
	// Build DiscoverPayload
	discoverPayload := consts.DiscoverPayload{
		Type:    payload.Type,
		Sender:  s.id,
		Payload: payload.Payload,
	}

	// Serialize as JSON
	jsonBytes, err := json.Marshal(discoverPayload)
	if err != nil {
		return fmt.Errorf("failed to serialize payload: %v", err)
	}

	// Use JSON string as a single TXT record
	txtRecords := []string{
		fmt.Sprintf("data=%s", string(jsonBytes)),
	}

	if s.server == nil {
		service := zeroconf.NewService(
			discoveryType,
			fmt.Sprintf("TransokService_%s", s.id),
			uint16(port),
		)
		service.Text = txtRecords

		publisher := zeroconf.New().Publish(service)
		_, err = publisher.Open() //Start broadcast
		if err != nil {
			logger.Error("Error broadcasting service", zap.Error(err))
			return err
		}

		s.server = publisher
	}

	s.server.Open() //Start broadcast

	go func() { //Stop broadcast when context is done
		<-s.ctx.Done()
		s.server.Close()
	}()

	/*
		// Modify service registration configuration
		server, err := zeroconf.Register(
			fmt.Sprintf("TransokService_%s", s.id),
			"_transok._tcp",
			"local.",
			port,
			txtRecords,
			nil,
		)
		if err != nil {
			logger.Error("Failed to register broadcast service", zap.Error(err))
			return err
		}

		// Ensure old instance is closed before saving new one
		if s.server != nil {
			s.server.Shutdown()
		}
		s.server = server

		logger.Info("Service started broadcasting",
			zap.String("id", s.id),
			zap.Int("port", port),
			zap.Any("payload", discoverPayload),
		)
	*/
	return nil
}

/* StartPeriodicBroadcast starts periodic broadcasting
func (s *DiscoverService) StartPeriodicBroadcast(port int, payload consts.DiscoverPayload, interval time.Duration) {
	if s.ticker != nil {
		return
	}

	s.ticker = time.NewTicker(interval)
	s.done = make(chan bool)

	go func() {
		_ = s.Broadcast(port, payload)

		for {
			select {
			case <-s.ticker.C:
				if s.server != nil {
					s.server.Close()
					s.server = nil
				}
				_ = s.Broadcast(port, payload)
			case <-s.done:
				if s.server != nil {
					s.server.Close()
					s.server = nil
				}
				return
			}
		}
	}()
}

// StopPeriodicBroadcast stops periodic broadcasting
func (s *DiscoverService) StopPeriodicBroadcast() {
	if s.ticker != nil {
		s.ticker.Stop()
		s.done <- true
		close(s.done)
		s.ticker = nil

		//		if s.server != nil {
		//			s.server.Close()
		//			s.server = nil
		//		}
	}
}
*/
// Stop method
func (s *DiscoverService) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}
