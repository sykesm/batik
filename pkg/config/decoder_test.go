// Copyright (c) 2015-2019 Carlos Alexandro Becker
// The MIT License (MIT)

package config

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

type unmarshaler struct {
	time.Duration
}

// TextUnmarshaler implements encoding.TextUnmarshaler
func (d *unmarshaler) UnmarshalText(data []byte) (err error) {
	if len(data) != 0 {
		d.Duration, err = time.ParseDuration(string(data))
	} else {
		d.Duration = 0
	}
	return err
}

type Config struct {
	String     string    `env:"STRING"`
	StringPtr  *string   `env:"STRING"`
	Strings    []string  `env:"STRINGS"`
	StringPtrs []*string `env:"STRINGS"`

	Bool     bool    `env:"BOOL"`
	BoolPtr  *bool   `env:"BOOL"`
	Bools    []bool  `env:"BOOLS"`
	BoolPtrs []*bool `env:"BOOLS"`

	Int     int    `env:"INT"`
	IntPtr  *int   `env:"INT"`
	Ints    []int  `env:"INTS"`
	IntPtrs []*int `env:"INTS"`

	Int8     int8    `env:"INT8"`
	Int8Ptr  *int8   `env:"INT8"`
	Int8s    []int8  `env:"INT8S"`
	Int8Ptrs []*int8 `env:"INT8S"`

	Int16     int16    `env:"INT16"`
	Int16Ptr  *int16   `env:"INT16"`
	Int16s    []int16  `env:"INT16S"`
	Int16Ptrs []*int16 `env:"INT16S"`

	Int32     int32    `env:"INT32"`
	Int32Ptr  *int32   `env:"INT32"`
	Int32s    []int32  `env:"INT32S"`
	Int32Ptrs []*int32 `env:"INT32S"`

	Int64     int64    `env:"INT64"`
	Int64Ptr  *int64   `env:"INT64"`
	Int64s    []int64  `env:"INT64S"`
	Int64Ptrs []*int64 `env:"INT64S"`

	Uint     uint    `env:"UINT"`
	UintPtr  *uint   `env:"UINT"`
	Uints    []uint  `env:"UINTS"`
	UintPtrs []*uint `env:"UINTS"`

	Uint8     uint8    `env:"UINT8"`
	Uint8Ptr  *uint8   `env:"UINT8"`
	Uint8s    []uint8  `env:"UINT8S"`
	Uint8Ptrs []*uint8 `env:"UINT8S"`

	Uint16     uint16    `env:"UINT16"`
	Uint16Ptr  *uint16   `env:"UINT16"`
	Uint16s    []uint16  `env:"UINT16S"`
	Uint16Ptrs []*uint16 `env:"UINT16S"`

	Uint32     uint32    `env:"UINT32"`
	Uint32Ptr  *uint32   `env:"UINT32"`
	Uint32s    []uint32  `env:"UINT32S"`
	Uint32Ptrs []*uint32 `env:"UINT32S"`

	Uint64     uint64    `env:"UINT64"`
	Uint64Ptr  *uint64   `env:"UINT64"`
	Uint64s    []uint64  `env:"UINT64S"`
	Uint64Ptrs []*uint64 `env:"UINT64S"`

	Float32     float32    `env:"FLOAT32"`
	Float32Ptr  *float32   `env:"FLOAT32"`
	Float32s    []float32  `env:"FLOAT32S"`
	Float32Ptrs []*float32 `env:"FLOAT32S"`

	Float64     float64    `env:"FLOAT64"`
	Float64Ptr  *float64   `env:"FLOAT64"`
	Float64s    []float64  `env:"FLOAT64S"`
	Float64Ptrs []*float64 `env:"FLOAT64S"`

	Duration     time.Duration    `env:"DURATION"`
	Durations    []time.Duration  `env:"DURATIONS"`
	DurationPtr  *time.Duration   `env:"DURATION"`
	DurationPtrs []*time.Duration `env:"DURATIONS"`

	Unmarshaler     unmarshaler    `env:"UNMARSHALER"`
	UnmarshalerPtr  *unmarshaler   `env:"UNMARSHALER"`
	Unmarshalers    []unmarshaler  `env:"UNMARSHALERS"`
	UnmarshalerPtrs []*unmarshaler `env:"UNMARSHALERS"`

	URL     url.URL    `env:"URL"`
	URLPtr  *url.URL   `env:"URL"`
	URLs    []url.URL  `env:"URLS"`
	URLPtrs []*url.URL `env:"URLS"`

	StringWithDefault string `env:"DATABASE_URL" example:"postgres://localhost:5432/db"`

	NonDefined struct {
		String string `env:"NONDEFINED_STR"`
	}

	NotAnEnv   string
	unexported string `env:"FOO"`
}

