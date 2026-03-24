package audio

import (
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/gen2brain/malgo"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type Recorder struct {
	context      *malgo.AllocatedContext
	device       *malgo.Device
	samples      []float32
	mu           sync.Mutex
	isRecording  bool
	amplitude    float32
	sampleRate   uint32
	recoveryPath string
	recoveryFile *os.File
}

func NewRecorder(sampleRate uint32, configDir string) (*Recorder, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, err
	}

	return &Recorder{
		context:      ctx,
		samples:      make([]float32, 0),
		sampleRate:   sampleRate,
		recoveryPath: filepath.Join(configDir, "recovery.bin"),
	}, nil
}

func (r *Recorder) Start() error {
	// Open recovery file for appending
	f, err := os.OpenFile(r.recoveryPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err == nil {
		r.recoveryFile = f
	}

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatF32
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = r.sampleRate

	onRec := func(pSampleOut, pSampleIn []byte, frameCount uint32) {
		if !r.isRecording {
			return
		}

		samples := make([]float32, frameCount)
		for i := uint32(0); i < frameCount; i++ {
			bits := uint32(pSampleIn[i*4]) | uint32(pSampleIn[i*4+1])<<8 | uint32(pSampleIn[i*4+2])<<16 | uint32(pSampleIn[i*4+3])<<24
			samples[i] = math.Float32frombits(bits)
		}

		r.mu.Lock()
		r.samples = append(r.samples, samples...)
		
		// Write to recovery file
		if r.recoveryFile != nil {
			_ = binary.Write(r.recoveryFile, binary.LittleEndian, samples)
		}

		var sum float32
		for _, s := range samples {
			if s < 0 {
				sum += -s
			} else {
				sum += s
			}
		}
		if frameCount > 0 {
			r.amplitude = sum / float32(frameCount)
		}
		r.mu.Unlock()
	}

	deviceCallbacks := malgo.DeviceCallbacks{
		Data: onRec,
	}

	device, err := malgo.InitDevice(r.context.Context, deviceConfig, deviceCallbacks)
	if err != nil {
		return err
	}

	err = device.Start()
	if err != nil {
		return err
	}

	r.device = device
	r.isRecording = true
	return nil
}

func (r *Recorder) Stop() {
	r.isRecording = false
	if r.recoveryFile != nil {
		r.recoveryFile.Close()
		r.recoveryFile = nil
	}
	if r.device != nil {
		r.device.Stop()
		r.device.Uninit()
	}
}

func (r *Recorder) Close() {
	if r.context != nil {
		r.context.Uninit()
		r.context.Free()
	}
}

func (r *Recorder) GetAmplitude() float32 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.amplitude
}

func (r *Recorder) Save(path string, format string) (string, error) {
	r.mu.Lock()
	samples := make([]float32, len(r.samples))
	copy(samples, r.samples)
	r.mu.Unlock()

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return "", err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	e := wav.NewEncoder(f, int(r.sampleRate), 16, 1, 1)
	
	intSamples := make([]int, len(samples))
	for i, s := range samples {
		intSamples[i] = int(s * 32767)
	}

	buf := &audio.IntBuffer{
		Data: intSamples,
		Format: &audio.Format{
			SampleRate: int(r.sampleRate),
			NumChannels: 1,
		},
	}

	if err := e.Write(buf); err != nil {
		return "", err
	}

	if err := e.Close(); err != nil {
		return "", err
	}

	// Successful save, delete recovery file
	_ = os.Remove(r.recoveryPath)

	return path, nil
}

func (r *Recorder) ClearRecovery() {
	_ = os.Remove(r.recoveryPath)
}

func LoadRecovery(path string) ([]float32, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	count := info.Size() / 4
	samples := make([]float32, count)
	err = binary.Read(f, binary.LittleEndian, &samples)
	if err != nil {
		return nil, err
	}

	return samples, nil
}

func SaveSamples(path string, samples []float32, sampleRate uint32) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	e := wav.NewEncoder(f, int(sampleRate), 16, 1, 1)
	
	intSamples := make([]int, len(samples))
	for i, s := range samples {
		intSamples[i] = int(s * 32767)
	}

	buf := &audio.IntBuffer{
		Data: intSamples,
		Format: &audio.Format{
			SampleRate: int(sampleRate),
			NumChannels: 1,
		},
	}

	if err := e.Write(buf); err != nil {
		return err
	}

	return e.Close()
}
