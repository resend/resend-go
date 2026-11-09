package resend

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

func concludeWebhookType(resource string) (reflect.Type, error) {
	if payloadType, ok := WebhookPayloadDataTypes[resource]; ok {
		return payloadType, nil
	}

	return reflect.TypeOf(nil), fmt.Errorf("`%s` not implemented type", resource)
}

func (s *WebhookPayload) UnmarshalJSON(b []byte) error {
	var (
		p     map[string]interface{}
		pInfo struct {
			CreatedAt time.Time `json:"created_at"`
			Type      string    `json:"type"`
		}
	)
	err := json.Unmarshal(b, &pInfo)
	if err != nil {
		return err
	}
	s.CreatedAt = pInfo.CreatedAt
	s.Type = pInfo.Type

	err = json.Unmarshal(b, &p)
	if err != nil {
		return err
	}
	if data, ok := p["data"].(map[string]interface{}); ok {
		s.Data, err = decodeWebhookPayload(s.Type, data)
		if err != nil {
			return fmt.Errorf("decode webhook payload data: %w", err)
		}
	} else {
		return fmt.Errorf("data must be a JSON object")
	}
	return nil
}

func decodeWebhookPayload(webhookType string, data map[string]interface{}) (WebhookPayloadData, error) {
	payloadType, err := concludeWebhookType(webhookType)
	if err != nil {
		return nil, fmt.Errorf("conclude webhook payload type: %w", err)
	}

	payload := reflect.New(payloadType).Interface()

	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal webhook payload data: %w", err)
	}

	if err := json.Unmarshal(b, payload); err != nil {
		return nil, fmt.Errorf("unmarshal webhook payload data: %w", err)
	}

	return payload.(WebhookPayloadData), nil
}

// marshalJSONWithReplacement marshals the normal generated struct and then
// replaces one or more fields with pre-marshaled JSON.
//
// This avoids the field-shadowing/collision issues that can occur when
// embedding an alias and declaring another field with the same JSON name.
func marshalJSONWithReplacement(
	v interface{},
	replacements map[string]json.RawMessage,
) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var obj map[string]json.RawMessage

	if err := json.Unmarshal(b, &obj); err != nil {
		return nil, err
	}

	for name, value := range replacements {
		obj[name] = value
	}

	return json.Marshal(obj)
}

// encodeWebhookPayloadData marshals a CollectionRuleConditionObject implementation
func encodeWebhookPayloadData(
	v WebhookPayloadData,
) (json.RawMessage, error) {
	if isNilInterface(v) {
		return json.RawMessage("null"), nil
	}

	return encodeObject(
		v,
		WebhookPayloadDataTypes,
		"WebhookPayloadData",
	)
}

// MarshalJSON encodes WebhookPayload
//
// This is the inverse of:
//
//	func (s *WebhookPayload) UnmarshalJSON(...)
//
// in unmarshal.go.
func (s WebhookPayload) MarshalJSON() ([]byte, error) {
	type alias WebhookPayload

	payloadObject, err := encodeWebhookPayloadData(s.Data)
	if err != nil {
		return nil, err
	}

	return marshalJSONWithReplacement(
		alias(s),
		map[string]json.RawMessage{
			"data": payloadObject,
		},
	)
}

func encodeObject(
	v interface{},
	types map[string]reflect.Type,
	kind string,
) (json.RawMessage, error) {
	if v == nil {
		return json.RawMessage("null"), nil
	}

	t := reflect.TypeOf(v)

	// Interface implementations returned by decodeByTypename are pointers,
	// while the registered reflect.Type values are the underlying structs.
	if t.Kind() == reflect.Ptr {
		if reflect.ValueOf(v).IsNil() {
			return json.RawMessage("null"), nil
		}
		t = t.Elem()
	}

	var typename string

	for name, registeredType := range types {
		if registeredType == t {
			typename = name
			break
		}
	}

	if typename == "" {
		return nil, fmt.Errorf(
			"encode %s: unsupported type %q",
			kind,
			t.String(),
		)
	}

	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf(
			"encode %s %s: %w",
			kind,
			typename,
			err,
		)
	}

	// The concrete object should always be a JSON object, but using
	// RawMessage here preserves the generated field encoding exactly.
	var obj map[string]json.RawMessage

	if err := json.Unmarshal(b, &obj); err != nil {
		return nil, fmt.Errorf(
			"encode %s %s: %w",
			kind,
			typename,
			err,
		)
	}

	b, err = json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf(
			"encode %s %s: %w",
			kind,
			typename,
			err,
		)
	}

	return b, nil
}

// isNilInterface handles the slightly annoying case where an interface contains
// a typed nil pointer:
//
//	var m Media = (*MediaImage)(nil)
//
// In that case m != nil, but it should still marshal as JSON null.
func isNilInterface(v interface{}) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Ptr,
		reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}
