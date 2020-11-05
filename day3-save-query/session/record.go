package session

import "lorm/day3-save-query/caluse"

//
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(values).refTable
		s.caluse.Set(caluse.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}
	s.caluse.Set(caluse.VALUES, recordValues)
	sql, vars := s.caluse.Build(caluse.INSERT, caluse.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
