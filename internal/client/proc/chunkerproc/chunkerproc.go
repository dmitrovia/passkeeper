package chunkerproc

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/dmitrovia/passkeeper/internal/client/models/procattrs/chunkerpa"
	"github.com/dmitrovia/passkeeper/internal/general/models/chunckmeta"
)

type ChunkerProc struct {
	attr *chunkerpa.ChunkerProcAttr
}

func NewProc(attr *chunkerpa.ChunkerProcAttr) *ChunkerProc {
	return &ChunkerProc{
		attr: attr,
	}
}

func (cp *ChunkerProc) RunProcess() error {
	if cp.attr == nil {
		cp.attr = &chunkerpa.ChunkerProcAttr{}
	}

	var waitGroup sync.WaitGroup

	var mutex sync.Mutex

	var chunks []chunckmeta.ChunkMeta

	file, err := os.Open(cp.attr.FilePath)
	if err != nil {
		return fmt.Errorf("RP->os.Open: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("RP->Stat: %w", err)
	}

	numChunks := int(fileInfo.Size() /
		int64(cp.attr.ChunkSize))
	if fileInfo.Size()%int64(cp.attr.ChunkSize) != 0 {
		numChunks++
	}

	chunkChan := make(chan chunckmeta.ChunkMeta, numChunks)
	errChan := make(chan error, numChunks)
	indexChan := make(chan int, numChunks)

	for i := range numChunks {
		indexChan <- i
	}

	close(indexChan)

	for range cp.attr.CountWorkers {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			for index := range indexChan {
				offset := int64(index) * int64(cp.attr.ChunkSize)
				buffer := make([]byte,
					cp.attr.ChunkSize)

				_, err := file.Seek(offset, 0)
				if err != nil {
					errChan <- err

					return
				}

				bytesRead, err := file.Read(buffer)
				if err != nil && errors.Is(err, io.EOF) {
					errChan <- err

					return
				}

				if bytesRead > 0 {
					hash := md5.Sum(buffer[:bytesRead])
					hashString := hex.EncodeToString(hash[:])

					chunkFileName := fmt.Sprintf("%s.chunk.%d",
						cp.attr.FilePath, index)

					chunkFile, err := os.Create(chunkFileName)
					if err != nil {
						errChan <- err

						return
					}

					_, err = chunkFile.Write(buffer[:bytesRead])
					if err != nil {
						errChan <- err

						return
					}

					chunk := chunckmeta.ChunkMeta{
						FileName: &chunkFileName,
						Hash:     &hashString,
						Index:    &index,
					}

					mutex.Lock()
					chunks = append(chunks, chunk)
					mutex.Unlock()

					chunkFile.Close()

					chunkChan <- chunk
				}
			}
		}()
	}

	go func() {
		waitGroup.Wait()
		close(chunkChan)
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
