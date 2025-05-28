package h264

import (
	"fmt"
	"log"
	"strconv"
)

// Profile represents H.264 profile
type Profile int

const (
	ConstrainedBaseline Profile = 1
	Baseline            Profile = 2
	Main                Profile = 3
	ConstrainedHigh     Profile = 4
	High                Profile = 5
	PredictiveHigh444   Profile = 6
)

// Level represents H.264 level
type Level int

const (
	L1_b Level = 0
	L1   Level = 10
	L1_1 Level = 11
	L1_2 Level = 12
	L1_3 Level = 13
	L2   Level = 20
	L2_1 Level = 21
	L2_2 Level = 22
	L3   Level = 30
	L3_1 Level = 31
	L3_2 Level = 32
	L4   Level = 40
	L4_1 Level = 41
	L4_2 Level = 42
	L5   Level = 50
	L5_1 Level = 51
	L5_2 Level = 52
)

// ProfileLevelId represents a parsed H.264 profile-level-id
type ProfileLevelId struct {
	Profile Profile
	Level   Level
}

// DefaultProfileLevelId is the default profile level ID
var DefaultProfileLevelId = ProfileLevelId{
	Profile: ConstrainedBaseline,
	Level:   L3_1,
}

// BitPattern matches bit patterns
type BitPattern struct {
	Mask        uint8
	MaskedValue uint8
}

// ProfilePattern maps profile_idc/profile_iop to Profile
type ProfilePattern struct {
	ProfileIdc uint8
	ProfileIop BitPattern
	Profile    Profile
}

// LevelConstraint defines level limits
type LevelConstraint struct {
	MaxMacroblocksPerSecond int
	MaxMacroblockFrameSize  int
	Level                   Level
}

var ProfilePatterns = []ProfilePattern{
	{0x42, BitPattern{^byteMaskString("x", "x1xx0000"), byteMaskString("1", "x1xx0000")}, ConstrainedBaseline},
	{0x4d, BitPattern{^byteMaskString("x", "1xxx0000"), byteMaskString("1", "1xxx0000")}, ConstrainedBaseline},
	{0x58, BitPattern{^byteMaskString("x", "11xx0000"), byteMaskString("1", "11xx0000")}, ConstrainedBaseline},
	{0x42, BitPattern{^byteMaskString("x", "x0xx0000"), byteMaskString("1", "x0xx0000")}, Baseline},
	{0x58, BitPattern{^byteMaskString("x", "10xx0000"), byteMaskString("1", "10xx0000")}, Baseline},
	{0x4d, BitPattern{^byteMaskString("x", "0x0x0000"), byteMaskString("1", "0x0x0000")}, Main},
	{0x64, BitPattern{^byteMaskString("x", "00000000"), byteMaskString("1", "00000000")}, High},
	{0x64, BitPattern{^byteMaskString("x", "00001100"), byteMaskString("1", "00001100")}, ConstrainedHigh},
	{0xf4, BitPattern{^byteMaskString("x", "00000000"), byteMaskString("1", "00000000")}, PredictiveHigh444},
}

var LevelConstraints = []LevelConstraint{
	{1485, 99, L1},
	{1485, 99, L1_b},
	{3000, 396, L1_1},
	{6000, 396, L1_2},
	{11880, 396, L1_3},
	{11880, 396, L2},
	{19800, 792, L2_1},
	{20250, 1620, L2_2},
	{40500, 1620, L3},
	{108000, 3600, L3_1},
	{216000, 5120, L3_2},
	{245760, 8192, L4},
	{245760, 8192, L4_1},
	{522240, 8704, L4_2},
	{589824, 22080, L5},
	{983040, 36864, L5_1},
	{2073600, 36864, L5_2},
}

