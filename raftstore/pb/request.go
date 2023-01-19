package pb

import (
	"encoding/binary"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

const (
	reqTypeCreateBucekt uint32 = 1
	reqTypeDeleteBucket uint32 = 2
	reqTypeTxn          uint32 = 3
	reqTypeCompact      uint32 = 4
	reqTypeSnapshot     uint32 = 5

	// 8 for uint64 ID, 4 for payload type and the last 4byte for payload size
	headerSize = 8 + 4 + 4
)

var (
	ErrUnknownRequestType    = errors.New("unknown request type")
	ErrDataTooShort          = errors.New("buf is too short to unmarshal")
	ErrUnexpectedRequestSize = errors.New("unexpected request size")
)

type request interface {
	proto.Marshaler
	Size() int
	MarshalToSizedBuffer(data []byte) (int, error)
}

type InternalRequest struct {
	ID uint64

	CreateBucket *CreateBucket
	DeleteBucket *DeleteBucket
	Txn          *Txn
	Compact      *Compact
	Snapshot     *Snapshot
}

func (r *InternalRequest) Marshal() ([]byte, error) {
	if r.CreateBucket != nil {
		return marshal(r.ID, reqTypeCreateBucekt, r.CreateBucket)
	} else if r.DeleteBucket != nil {
		return marshal(r.ID, reqTypeDeleteBucket, r.DeleteBucket)
	} else if r.Txn != nil {
		return marshal(r.ID, reqTypeTxn, r.Txn)
	} else if r.Compact != nil {
		return marshal(r.ID, reqTypeCompact, r.Compact)
	} else if r.Snapshot != nil {
		return marshal(r.ID, reqTypeSnapshot, r.Snapshot)
	} else {
		return nil, ErrUnknownRequestType
	}
}

func marshal(id uint64, typ uint32, m request) ([]byte, error) {
	ps := m.Size()

	buf := make([]byte, headerSize+ps)
	binary.BigEndian.PutUint64(buf, id)
	binary.BigEndian.PutUint32(buf[8:], typ)
	binary.BigEndian.PutUint32(buf[12:], uint32(ps))

	_, err := m.MarshalToSizedBuffer(buf[headerSize:])
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (r *InternalRequest) Unmarshal(data []byte) error {
	if len(data) <= headerSize {
		return ErrDataTooShort
	}

	r.ID = binary.BigEndian.Uint64(data)
	typ := binary.BigEndian.Uint32(data[8:])
	ps := binary.BigEndian.Uint32(data[12:])

	if ps+headerSize != uint32(len(data)) {
		return ErrUnexpectedRequestSize
	}

	switch typ {
	case reqTypeCreateBucekt:
		r.CreateBucket = &CreateBucket{}
		return r.CreateBucket.Unmarshal(data[headerSize:])
	case reqTypeDeleteBucket:
		r.DeleteBucket = &DeleteBucket{}
		return r.DeleteBucket.Unmarshal(data[headerSize:])
	case reqTypeTxn:
		r.Txn = &Txn{}
		return r.Txn.Unmarshal(data[headerSize:])
	case reqTypeCompact:
		r.Compact = &Compact{}
		return r.Compact.Unmarshal(data[headerSize:])
	case reqTypeSnapshot:
		r.Snapshot = &Snapshot{}
		return r.Snapshot.Unmarshal(data[headerSize:])
	default:
		return ErrUnknownRequestType
	}
}
