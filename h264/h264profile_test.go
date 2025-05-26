package h264

import (
	"strings"
	"testing"
)

func TestParseProfileLevelId(t *testing.T) {
	tests := []struct {
		input           string
		expectedLevel   Level
		expectedProfile Profile
		shouldBeNil     bool
	}{
		{"42e01f", L3_1, ConstrainedBaseline, false},
		{"42e00b", L1_1, ConstrainedBaseline, false},
		{"42f00b", L1_b, ConstrainedBaseline, false},
		{"42C02A", L4_2, ConstrainedBaseline, false},
		{"640c34", L5_2, ConstrainedHigh, false},
		{"4de01f", L3_1, ConstrainedBaseline, false},
		{"58f01f", L3_1, ConstrainedBaseline, false},
		{"42a01f", L3_1, Baseline, false},
		{"58A01F", L3_1, Baseline, false},
		{"4D401f", L3_1, Main, false},
		{"64001f", L3_1, High, false},
		{"640c1f", L3_1, ConstrainedHigh, false},
		{"", 0, 0, true},
		{" 42e01f", 0, 0, true},
		{"4242e01f", 0, 0, true},
		{"e01f", 0, 0, true},
		{"gggggg", 0, 0, true},
		{"42e000", 0, 0, true},
		{"42e00f", 0, 0, true},
		{"42e0ff", 0, 0, true},
		{"42e11f", 0, 0, true},
		{"58601f", 0, 0, true},
		{"64e01f", 0, 0, true},
	}

	for _, test := range tests {
		result := ParseProfileLevelId(test.input)
		if test.shouldBeNil {
			if result != nil {
				t.Errorf("ParseProfileLevelId(%q) should return nil", test.input)
			}
			continue
		}
		if result == nil {
			t.Errorf("ParseProfileLevelId(%q) should not return nil", test.input)
			continue
		}
		if result.Level != test.expectedLevel {
			t.Errorf("ParseProfileLevelId(%q) level = %v, want %v", test.input, result.Level, test.expectedLevel)
		}
		if result.Profile != test.expectedProfile {
			t.Errorf("ParseProfileLevelId(%q) profile = %v, want %v", test.input, result.Profile, test.expectedProfile)
		}
	}
}

func TestProfileLevelIdToString(t *testing.T) {
	tests := []struct {
		input         ProfileLevelId
		expected      string
		shouldBeEmpty bool
	}{
		{ProfileLevelId{ConstrainedBaseline, L3_1}, "42e01f", false},
		{ProfileLevelId{Baseline, L1}, "42000a", false},
		{ProfileLevelId{Main, L3_1}, "4d001f", false},
		{ProfileLevelId{ConstrainedHigh, L4_2}, "640c2a", false},
		{ProfileLevelId{High, L4_2}, "64002a", false},
		{ProfileLevelId{ConstrainedBaseline, L1_b}, "42f00b", false},
		{ProfileLevelId{Baseline, L1_b}, "42100b", false},
		{ProfileLevelId{Main, L1_b}, "4d100b", false},
		{ProfileLevelId{High, L1_b}, "", true},
		{ProfileLevelId{ConstrainedHigh, L1_b}, "", true},
		{ProfileLevelId{255, L3_1}, "", true},
	}

	for _, test := range tests {
		result := ProfileLevelIdToString(test.input)
		if test.shouldBeEmpty {
			if result != "" {
				t.Errorf("ProfileLevelIdToString(%v) should return empty string", test.input)
			}
			continue
		}
		if result != test.expected {
			t.Errorf("ProfileLevelIdToString(%v) = %q, want %q", test.input, result, test.expected)
		}
	}

	// Round trip tests
	roundTripTests := []string{"42e01f", "42E01F", "4d100b", "4D100B", "640c2a", "640C2A"}
	for _, input := range roundTripTests {
		parsed := ParseProfileLevelId(input)
		if parsed == nil {
			t.Errorf("ParseProfileLevelId(%q) returned nil", input)
			continue
		}
		result := ProfileLevelIdToString(*parsed)
		if result != strings.ToLower(input) {
			t.Errorf("Round trip ParseProfileLevelId(%q) -> ProfileLevelIdToString = %q, want %q", input, result, strings.ToLower(input))
		}
	}
}

func TestParseSdpProfileLevelId(t *testing.T) {
	tests := []struct {
		params          map[string]interface{}
		expectedProfile Profile
		expectedLevel   Level
		shouldBeNil     bool
	}{
		{nil, ConstrainedBaseline, L3_1, false},
		{map[string]interface{}{"profile-level-id": "640c2a"}, ConstrainedHigh, L4_2, false},
		{map[string]interface{}{"profile-level-id": "foobar"}, 0, 0, true},
	}

	for _, test := range tests {
		result := ParseSdpProfileLevelId(test.params)
		if test.shouldBeNil {
			if result != nil {
				t.Errorf("ParseSdpProfileLevelId(%v) should return nil", test.params)
			}
			continue
		}
		if result == nil {
			t.Errorf("ParseSdpProfileLevelId(%v) should not return nil", test.params)
			continue
		}
		if result.Profile != test.expectedProfile {
			t.Errorf("ParseSdpProfileLevelId(%v) profile = %v, want %v", test.params, result.Profile, test.expectedProfile)
		}
		if result.Level != test.expectedLevel {
			t.Errorf("ParseSdpProfileLevelId(%v) level = %v, want %v", test.params, result.Level, test.expectedLevel)
		}
	}
}

