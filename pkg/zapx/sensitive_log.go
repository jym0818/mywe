package zapx

import (
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type MaskingEncoder struct {
	zapcore.Encoder
}

func NewMaskingEncoder(encoder zapcore.Encoder) zapcore.Encoder {
	return &MaskingEncoder{
		Encoder: encoder,
	}
}

func (e *MaskingEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// 在编码前处理字段
	maskedFields := make([]zapcore.Field, len(fields))
	for i, field := range fields {
		if field.Key == "phone" && field.Type == zapcore.StringType {
			// 复制字段并修改值
			maskedField := field
			phone := field.String
			if len(phone) >= 11 {
				maskedField.String = phone[:3] + "****" + phone[7:]
			}
			maskedFields[i] = maskedField
		} else {
			maskedFields[i] = field
		}
	}
	return e.Encoder.EncodeEntry(entry, maskedFields)
}
