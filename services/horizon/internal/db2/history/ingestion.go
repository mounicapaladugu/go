package history

import sq "github.com/Masterminds/squirrel"

// TruncateExpingestStateTables clears out ingestion state tables.
// Ingestion state tables are horizon database tables populated by
// the experimental ingestion system using history archive snapshots.
// Any horizon database tables which cannot be populated using
// history archive snapshots will not be truncated.
func (q *Q) TruncateExpingestStateTables() error {
	return q.TruncateTables([]string{
		"accounts",
		"accounts_data",
		"accounts_signers",
		"exp_asset_stats",
		"offers",
		"trust_lines",
	})
}

// ExpIngestRemovalSummary describes how many rows in the experimental ingestion
// history tables have been deleted by RemoveExpIngestHistory()
type ExpIngestRemovalSummary struct {
	LedgersRemoved      int64
	TransactionsRemoved int64
}

// RemoveExpIngestHistory removes all rows in the experimental ingestion
// history tables which have a ledger sequence higher than `newerThanSequence`
func (q *Q) RemoveExpIngestHistory(newerThanSequence uint32) (ExpIngestRemovalSummary, error) {
	summary := ExpIngestRemovalSummary{}

	result, err := q.Exec(
		sq.Delete("exp_history_ledgers").
			Where("sequence > ?", newerThanSequence),
	)
	if err != nil {
		return summary, err
	}

	summary.LedgersRemoved, err = result.RowsAffected()
	if err != nil {
		return summary, err
	}

	result, err = q.Exec(
		sq.Delete("exp_history_transactions").
			Where("ledger_sequence > ?", newerThanSequence),
	)
	if err != nil {
		return summary, err
	}

	summary.TransactionsRemoved, err = result.RowsAffected()
	return summary, err
}
