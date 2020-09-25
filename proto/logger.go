package proto

import (
	"encoding/json"
	"os"
)

type Logger struct {
	File os.File
}

func (l *Logger) LogPacketOut(id PacketID, packet EncodableAsPacket) error {
	row := map[string]interface{}{
		"id":   id,
		"sent": packet,
	}

	rowData, err := json.Marshal(row)
	if err != nil {
		return err
	}
	_, err = l.File.Write(rowData)
	return err
}