type ParentStruct struct {
	InnerStruct *InnerStruct
	unexported  *InnerStruct
	Ignored     *http.Client
}

type InnerStruct struct {
	Inner  string `env:"innervar"`
	Number uint   `env:"innernum"`
}

type ForNestedStruct struct {
	NestedStruct
}
type NestedStruct struct {
	NestedVar string `env:"nestedvar"`
}

func TestParse(t *testing.T) {
	gt := NewGomegaWithT(t)

	var tos = func(v interface{}) string {
		return fmt.Sprintf("%v", v)
	}

	var toss = func(v ...interface{}) string {
		var ss = []string{}
		for _, s := range v {
			ss = append(ss, tos(s))
		}
		return strings.Join(ss, ",")
	}

	envMap := make(EnvMap)
	var str1 = "str1"
	var str2 = "str2"
	envMap["STRING"] = str1
	envMap["STRINGS"] = toss(str1, str2)

	var bool1 = true
	var bool2 = false
	envMap["BOOL"] = tos(bool1)
	envMap["BOOLS"] = toss(bool1, bool2)

	var int1 = -1
	var int2 = 2
	envMap["INT"] = tos(int1)
	envMap["INTS"] = toss(int1, int2)

	var int81 int8 = -2
	var int82 int8 = 5
	envMap["INT8"] = tos(int81)
	envMap["INT8S"] = toss(int81, int82)

	var int161 int16 = -24
	var int162 int16 = 15
	envMap["INT16"] = tos(int161)
	envMap["INT16S"] = toss(int161, int162)

	var int321 int32 = -14
	var int322 int32 = 154
	envMap["INT32"] = tos(int321)
	envMap["INT32S"] = toss(int321, int322)

	var int641 int64 = -12
	var int642 int64 = 150
	envMap["INT64"] = tos(int641)
	envMap["INT64S"] = toss(int641, int642)

	var uint1 uint = 1
	var uint2 uint = 2
	envMap["UINT"] = tos(uint1)
	envMap["UINTS"] = toss(uint1, uint2)

	var uint81 uint8 = 15
	var uint82 uint8 = 51
	envMap["UINT8"] = tos(uint81)
	envMap["UINT8S"] = toss(uint81, uint82)

	var uint161 uint16 = 532
	var uint162 uint16 = 123
	envMap["UINT16"] = tos(uint161)
	envMap["UINT16S"] = toss(uint161, uint162)

	var uint321 uint32 = 93
	var uint322 uint32 = 14
	envMap["UINT32"] = tos(uint321)
	envMap["UINT32S"] = toss(uint321, uint322)

	var uint641 uint64 = 5
	var uint642 uint64 = 43
	envMap["UINT64"] = tos(uint641)
	envMap["UINT64S"] = toss(uint641, uint642)

	var float321 float32 = 9.3
	var float322 float32 = 1.1
	envMap["FLOAT32"] = tos(float321)
	envMap["FLOAT32S"] = toss(float321, float322)

	var float641 = 1.53
	var float642 = 0.5
	envMap["FLOAT64"] = tos(float641)
	envMap["FLOAT64S"] = toss(float641, float642)

	var duration1 = time.Second
	var duration2 = time.Second * 4
	envMap["DURATION"] = tos(duration1)
	envMap["DURATIONS"] = toss(duration1, duration2)

	var unmarshaler1 = unmarshaler{time.Minute}
	var unmarshaler2 = unmarshaler{time.Millisecond * 1232}
	envMap["UNMARSHALER"] = tos(unmarshaler1.Duration)
	envMap["UNMARSHALERS"] = toss(unmarshaler1.Duration, unmarshaler2.Duration)

	var url1s = "https://goreleaser.com"
	var url2s = "https://caarlos0.dev"
	url1, err := url.Parse(url1s)
	gt.Expect(err).NotTo(HaveOccurred())
	url2, err := url.Parse(url2s)
	gt.Expect(err).NotTo(HaveOccurred())
	envMap["URL"] = tos(url1)
	envMap["URLS"] = toss(url1, url2)

	envMap["SEPSTRINGS"] = strings.Join([]string{str1, str2}, ":")

	nonDefinedStr := "nonDefinedStr"
	envMap["NONDEFINED_STR"] = nonDefinedStr

	var cfg = Config{}

	decoder := Decoder{
		lookuper:   envMap,
		parseTag:   "env",
		defaultTag: "example",
	}
	err = decoder.Parse(&cfg)
	gt.Expect(err).NotTo(HaveOccurred())

	gt.Expect(cfg).To(Equal(Config{
		String:     str1,
		StringPtr:  &str1,
		Strings:    []string{str1, str2},
		StringPtrs: []*string{&str1, &str2},

		Bool:     bool1,
		BoolPtr:  &bool1,
		Bools:    []bool{bool1, bool2},
		BoolPtrs: []*bool{&bool1, &bool2},

		Int:     int1,
		IntPtr:  &int1,
		Ints:    []int{int1, int2},
		IntPtrs: []*int{&int1, &int2},

		Int8:     int81,
		Int8Ptr:  &int81,
		Int8s:    []int8{int81, int82},
		Int8Ptrs: []*int8{&int81, &int82},

		Int16:     int161,
		Int16Ptr:  &int161,
		Int16s:    []int16{int161, int162},
		Int16Ptrs: []*int16{&int161, &int162},

		Int32:     int321,
		Int32Ptr:  &int321,
		Int32s:    []int32{int321, int322},
		Int32Ptrs: []*int32{&int321, &int322},

		Int64:     int641,
		Int64Ptr:  &int641,
		Int64s:    []int64{int641, int642},
		Int64Ptrs: []*int64{&int641, &int642},

		Uint:     uint1,
		UintPtr:  &uint1,
		Uints:    []uint{uint1, uint2},
		UintPtrs: []*uint{&uint1, &uint2},

		Uint8:     uint81,
		Uint8Ptr:  &uint81,
		Uint8s:    []uint8{uint81, uint82},
		Uint8Ptrs: []*uint8{&uint81, &uint82},

		Uint16:     uint161,
		Uint16Ptr:  &uint161,
		Uint16s:    []uint16{uint161, uint162},
		Uint16Ptrs: []*uint16{&uint161, &uint162},

		Uint32:     uint321,
		Uint32Ptr:  &uint321,
		Uint32s:    []uint32{uint321, uint322},
		Uint32Ptrs: []*uint32{&uint321, &uint322},

		Uint64:     uint641,
		Uint64Ptr:  &uint641,
		Uint64s:    []uint64{uint641, uint642},
		Uint64Ptrs: []*uint64{&uint641, &uint642},

		Float32:     float321,
		Float32Ptr:  &float321,
		Float32s:    []float32{float321, float322},
		Float32Ptrs: []*float32{&float321, &float322},

		Float64:     float641,
		Float64Ptr:  &float641,
		Float64s:    []float64{float641, float642},
		Float64Ptrs: []*float64{&float641, &float642},

		Duration:     duration1,
		DurationPtr:  &duration1,
		Durations:    []time.Duration{duration1, duration2},
		DurationPtrs: []*time.Duration{&duration1, &duration2},

		Unmarshaler:     unmarshaler1,
		UnmarshalerPtr:  &unmarshaler1,
		Unmarshalers:    []unmarshaler{unmarshaler1, unmarshaler2},
		UnmarshalerPtrs: []*unmarshaler{&unmarshaler1, &unmarshaler2},

		URL:     *url1,
		URLPtr:  url1,
		URLs:    []url.URL{*url1, *url2},
		URLPtrs: []*url.URL{url1, url2},

		StringWithDefault: "postgres://localhost:5432/db",
		NonDefined: struct {
			String string `env:"NONDEFINED_STR"`
		}{
			String: nonDefinedStr,
		},

		NotAnEnv:   "",
		unexported: "",
	}))

	// gt.Expect(cfg.URL.String()).To(Equal(url1))
	// gt.Expect(cfg.URLPtr.String()).To(Equal(url1))
	// gt.Expect(cfg.URLs[0].String()).To(Equal(url1))
	// gt.Expect(cfg.URLs[1].String()).To(Equal(url2))
	// gt.Expect(cfg.URLPtrs[0].String()).To(Equal(url1))
	// gt.Expect(cfg.URLPtrs[1].String()).To(Equal(url2))
}

