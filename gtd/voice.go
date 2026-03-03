package gtd

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"time"
)

type VoiceProcessor struct {
	dataDir string
}

func NewVoiceProcessor(dataDir string) *VoiceProcessor {
	voiceDir := filepath.Join(dataDir, "voice_notes")
	os.MkdirAll(voiceDir, 0755)

	return &VoiceProcessor{
		dataDir: voiceDir,
	}
}

func (vp *VoiceProcessor) SaveVoiceNote(audioData []byte) (string, error) {
	filename := filepath.Join(vp.dataDir, time.Now().Format("20060102_150405")+".m4a")

	if err := os.WriteFile(filename, audioData, 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func (vp *VoiceProcessor) ProcessVoiceInput(audioData []byte) (string, error) {
	return "Распознанный текст из голосового сообщения", nil
}

func (vp *VoiceProcessor) GetVoiceNoteContent(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (vp *VoiceProcessor) EncodeVoiceNoteForPlayback(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}
