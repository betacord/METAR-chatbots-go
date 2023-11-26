package main

import (
	"testing"
	"time"
)

func TestDecodeMETAR(t *testing.T) {
	testCases := []struct {
		input          string
		expectedOutput METARReport
		expectedError  bool
	}{
		// Przykłady testowe...

		{"METAR EDDF 041600Z 25010KT 9999 SCT025 10/M02 Q1020 NOSIG", METARReport{
			AirportCode:          "EDDF",
			ObservationTime:      time.Date(2023, time.April, 4, 16, 0, 0, 0, time.UTC),
			WindDirection:        250,
			WindSpeed:            10,
			Visibility:           9999,
			AtmosphericPhenomena: nil,
			Clouds:               []string{"SCT"},
			CloudBase:            []int{2500},
			Temperature:          10,
			DewPoint:             -2,
			Pressure:             1020,
			SignificantChanges:   false,
		}, false},

		// Dodaj więcej przypadków testowych...

	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			decodedMetar, err := decodeMETAR(tc.input)

			if tc.expectedError && err == nil {
				t.Errorf("Expected an error, but got none.")
			}

			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}

			if !tc.expectedError && !decodedMetar.compare(&tc.expectedOutput) {
				t.Errorf("Expected %v, but got %v", tc.expectedOutput, *decodedMetar)
			}
		})
	}
}