func TestParsesEnvInner(t *testing.T) {
	gt := NewGomegaWithT(t)

	cfg := ParentStruct{
		InnerStruct: &InnerStruct{},
		unexported:  &InnerStruct{},
	}

	envMap := EnvMap{
		"innervar": "someinnervalue",
		"innernum": "8",
	}

	decoder := Decoder{
		lookuper:   envMap,
		parseTag:   "env",
		defaultTag: "example",
	}

	err := decoder.Parse(&cfg)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(cfg.InnerStruct.Inner).To(Equal("someinnervalue"))
	gt.Expect(cfg.InnerStruct.Number).To(Equal(uint(8)))
}

func TestParsesEnvInnerFails(t *testing.T) {
	gt := NewGomegaWithT(t)

	type config struct {
		Foo struct {
			Number int `env:"NUMBER"`
		}
	}
	var cfg config

	envMap := EnvMap{
		"NUMBER": "not-a-number",
	}

	decoder := Decoder{
		lookuper:   envMap,
		parseTag:   "env",
		defaultTag: "example",
	}

	err := decoder.Parse(&cfg)
	gt.Expect(err).To(MatchError("decode: parse error on field \"Number\" of type \"int\": strconv.ParseInt: parsing \"not-a-number\": invalid syntax"))
}

