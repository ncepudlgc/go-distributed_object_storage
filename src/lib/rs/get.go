package rs

import (
	"Distributed_Object_Storage/src/lib/objectstream"
	"fmt"
	"io"
)

type RSGetStream struct {
	*decoder
}

func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
	// 检查locateInfo和dataServers的长度总和是否等于ALL_SHARDS
	if len(locateInfo)+len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}

	// 创建一个长度为ALL_SHARDS的io.Reader切片
	readers := make([]io.Reader, ALL_SHARDS)
	for i := 0; i < ALL_SHARDS; i++ {
		// 获取当前索引i对应的服务器地址
		server := locateInfo[i]
		if server == "" {
			// 如果服务器地址为空，则将dataServers的第一个地址赋值给locateInfo[i]
			locateInfo[i] = dataServers[0]
			// 移除dataServers的第一个地址
			dataServers = dataServers[1:]
			continue
		}
		// 调用objectstream.NewGetStream函数创建一个GetStream对象，并将其赋值给reader
		reader, e := objectstream.NewGetStream(server, fmt.Sprintf("%s.%d", hash, i))
		if e == nil {
			// 如果创建成功，则将reader添加到readers切片中
			readers[i] = reader
		}
	}

	// 创建一个长度为ALL_SHARDS的io.Writer切片
	writers := make([]io.Writer, ALL_SHARDS)
	// 计算每个分片的大小
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	var e error
	for i := range readers {
		// 如果readers[i]为空
		if readers[i] == nil {
			// 调用objectstream.NewTempPutStream函数创建一个TempPutStream对象，并将其赋值给writer
			writers[i], e = objectstream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", hash, i), perShard)
			if e != nil {
				// 如果创建失败，则返回错误
				return nil, e
			}
		}
	}

	// 创建一个NewDecoder对象，并将readers、writers和size作为参数传入
	dec := NewDecoder(readers, writers, size)
	// 返回一个包含dec的RSGetStream对象和一个nil错误
	return &RSGetStream{dec}, nil
}

func (s *RSGetStream) Close() {
	for i := range s.writers {
		if s.writers[i] != nil {
			s.writers[i].(*objectstream.TempPutStream).Commit(true)
		}
	}
}

func (s *RSGetStream) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekCurrent {
		panic("only support SeekCurrent")
	}
	if offset < 0 {
		panic("only support forward seek")
	}
	for offset != 0 {
		length := int64(BLOCK_SIZE)
		if offset < length {
			length = offset
		}
		buf := make([]byte, length)
		io.ReadFull(s, buf)
		offset -= length
	}
	return offset, nil
}