func TestIsSameProfile(t *testing.T) {
	tests := []struct {
		params1  map[string]interface{}
		params2  map[string]interface{}
		expected bool
	}{
		{map[string]interface{}{"foo": "foo"}, map[string]interface{}{"bar": "bar"}, true},
		{map[string]interface{}{"profile-level-id": "42e01f"}, map[string]interface{}{"profile-level-id": "42C02A"}, true},
		{map[string]interface{}{"profile-level-id": "42a01f"}, map[string]interface{}{"profile-level-id": "58A01F"}, true},
		{map[string]interface{}{"profile-level-id": "42e01f"}, nil, true},
		{map[string]interface{}{"profile-level-id": "42a01f"}, map[string]interface{}{"profile-level-id": "640c1f"}, false},
		{map[string]interface{}{"profile-level-id": "42000a"}, map[string]interface{}{"profile-level-id": "64002a"}, false},
	}

	for _, test := range tests {
		result := IsSameProfile(test.params1, test.params2)
		if result != test.expected {
			t.Errorf("IsSameProfile(%v, %v) = %v, want %v", test.params1, test.params2, result, test.expected)
		}
	}
}

func TestIsSameProfileAndLevel(t *testing.T) {
	tests := []struct {
		params1  map[string]interface{}
		params2  map[string]interface{}
		expected bool
	}{
		{map[string]interface{}{"foo": "foo"}, map[string]interface{}{"bar": "bar"}, true},
		{map[string]interface{}{"profile-level-id": "42e01f"}, map[string]interface{}{"profile-level-id": "42f01f"}, true},
		{map[string]interface{}{"profile-level-id": "42a01f"}, map[string]interface{}{"profile-level-id": "58A01F"}, true},
		{map[string]interface{}{"profile-level-id": "42e01f"}, nil, true},
		{nil, map[string]interface{}{"profile-level-id": "4d001f"}, false},
		{map[string]interface{}{"profile-level-id": "42a01f"}, map[string]interface{}{"profile-level-id": "640c1f"}, false},
		{map[string]interface{}{"profile-level-id": "42000a"}, map[string]interface{}{"profile-level-id": "64002a"}, false},
	}

	for _, test := range tests {
		result := IsSameProfileAndLevel(test.params1, test.params2)
		if result != test.expected {
			t.Errorf("IsSameProfileAndLevel(%v, %v) = %v, want %v", test.params1, test.params2, result, test.expected)
		}
	}
}

func TestGenerateProfileLevelIdStringForAnswer(t *testing.T) {
	tests := []struct {
		localParams  map[string]interface{}
		remoteParams map[string]interface{}
		expected     string
		shouldError  bool
	}{
		{nil, nil, "", false},
		{map[string]interface{}{"profile-level-id": "42e015"}, map[string]interface{}{"profile-level-id": "42e01f"}, "42e015", false},
		{map[string]interface{}{"profile-level-id": "42e01f"}, map[string]interface{}{"profile-level-id": "42e015"}, "42e015", false},
		{map[string]interface{}{"profile-level-id": "42e01f", "level-asymmetry-allowed": "1"}, map[string]interface{}{"profile-level-id": "42e015", "level-asymmetry-allowed": "1"}, "42e01f", false},
	}

	for _, test := range tests {
		result, err := GenerateProfileLevelIdStringForAnswer(test.localParams, test.remoteParams)
		if test.shouldError {
			if err == nil {
				t.Errorf("GenerateProfileLevelIdStringForAnswer(%v, %v) should return error", test.localParams, test.remoteParams)
			}
			continue
		}
		if err != nil {
			t.Errorf("GenerateProfileLevelIdStringForAnswer(%v, %v) returned error: %v", test.localParams, test.remoteParams, err)
		}
		if result != test.expected {
			t.Errorf("GenerateProfileLevelIdStringForAnswer(%v, %v) = %q, want %q", test.localParams, test.remoteParams, result, test.expected)
		}
	}
}

func TestSupportedLevel(t *testing.T) {
	tests := []struct {
		pixelCount int
		fps        int
		expected   *Level
	}{
		{640 * 480, 25, levelPtr(L2_1)},
		{1280 * 720, 30, levelPtr(L3_1)},
		{1920 * 1280, 60, levelPtr(L4_2)},
		{0, 0, nil},
		{1280 * 720, 5, nil},
		{183 * 137, 30, nil},
	}

	for _, test := range tests {
		result := SupportedLevel(test.pixelCount, test.fps)
		if test.expected == nil {
			if result != nil {
				t.Errorf("SupportedLevel(%d, %d) should return nil", test.pixelCount, test.fps)
			}
			continue
		}
		if result == nil || *result != *test.expected {
			t.Errorf("SupportedLevel(%d, %d) = %v, want %v", test.pixelCount, test.fps, result, *test.expected)
		}
	}
}

// Helper function to create Level pointer
func levelPtr(l Level) *Level {
	return &l
}