func TestParsesEnvInnerNil(t *testing.T) {
	gt := NewGomegaWithT(t)

	envMap := EnvMap{
		"innervar": "someinnervalue",
	}

	decoder := Decoder{
		lookuper:   envMap,
		parseTag:   "env",
		defaultTag: "example",
	}

	cfg := ParentStruct{}
	err := decoder.Parse(&cfg)
	gt.Expect(err).NotTo(HaveOccurred())
}

func TestParsesEnvInnerInvalid(t *testing.T) {
	gt := NewGomegaWithT(t)

	envMap := EnvMap{
		"innernum": "-547",
	}

	decoder := Decoder{
		lookuper:   envMap,
		parseTag:   "env",
		defaultTag: "example",
	}

	cfg := ParentStruct{
		InnerStruct: &InnerStruct{},
	}
	err := decoder.Parse(&cfg)
	gt.Expect(err).To(MatchError("decode: parse error on field \"Number\" of type \"uint\": strconv.ParseUint: parsing \"-547\": invalid syntax"))
}

func TestParsesEnvNested(t *testing.T) {
	gt := NewGomegaWithT(t)

	envMap := EnvMap{
		"nestedvar": "somenestedvalue",
	}

	decoder := Decoder{
		lookuper:   envMap,
		parseTag:   "env",
		defaultTag: "example",
	}

	var cfg ForNestedStruct
	err := decoder.Parse(&cfg)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(cfg.NestedVar).To(Equal("somenestedvalue"))
}

func TestEmptyVars(t *testing.T) {
	gt := NewGomegaWithT(t)

	var cfg Config

	decoder := Decoder{
		lookuper:   EnvMap{},
		parseTag:   "env",
		defaultTag: "example",
	}

	err := decoder.Parse(&cfg)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(cfg.String).To(BeEmpty())
	gt.Expect(cfg.Bool).To(BeFalse())
	gt.Expect(cfg.Int).To(BeZero())
	gt.Expect(cfg.Uint).To(BeZero())
	gt.Expect(cfg.Uint64).To(BeZero())
	gt.Expect(cfg.Int64).To(BeZero())
	gt.Expect(cfg.Strings).To(BeEmpty())
	gt.Expect(cfg.Ints).To(BeEmpty())
	gt.Expect(cfg.Bools).To(BeEmpty())
}

