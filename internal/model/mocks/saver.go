package mocks

import "github.com/bassosimone/oonidsl/internal/model"

// Saver saves a measurement on some persistent storage.
type Saver struct {
	MockSaveMeasurement func(m *model.Measurement) error
}

func (s *Saver) SaveMeasurement(m *model.Measurement) error {
	return s.MockSaveMeasurement(m)
}
