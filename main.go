package main

import "log"
import "fmt"
import "os"
import "bytes"
import "encoding/binary"

type OpenTypeHeader struct {
  numTables uint16
  searchRange uint16
  entrySelector uint16
  rangeShift uint16
}

type TableRecord struct {
  tag string
  checkSum uint32
  offset uint32
  length uint32
}

func readOpenTypeHeader(file *os.File, header *OpenTypeHeader) (error) {
  _, _ = getString(file, 4)

  numTables, numErr := getUshort(file)
  if numErr != nil {
    return numErr
  }

  searchRange, searchErr := getUshort(file)
  if searchErr != nil {
    return searchErr
  }

  entrySelector, entryErr := getUshort(file)
  if entryErr != nil {
    return entryErr
  }

  rangeShift, rangeErr := getUshort(file)
  if rangeErr != nil {
    return rangeErr
  }

  header.numTables = numTables
  header.searchRange = searchRange
  header.entrySelector = entrySelector
  header.rangeShift = rangeShift

  return nil
}

func readTable(file *os.File, table *TableRecord) error {
  tag, err := getString(file, 4)

  if err != nil {
    return err
  }

  checkSum, err := getUlong(file)

  if err != nil {
    return err
  }

  offset, err := getUlong(file)

  if err != nil {
    return err
  }

  length, err := getUlong(file)

  if err != nil {
    return err
  }

  table.offset = offset
  table.tag = tag
  table.checkSum = checkSum
  table.length = length

  return nil
}

func readTables(file *os.File, header *OpenTypeHeader, tables []TableRecord) error {
  var numTables int
  numTables = int(header.numTables)
  for i := 0; i < numTables; i++ {
    readTable(file, &tables[i])
  }
  return nil
}

func getString(file *os.File, byteCount int) (string, error) {
  bytes := make([]byte, byteCount)
  _, err := file.Read(bytes)

  if err != nil {
    return "", err
  }

  return string(bytes), nil
}

func getUshort(file *os.File) (uint16, error) {
  readBytes := make([]byte, 2)
  _, fileErr := file.Read(readBytes)

  if fileErr != nil {
    return 0, fileErr
  }

  var num uint16

  buf := bytes.NewReader(readBytes)
  binErr := binary.Read(buf, binary.BigEndian, &num)

  if binErr != nil {
    return 0, binErr
  }

  return num, nil
}

func getUlong(file *os.File) (uint32, error) {
  readBytes := make([]byte, 4)
  _, fileErr := file.Read(readBytes)

  if fileErr != nil {
    return 0, fileErr
  }

  var num uint32

  buf := bytes.NewReader(readBytes)
  binErr := binary.Read(buf, binary.BigEndian, &num)

  if binErr != nil {
    return 0, binErr
  }

  return num, nil
}

func main() {
  file, err := os.Open("font.otf")

  if err != nil {
    log.Fatal(err)
  }

  header := OpenTypeHeader{}

  readOpenTypeHeader(file, &header)

  tables := make([]TableRecord, header.numTables)

  tableErr := readTables(file, &header, tables)

  if tableErr != nil {
    log.Fatal(tableErr)
  }

  for i := 0; i < len(tables); i++ {
    tr := tables[i]
    b := make([]byte, tr.length)
    _, _ = file.ReadAt(b, int64(tr.offset))
  }

  file.Close()
}