func TestPassAnInvalidPtr(t *testing.T) {
	gt := NewGomegaWithT(t)

	var thisShouldBreak int

	decoder := Decoder{
		lookuper: EnvMap{},
	}

	err := decoder.Parse(&thisShouldBreak)
	gt.Expect(err).To(MatchError("decode: expected a pointer to a Struct"))
}

func TestPassReference(t *testing.T) {
	gt := NewGomegaWithT(t)

	var cfg Config

	decoder := Decoder{
		lookuper: EnvMap{},
	}

	err := decoder.Parse(cfg)
	gt.Expect(err).To(MatchError("decode: expected a pointer to a Struct"))
}

func TestInvalidTypes(t *testing.T) {
	tests := []struct {
		testName    string
		envMap      EnvMap
		expectedErr string
	}{
		{
			testName:    "invalid int",
			envMap:      EnvMap{"INT": "should-be-an-int"},
			expectedErr: `decode: parse error on field "Int" of type "int": strconv.ParseInt: parsing "should-be-an-int": invalid syntax`,
		},
		{
			testName:    "invalid uint",
			envMap:      EnvMap{"UINT": "-44"},
			expectedErr: `decode: parse error on field "Uint" of type "uint": strconv.ParseUint: parsing "-44": invalid syntax`,
		},
		{
			testName:    "invalid float32",
			envMap:      EnvMap{"FLOAT32": "AAA"},
			expectedErr: `decode: parse error on field "Float32" of type "float32": strconv.ParseFloat: parsing "AAA": invalid syntax`,
		},
		{
			testName:    "invalid float64",
			envMap:      EnvMap{"FLOAT64": "AAA"},
			expectedErr: `decode: parse error on field "Float64" of type "float64": strconv.ParseFloat: parsing "AAA": invalid syntax`,
		},
		{
			testName:    "invalid uint64",
			envMap:      EnvMap{"UINT64": "AAA"},
			expectedErr: `decode: parse error on field "Uint64" of type "uint64": strconv.ParseUint: parsing "AAA": invalid syntax`,
		},
		{
			testName:    "invalid int64",
			envMap:      EnvMap{"INT64": "AAA"},
			expectedErr: `decode: parse error on field "Int64" of type "int64": strconv.ParseInt: parsing "AAA": invalid syntax`,
		},
		{
			testName:    "invalid int64 slice",
			envMap:      EnvMap{"INT64S": "A,2,3"},
			expectedErr: `decode: parse error on field "Int64s" of type "\[\]int64": strconv.ParseInt: parsing "A": invalid syntax`,
		},
		{
			testName:    "invalid uint64 slice",
			envMap:      EnvMap{"UINT64S": "A,2,3"},
			expectedErr: `decode: parse error on field "Uint64s" of type "\[\]uint64": strconv.ParseUint: parsing "A": invalid syntax`,
		},
		{
			testName:    "invalid float32 slice",
			envMap:      EnvMap{"FLOAT32S": "A,2.0,3.0"},
			expectedErr: `decode: parse error on field "Float32s" of type "\[\]float32": strconv.ParseFloat: parsing "A": invalid syntax`,
		},
		{
			testName:    "invalid float64 slice",
			envMap:      EnvMap{"FLOAT64S": "A,2.0,3.0"},
			expectedErr: `decode: parse error on field "Float64s" of type "\[\]float64": strconv.ParseFloat: parsing "A": invalid syntax`,
		},
		{
			testName:    "invalid bool",
			envMap:      EnvMap{"BOOL": "should-be-a-bool"},
			expectedErr: `decode: parse error on field "Bool" of type "bool": strconv.ParseBool: parsing "should-be-a-bool": invalid syntax`,
		},
		{
			testName:    "invalid bool slice",
			envMap:      EnvMap{"BOOLS": "t,f,TRUE,faaaalse"},
			expectedErr: `decode: parse error on field "Bools" of type "\[\]bool": strconv.ParseBool: parsing "faaaalse": invalid syntax`,
		},
		{
			testName:    "invalid duration",
			envMap:      EnvMap{"DURATION": "should-be-a-valid-duration"},
			expectedErr: `decode: parse error on field "Duration" of type "time.Duration": unable to parse duration: time: invalid duration "?should-be-a-valid-duration"?`,
		},
		{
			testName:    "invalid duration slice",
			envMap:      EnvMap{"DURATIONS": "1s,contains-an-invalid-duration,3s"},
			expectedErr: `decode: parse error on field "Durations" of type "\[\]time.Duration": unable to parse duration: time: invalid duration "?contains-an-invalid-duration"?`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			var cfg Config

			decoder := Decoder{
				lookuper:   tt.envMap,
				parseTag:   "env",
				defaultTag: "example",
			}

			err := decoder.Parse(&cfg)
			gt.Expect(err).To(MatchError(MatchRegexp(tt.expectedErr)))
		})
	}
}

