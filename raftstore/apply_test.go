package raftstore

import "testing"

func TestApplyBatch(t *testing.T) {

}

/*
var (
		appliedIndex, appliedTerm uint64
		shouldStop                bool
	)

	s.logger.Debug("applying entries",
		zap.Int("count", len(entries)))

	for i := range entries {
		entry := entries[i]
		s.logger.Debug("applying entry",
			zap.Uint64("index", entry.Index),
			zap.Uint64("term", entry.Term),
			zap.Stringer("type", entry.Type))

		switch entry.Type {
		case raftpb.EntryNormal:
			s.applyEntryNormal(&entry)
			s.setAppliedIndex(entry.Index)
			s.setTerm(entry.Term)

		case raftpb.EntryConfChange:
			var cc raftpb.ConfChange
			pbutil.MustUnmarshal(&cc, entry.Data)
			removedSelf, err := s.applyConfChange(cc, confState)
			s.setAppliedIndex(entry.Index)
			s.setTerm(entry.Term)

			shouldStop = shouldStop || removedSelf
			s.wait.Trigger(cc.ID, &confChangeResponse{s.cluster.Members(), err})

		default:
			s.logger.Fatal("unknown entry type, it must be EntryNormal or EntryConfChange",
				zap.Stringer("type", entry.Type))
		}

		appliedIndex, appliedTerm = entry.Index, entry.Term
	}

	return appliedTerm, appliedIndex, shouldStop
*/