// ParseProfileLevelId parses a profile-level-id string
func ParseProfileLevelId(str string) *ProfileLevelId {
	const ConstraintSet3Flag = 0x10

	if len(str) != 6 {
		return nil
	}

	profileLevelIdNumeric, err := strconv.ParseInt(str, 16, 32)
	if err != nil || profileLevelIdNumeric == 0 {
		return nil
	}

	levelIdc := Level(profileLevelIdNumeric & 0xff)
	profileIop := uint8((profileLevelIdNumeric >> 8) & 0xff)
	profileIdc := uint8((profileLevelIdNumeric >> 16) & 0xff)

	var level Level
	switch levelIdc {
	case L1_1:
		if profileIop&ConstraintSet3Flag != 0 {
			level = L1_b
		} else {
			level = L1_1
		}
	case L1, L1_2, L1_3, L2, L2_1, L2_2, L3, L3_1, L3_2, L4, L4_1, L4_2, L5, L5_1, L5_2:
		level = levelIdc
	default:
		log.Printf("parseProfileLevelId() | unrecognized level_idc [str:%s, level_idc:%d]", str, levelIdc)
		return nil
	}

	for _, pattern := range ProfilePatterns {
		if profileIdc == pattern.ProfileIdc && pattern.ProfileIop.IsMatch(profileIop) {
			log.Printf("parseProfileLevelId() | result [str:%s, profile:%d, level:%d]", str, pattern.Profile, level)
			return &ProfileLevelId{Profile: pattern.Profile, Level: level}
		}
	}

	log.Printf("parseProfileLevelId() | unrecognized profile_idc/profile_iop [str:%s, profile_idc:%d, profile_iop:%d]", str, profileIdc, profileIop)
	return nil
}

// ProfileLevelIdToString converts ProfileLevelId to string
func ProfileLevelIdToString(profileLevelId ProfileLevelId) string {
	if profileLevelId.Level == L1_b {
		switch profileLevelId.Profile {
		case ConstrainedBaseline:
			return "42f00b"
		case Baseline:
			return "42100b"
		case Main:
			return "4d100b"
		default:
			log.Printf("profileLevelIdToString() | Level 1_b not allowed for profile %d", profileLevelId.Profile)
			return ""
		}
	}

	var profileIdcIopString string
	switch profileLevelId.Profile {
	case ConstrainedBaseline:
		profileIdcIopString = "42e0"
	case Baseline:
		profileIdcIopString = "4200"
	case Main:
		profileIdcIopString = "4d00"
	case ConstrainedHigh:
		profileIdcIopString = "640c"
	case High:
		profileIdcIopString = "6400"
	case PredictiveHigh444:
		profileIdcIopString = "f400"
	default:
		log.Printf("profileLevelIdToString() | unrecognized profile %d", profileLevelId.Profile)
		return ""
	}

	levelStr := fmt.Sprintf("%x", profileLevelId.Level)
	if len(levelStr) == 1 {
		levelStr = "0" + levelStr
	}

	return profileIdcIopString + levelStr
}

// ProfileToString converts Profile to string
func ProfileToString(profile Profile) string {
	switch profile {
	case ConstrainedBaseline:
		return "ConstrainedBaseline"
	case Baseline:
		return "Baseline"
	case Main:
		return "Main"
	case ConstrainedHigh:
		return "ConstrainedHigh"
	case High:
		return "High"
	case PredictiveHigh444:
		return "PredictiveHigh444"
	default:
		log.Printf("profileToString() | unrecognized profile %d", profile)
		return ""
	}
}

// LevelToString converts Level to string
func LevelToString(level Level) string {
	switch level {
	case L1_b:
		return "1b"
	case L1:
		return "1"
	case L1_1:
		return "1.1"
	case L1_2:
		return "1.2"
	case L1_3:
		return "1.3"
	case L2:
		return "2"
	case L2_1:
		return "2.1"
	case L2_2:
		return "2.2"
	case L3:
		return "3"
	case L3_1:
		return "3.1"
	case L3_2:
		return "3.2"
	case L4:
		return "4"
	case L4_1:
		return "4.1"
	case L4_2:
		return "4.2"
	case L5:
		return "5"
	case L5_1:
		return "5.1"
	case L5_2:
		return "5.2"
	default:
		log.Printf("levelToString() | unrecognized level %d", level)
		return ""
	}
}

// ParseSdpProfileLevelId parses profile-level-id from SDP parameters
func ParseSdpProfileLevelId(params map[string]interface{}) *ProfileLevelId {
	if profileLevelId, ok := params["profile-level-id"].(string); ok {
		return ParseProfileLevelId(profileLevelId)
	}
	return &DefaultProfileLevelId
}

// IsSameProfile checks if two sets of parameters have the same H.264 profile
func IsSameProfile(params1, params2 map[string]interface{}) bool {
	profileLevelId1 := ParseSdpProfileLevelId(params1)
	profileLevelId2 := ParseSdpProfileLevelId(params2)
	return profileLevelId1 != nil && profileLevelId2 != nil && profileLevelId1.Profile == profileLevelId2.Profile
}