func TestParseStructWithoutEnvTag(t *testing.T) {
	gt := NewGomegaWithT(t)

	cfg := Config{}

	decoder := Decoder{
		lookuper:   EnvMap{},
		defaultTag: "example",
	}

	err := decoder.Parse(&cfg)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(cfg.NotAnEnv).To(BeEmpty())
}

func TestParseStructWithInvalidFieldKind(t *testing.T) {
	gt := NewGomegaWithT(t)

	type config struct {
		WontWorkByte byte `env:"BLAH"`
	}
	var cfg config

	envMap := EnvMap{
		"BLAH": "a",
	}

	decoder := Decoder{
		lookuper:   envMap,
		parseTag:   "env",
		defaultTag: "example",
	}

	err := decoder.Parse(&cfg)
	gt.Expect(err).To(MatchError("decode: parse error on field \"WontWorkByte\" of type \"uint8\": strconv.ParseUint: parsing \"a\": invalid syntax"))
}

func TestUnsupportedSliceType(t *testing.T) {
	gt := NewGomegaWithT(t)

	type config struct {
		WontWork []map[int]int `env:"WONTWORK"`
	}
	var cfg config

	envMap := EnvMap{
		"WONTWORK": "1,2,3",
	}

	decoder := Decoder{
		lookuper:   envMap,
		parseTag:   "env",
		defaultTag: "example",
	}

	err := decoder.Parse(&cfg)
	gt.Expect(err).To(MatchError("decode: no parser found for field \"WontWork\" of type \"[]map[int]int\""))
}

func TestCustomParserBasicUnsupported(t *testing.T) {
	gt := NewGomegaWithT(t)

	type ConstT struct {
		A int
	}

	type config struct {
		Const ConstT `env:"CONST_"`
	}
	var cfg config

	envMap := EnvMap{
		"CONST_": "42",
	}

	decoder := Decoder{
		lookuper:   envMap,
		parseTag:   "env",
		defaultTag: "example",
	}

	err := decoder.Parse(&cfg)

	gt.Expect(cfg.Const).To(BeZero())
	gt.Expect(err).To(MatchError("decode: no parser found for field \"Const\" of type \"config.ConstT\""))
}

func TestUnsupportedStructType(t *testing.T) {
	gt := NewGomegaWithT(t)

	type config struct {
		Foo http.Client `env:"FOO"`
	}
	var cfg config

	envMap := EnvMap{
		"FOO": "foo",
	}

	decoder := Decoder{
		lookuper:   envMap,
		parseTag:   "env",
		defaultTag: "example",
	}

	err := decoder.Parse(&cfg)

	gt.Expect(err).To(MatchError("decode: no parser found for field \"Foo\" of type \"http.Client\""))
}

func TestEmptyOption(t *testing.T) {
	gt := NewGomegaWithT(t)

	type config struct {
		Var string `env:"VAR,"`
	}
	var cfg config

	envMap := EnvMap{
		"VAR": "",
	}
	decoder := Decoder{
		lookuper:   envMap,
		parseTag:   "env",
		defaultTag: "example",
	}

	err := decoder.Parse(&cfg)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(cfg.Var).To(Equal(""))
}

func TestErrorOptionNotRecognized(t *testing.T) {
	gt := NewGomegaWithT(t)

	type config struct {
		Var string `env:"VAR,not_supported!"`
	}
	var cfg config

	decoder := Decoder{
		lookuper:   EnvMap{},
		parseTag:   "env",
		defaultTag: "example",
	}

	err := decoder.Parse(&cfg)
	gt.Expect(err).To(MatchError("decode: tag option \"not_supported!\" not supported"))
}

