package message

import "io"

type DataRow struct {
	ColumnValues []Value
}

func (m *DataRow) Reader() io.Reader {
	length := 2
	for _, v := range m.ColumnValues {
		length += v.Length()
	}

	b := NewBase(length)
	b.WriteUint16(uint16(len(m.ColumnValues)))
	for _, v := range m.ColumnValues {
		b.WriteValue(v)
	}
	return b.SetType('D').Reader()
}

func ReadDataRow(raw []byte) *DataRow {
	b := NewBaseFromBytes(raw)
	row := &DataRow{}
	row.ColumnValues = make([]Value, b.ReadUint16())
	for i := 0; i < len(row.ColumnValues); i++ {
		row.ColumnValues[i] = ReadValue(b)
	}
	return row
}
