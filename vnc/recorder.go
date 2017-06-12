package vnc

import (
	"os"
)

type Recorder struct {
	RBSFileName string
	fileHandle  *os.File
	logger      Logger
}

func NewRecorder(saveFilePath string, logger Logger) *Recorder {
	rec := Recorder{RBSFileName: saveFilePath}
	var err error
	rec.fileHandle, err = os.OpenFile(saveFilePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		logger.Errorf("unable to open file: %s, error: %v", saveFilePath, err)
		return nil
	}
	return &rec
}

func (r *Recorder) Write(data []byte) error {
	_, err := r.fileHandle.Write(data)
	return err
}

// func (r *Recorder) WriteUInt8(data uint8) error {
// 	buf := make([]byte, 1)
// 	buf[0] = byte(data) // cast int8 to byte
// 	return r.Write(buf)
// }

func (r *Recorder) Close() {
	r.fileHandle.Close()
}
