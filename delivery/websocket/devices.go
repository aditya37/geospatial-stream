package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	geofence_logger "github.com/aditya37/geofence-service/util"
	"github.com/aditya37/geospatial-stream/repository/channel"
	"github.com/aditya37/geospatial-stream/usecase/device"
	gws "github.com/gorilla/websocket"
)

type (
	DeviceWebsocketDeliver struct {
		tickTrigger             *time.Ticker
		upgrader                gws.Upgrader
		tickerDuration          time.Duration
		deviceManagerCase       *device.DeviceManageUscase
		deviceLogsChannelStream *channel.ChannelStreamDeviceLogs
	}
)

func NewDeviceWebsocketDeliver(
	intervalTrigger int,
	deviceManagerCase *device.DeviceManageUscase,
	deviceLogsChannelStream *channel.ChannelStreamDeviceLogs,
) *DeviceWebsocketDeliver {
	// trigger ticker...
	duration := time.Duration(intervalTrigger * int(time.Second))
	triggerTick := time.NewTicker(duration)

	// http socket upgrader...
	socketUpgrade := gws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	return &DeviceWebsocketDeliver{
		tickTrigger:             triggerTick,
		upgrader:                socketUpgrade,
		tickerDuration:          duration,
		deviceManagerCase:       deviceManagerCase,
		deviceLogsChannelStream: deviceLogsChannelStream,
	}
}

//StreamDeviceLog...
func (dd *DeviceWebsocketDeliver) StreamDeviceLogs(w http.ResponseWriter, r *http.Request) {
	wsconn, err := dd.upgrader.Upgrade(w, r, nil)
	if err != nil {
		geofence_logger.Logger().Error(err)
		return
	}

	// reader...
	go dd.readSocketDeviceLogs(wsconn)
	go dd.writeSocketDeviceLog(wsconn)

}

// writeSocketDeviceLog....
func (dd *DeviceWebsocketDeliver) writeSocketDeviceLog(conn *gws.Conn) {
	// recreate ticker...
	dd.tickTrigger.Reset(dd.tickerDuration)

	// forever run and write...
	for {
		select {
		case data := <-dd.deviceLogsChannelStream.OutcomeData:
			j, _ := json.Marshal(data)
			if err := conn.WriteMessage(gws.TextMessage, j); err != nil {
				geofence_logger.Logger().Error(err)
				return
			}
		case <-dd.tickTrigger.C:
			if _, err := dd.deviceManagerCase.TriggerDeviceLogs(context.Background()); err != nil {
				geofence_logger.Logger().Error(err)
				return
			}
		}
	}
}

// readSocketDeviceLogs....
func (dd *DeviceWebsocketDeliver) readSocketDeviceLogs(conn *gws.Conn) {
	for {
		if _, _, err := conn.NextReader(); err != nil {
			geofence_logger.Logger().Error(err)
			conn.Close()
			break
		}
	}
}
