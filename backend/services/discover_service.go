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
	ctx     context.Context
	cancel  context.CancelFunc
	browser *zeroconf.Client // Added browser client
	server  *zeroconf.Client // Used as publisher
	ticker  *time.Ticker
	done    chan bool
	id      string
}

var (
	discoverService *DiscoverService
	once            sync.Once
	discoveryType   = zeroconf.NewType("_transok._tcp")
)

func GetDiscoverService() *DiscoverService {
	once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		discoverService = &DiscoverService{
			id:     uuid.New().String(),
			ctx:    ctx,
			cancel: cancel,
		}
	})
	return discoverService
}

/* Start starts listening for mDNS broadcasts */
func (s *DiscoverService) Start() error {
	if s.browser != nil {
		return nil // Already started
	}

	logger.Info("Starting discovery...")

	browser := zeroconf.New().Browse(func(e zeroconf.Event) {
		//		if e.Op != zeroconf.OpAdded {
		//			return
		//		}

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
		//if data.Sender == s.id {
		//	return
		//}

		logger.Info("Service discovered", zap.Any("data", data))
		mdns.GetDispatcher().Dispatch(data)
	}, discoveryType)

	client, err := browser.Open()
	if err != nil {
		logger.Error("Error discovering devices", zap.Error(err))
		return err
	}
	s.browser = client

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

	// If already broadcasting, close old one to update payload
	if s.server != nil {
		s.server.Close()
		s.server = nil
	}

	service := zeroconf.NewService(
		discoveryType,
		fmt.Sprintf("TransokService_%s", s.id),
		uint16(port),
	)
	service.Text = txtRecords

	publisher := zeroconf.New().Publish(service)
	_, err = publisher.Open()
	if err != nil {
		logger.Error("Error broadcasting service", zap.Error(err))
		return err
	}

	go func() { //Stop broadcast when done
		<-s.ctx.Done()
		publisher.Close()
	}()

	s.server = publisher

	logger.Info("Service started broadcasting",
		zap.String("id", s.id),
		zap.Int("port", port),
		zap.Any("payload", discoverPayload),
	)

	return nil
}

// Stop method
func (s *DiscoverService) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	if s.server != nil {
		s.server.Close()
		s.server = nil
	}
	if s.browser != nil {
		s.browser.Close()
		s.browser = nil
	}
}
