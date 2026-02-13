package realtime

const (
	MessageSync      = 0
	MessageAwareness = 1
)

const (
	SyncStep1  = 0
	SyncStep2  = 1
	SyncUpdate = 2
)

func ParseYjsMessage(msg []byte) (msgType byte, payload []byte) {
	if len(msg) == 0 {
		return 0, nil
	}
	return msg[0], msg[1:]
}

func ParseYjsSyncMessage(payload []byte) (syncType byte, data []byte) {
	if len(payload) == 0 {
		return 0, nil
	}
	return payload[0], payload[1:]
}

func CreateSyncStep2Message(update []byte) []byte {
	result := make([]byte, 2+len(update))
	result[0] = MessageSync
	result[1] = SyncStep2
	copy(result[2:], update)
	return result
}

func CreateSyncUpdateMessage(update []byte) []byte {
	result := make([]byte, 2+len(update))
	result[0] = MessageSync
	result[1] = SyncUpdate
	copy(result[2:], update)
	return result
}
