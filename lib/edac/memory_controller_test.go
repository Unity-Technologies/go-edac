package edac

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	testMemInfo = MemoryInfo{
		Name:                "mc0",
		SinceReset:          time.Minute,
		Type:                "test-type",
		Size:                10,
		Uncorrectable:       1,
		UncorrectableNoInfo: 2,
		Correctable:         3,
		CorrectableNoInfo:   4,
		ScrubRate:           5,
		MaxLocation:         "test-location",
	}

	testDimmRank = DimmRank{
		Name:          "dimm0",
		Size:          10,
		DeviceType:    "test-device-type",
		Mode:          "test-mode",
		Label:         "test-label",
		Location:      "test-location",
		MemoryType:    "test-mem-type",
		Correctable:   6,
		Uncorrectable: 7,
	}
)

func mockSysfs(t *testing.T) func() {
	var err error
	root, err = ioutil.TempDir("", "edac")
	require.NoError(t, err)
	var ok bool
	clean := func() {
		require.NoError(t, os.RemoveAll(root))
	}

	defer func() {
		if !ok {
			clean()
		}
	}()

	// Create a single memory controller
	mc0 := filepath.Join(root, "mc0")
	require.NoError(t, os.Mkdir(mc0, 0755))

	// Create mc0 readable entries.
	v := reflect.ValueOf(testMemInfo)
	tp := v.Type()
	for j := tp.NumField() - 1; j >= 0; j-- {
		tf := tp.Field(j)
		if tf.Name == "Name" {
			continue
		}
		file := tf.Tag.Get("file")
		require.NotEmpty(t, file, tf.Name)
		fv := v.Field(j)

		var data []byte
		perm := os.FileMode(0444)
		switch tf.Name {
		case "SinceReset":
			// Convert to seconds.
			data = []byte(fmt.Sprintf("%d\n", int64(fv.Interface().(time.Duration).Seconds())))
		case "ScrubRate":
			// ScrubRate is read / write.
			perm = 0666
			fallthrough
		default:
			data = []byte(fmt.Sprintf("%v\n", fv.Interface()))
		}

		require.NoError(t, ioutil.WriteFile(filepath.Join(mc0, file), data, perm))
	}

	// Create mc0 write only entries (readable so we can validate the written value).
	require.NoError(t, ioutil.WriteFile(filepath.Join(mc0, "reset_counters"), nil, 0666))

	// Create mc0 dimm0 entries.
	dimm0 := filepath.Join(mc0, testDimmRank.Name)
	require.NoError(t, os.Mkdir(dimm0, 0755))
	v = reflect.ValueOf(testDimmRank)
	tp = v.Type()
	for j := tp.NumField() - 1; j >= 0; j-- {
		tf := tp.Field(j)
		if tf.Name == "Name" {
			continue
		}
		file := tf.Tag.Get("file")
		require.NotEmpty(t, file, tf.Name)
		fv := v.Field(j)

		perm := os.FileMode(0444)
		data := []byte(fmt.Sprintf("%v\n", fv.Interface()))
		require.NoError(t, ioutil.WriteFile(filepath.Join(dimm0, file), data, perm))
	}

	ok = true
	return clean
}

var (
	subTests = []struct {
		name string
		f    func(t *testing.T, mc *MemoryController)
	}{
		{
			name: "info",
			f:    testInfo,
		},
		{
			name: "reset-counters",
			f:    testResetCounters,
		},
		{
			name: "set-scrub-rate",
			f:    testSetScrubRate,
		},
		{
			name: "dimm-ranks",
			f:    testDimmRanks,
		},
	}
)

func TestMemoryControllers(t *testing.T) {
	t.Run("not-supported", func(t *testing.T) {
		root = "/notexists/edac"
		_, err := MemoryControllers()
		require.EqualError(t, err, ErrMemoryNotSupported.Error())
	})

	clean := mockSysfs(t)
	defer clean()

	mcs, err := MemoryControllers()
	require.NoError(t, err)
	require.Len(t, mcs, 1)

	mc := mcs[0]
	require.Equal(t, "mc0", mc.Name)

	for _, tc := range subTests {
		t.Run(tc.name, func(t *testing.T) {
			tc.f(t, &mc)
		})
	}
}

func testInfo(t *testing.T, mc *MemoryController) {
	info, err := mc.Info()
	require.NoError(t, err)
	require.Equal(t, testMemInfo, *info)
	require.Contains(t, info.String(), mc.Name)
}

func testResetCounters(t *testing.T, mc *MemoryController) {
	require.NoError(t, mc.ResetCounters())
	d, err := ioutil.ReadFile(filepath.Join(root, mc.Name, "reset_counters"))
	require.NoError(t, err)
	require.Equal(t, []byte("1"), d)
}

func testSetScrubRate(t *testing.T, mc *MemoryController) {
	require.NoError(t, mc.SetScrubRate(100))
	d, err := ioutil.ReadFile(filepath.Join(root, mc.Name, "sdram_scrub_rate"))
	require.NoError(t, err)
	require.Equal(t, []byte("100"), d)
}

func testDimmRanks(t *testing.T, mc *MemoryController) {
	drs, err := mc.DimmRanks()
	require.NoError(t, err)
	require.Len(t, drs, 1)
	require.Equal(t, testDimmRank, drs[0])
}
