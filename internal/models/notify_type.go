package models

import "errors"

const (
	NotifyTypeMailMask int8 = 1 << iota
	NotifyTypeSlackMask
	NotifyTypeWebhookMask
	NotifyTypeServerChan3Mask
)

const NotifyTypeAllMask = NotifyTypeMailMask | NotifyTypeSlackMask | NotifyTypeWebhookMask | NotifyTypeServerChan3Mask

func NormalizeNotifyTypeMask(v int8) (int8, error) {
	if v < 0 {
		return 0, errors.New("invalid notify type")
	}
	if v >= 0 && v <= 3 {
		return int8(1 << v), nil
	}
	if v&^NotifyTypeAllMask != 0 {
		return 0, errors.New("invalid notify type")
	}
	return v, nil
}

func NotifyTypeMaskToTypes(mask int8) []int8 {
	types := make([]int8, 0, 4)
	if mask&NotifyTypeMailMask != 0 {
		types = append(types, 0)
	}
	if mask&NotifyTypeSlackMask != 0 {
		types = append(types, 1)
	}
	if mask&NotifyTypeWebhookMask != 0 {
		types = append(types, 2)
	}
	if mask&NotifyTypeServerChan3Mask != 0 {
		types = append(types, 3)
	}
	return types
}
