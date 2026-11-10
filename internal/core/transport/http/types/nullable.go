package core_http_types

import (
	"bytes"
	"encoding/json"

	"github.com/FITLPI/TODOGOLANG/internal/core/domain"
)

// Nullable — HTTP-версия доменного Nullable[T].
// Сам факт вызова UnmarshalJSON означает, что поле было передано в теле
// запроса (Set = true). Значение "null" очищает Value (явный null),
// любое другое значение декодируется в Value.
type Nullable[T any] struct {
	domain.Nullable[T]
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	n.Set = true

	if bytes.Equal(data, []byte("null")) {
		n.Value = nil
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	n.Value = &value

	return nil
}

func (n Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{Set: n.Set, Value: n.Value}
}