// IsSameProfileAndLevel checks if two sets of parameters have the same H.264 profile and level
func IsSameProfileAndLevel(params1, params2 map[string]interface{}) bool {
	profileLevelId1 := ParseSdpProfileLevelId(params1)
	profileLevelId2 := ParseSdpProfileLevelId(params2)
	return profileLevelId1 != nil && profileLevelId2 != nil &&
		profileLevelId1.Profile == profileLevelId2.Profile &&
		profileLevelId1.Level == profileLevelId2.Level
}

// GenerateProfileLevelIdStringForAnswer generates profile-level-id for SDP answer
func GenerateProfileLevelIdStringForAnswer(localSupportedParams, remoteOfferedParams map[string]interface{}) (string, error) {
	if localSupportedParams["profile-level-id"] == nil && remoteOfferedParams["profile-level-id"] == nil {
		log.Println("generateProfileLevelIdStringForAnswer() | profile-level-id missing in local and remote params")
		return "", nil
	}

	localProfileLevelId := ParseSdpProfileLevelId(localSupportedParams)
	remoteProfileLevelId := ParseSdpProfileLevelId(remoteOfferedParams)

	if localProfileLevelId == nil {
		return "", fmt.Errorf("invalid local_profile_level_id")
	}
	if remoteProfileLevelId == nil {
		return "", fmt.Errorf("invalid remote_profile_level_id")
	}
	if localProfileLevelId.Profile != remoteProfileLevelId.Profile {
		return "", fmt.Errorf("H264 Profile mismatch")
	}

	levelAsymmetryAllowed := IsLevelAsymmetryAllowed(localSupportedParams) && IsLevelAsymmetryAllowed(remoteOfferedParams)
	localLevel := localProfileLevelId.Level
	remoteLevel := remoteProfileLevelId.Level
	minLevel := MinLevel(localLevel, remoteLevel)

	answerLevel := minLevel
	if levelAsymmetryAllowed {
		answerLevel = localLevel
	}

	log.Printf("generateProfileLevelIdStringForAnswer() | result [profile:%d, level:%d]", localProfileLevelId.Profile, answerLevel)
	return ProfileLevelIdToString(ProfileLevelId{localProfileLevelId.Profile, answerLevel}), nil
}

// SupportedLevel determines the highest supported H.264 level
func SupportedLevel(maxFramePixelCount, maxFps int) *Level {
	const PixelsPerMacroblock = 16 * 16

	for i := len(LevelConstraints) - 1; i >= 0; i-- {
		constraint := LevelConstraints[i]
		if constraint.MaxMacroblockFrameSize*PixelsPerMacroblock <= maxFramePixelCount &&
			constraint.MaxMacroblocksPerSecond <= maxFps*constraint.MaxMacroblockFrameSize {
			log.Printf("supportedLevel() | result [max_frame_pixel_count:%d, max_fps:%d, level:%d]", maxFramePixelCount, maxFps, constraint.Level)
			return &constraint.Level
		}
	}

	log.Printf("supportedLevel() | no level supported [max_frame_pixel_count:%d, max_fps:%d]", maxFramePixelCount, maxFps)
	return nil
}

// byteMaskString converts a string pattern to a byte
func byteMaskString(c string, str string) uint8 {
	if len(str) != 8 {
		return 0
	}
	var result uint8
	for i, char := range str {
		if string(char) == c {
			result |= 1 << (7 - i)
		}
	}
	return result
}

// IsMatch checks if a value matches the bit pattern
func (bp BitPattern) IsMatch(value uint8) bool {
	return bp.MaskedValue == (value & bp.Mask)
}

// IsLessLevel compares H.264 levels
func IsLessLevel(a, b Level) bool {
	if a == L1_b {
		return b != L1 && b != L1_b
	}
	if b == L1_b {
		return a != L1
	}
	return a < b
}

// MinLevel returns the minimum H.264 level
func MinLevel(a, b Level) Level {
	if IsLessLevel(a, b) {
		return a
	}
	return b
}

// IsLevelAsymmetryAllowed checks if level asymmetry is allowed
func IsLevelAsymmetryAllowed(params map[string]interface{}) bool {
	if val, ok := params["level-asymmetry-allowed"]; ok {
		switch v := val.(type) {
		case bool:
			return v
		case int:
			return v == 1
		case string:
			return v == "1"
		}
	}
	return false
}