func TestTextUnmarshalerError(t *testing.T) {
	gt := NewGomegaWithT(t)

	type config struct {
		Unmarshaler unmarshaler `env:"UNMARSHALER"`
	}
	var cfg config

	envMap := EnvMap{
		"UNMARSHALER": "invalid",
	}

	decoder := Decoder{
		lookuper:   envMap,
		parseTag:   "env",
		defaultTag: "example",
	}
	err := decoder.Parse(&cfg)
	gt.Expect(err).To(MatchError(MatchRegexp(`decode: parse error on field "Unmarshaler" of type "config.unmarshaler": time: invalid duration "?invalid"?`)))
}

func TestTextUnmarshalersError(t *testing.T) {
	gt := NewGomegaWithT(t)

	type config struct {
		Unmarshalers []unmarshaler `env:"UNMARSHALERS"`
	}
	var cfg config

	envMap := EnvMap{
		"UNMARSHALERS": "1s,invalid",
	}

	decoder := Decoder{
		lookuper:   envMap,
		parseTag:   "env",
		defaultTag: "example",
	}

	err := decoder.Parse(&cfg)
	gt.Expect(err).To(MatchError(MatchRegexp(`decode: parse error on field "Unmarshalers" of type "\[\]config.unmarshaler": time: invalid duration "?invalid"?`)))
}

func TestParseURL(t *testing.T) {
	gt := NewGomegaWithT(t)

	type config struct {
		ExampleURL url.URL `env:"EXAMPLE_URL" example:"https://google.com"`
	}
	var cfg config

	decoder := Decoder{
		lookuper:   EnvMap{},
		parseTag:   "env",
		defaultTag: "example",
	}
	err := decoder.Parse(&cfg)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(cfg.ExampleURL.String()).To(Equal("https://google.com"))
}

func TestParseURLFailure(t *testing.T) {
	gt := NewGomegaWithT(t)

	type config struct {
		ExampleURL url.URL `env:"EXAMPLE_URL_2"`
	}
	var cfg config

	envMap := EnvMap{
		"EXAMPLE_URL_2": "nope://s s/",
	}

	decoder := Decoder{
		lookuper:   envMap,
		parseTag:   "env",
		defaultTag: "example",
	}

	err := decoder.Parse(&cfg)
	gt.Expect(err).To(MatchError("decode: parse error on field \"ExampleURL\" of type \"url.URL\": unable to parse URL: parse \"nope://s s/\": invalid character \" \" in host name"))
}

func TestIgnoresUnexported(t *testing.T) {
	gt := NewGomegaWithT(t)

	type unexportedConfig struct {
		home  string `env:"HOME"`
		Home2 string `env:"HOME"`
	}
	var cfg unexportedConfig

	envMap := EnvMap{
		"HOME": "/tmp/fakehome",
	}

	decoder := Decoder{
		lookuper:   envMap,
		parseTag:   "env",
		defaultTag: "example",
	}

	err := decoder.Parse(&cfg)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(cfg.home).To(BeEmpty())
	gt.Expect(cfg.Home2).To(Equal("/tmp/fakehome"))
}

type LogLevel int8

func (l *LogLevel) UnmarshalText(text []byte) error {
	txt := string(text)
	switch txt {
	case "debug":
		*l = DebugLevel
	case "info":
		*l = InfoLevel
	default:
		return fmt.Errorf("unknown level: %q", txt)
	}

	return nil
}

const (
	DebugLevel LogLevel = iota - 1
	InfoLevel
)

func TestPrecedenceUnmarshalText(t *testing.T) {
	gt := NewGomegaWithT(t)

	type config struct {
		LogLevel  LogLevel   `env:"LOG_LEVEL"`
		LogLevels []LogLevel `env:"LOG_LEVELS"`
	}
	var cfg config

	envMap := EnvMap{
		"LOG_LEVEL":  "debug",
		"LOG_LEVELS": "debug,info",
	}

	decoder := Decoder{
		lookuper:   envMap,
		parseTag:   "env",
		defaultTag: "example",
	}

	err := decoder.Parse(&cfg)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(cfg.LogLevel).To(Equal(DebugLevel))
	gt.Expect(cfg.LogLevels).To(Equal([]LogLevel{DebugLevel, InfoLevel}))
}
