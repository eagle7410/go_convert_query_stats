package json

import (
  "encoding/json"
  "io/ioutil"
  "strings"
)

func addExtJson(s *string) {
  if strings.Index(*s, ".json") == -1 {
    *s += ".json"
  }
}

func To(v interface{}) ([]byte, error)  {
  return json.Marshal(v)
}

func ToFile(v interface{}, filePath string)  ([]byte, error)   {

  json, err := To(v)

  addExtJson(&filePath)

  if err != nil {
    return json, err
  }

  err = ioutil.WriteFile(filePath, json, 0644)

  return json, err

}

func From(data []byte, v interface{}) error  {
  return json.Unmarshal(data, v)
}

func FromFile (filePath string, v interface{}) error {
  addExtJson(&filePath)

  bytesRead, err := ioutil.ReadFile(filePath)

  if err != nil {
    return err
  }

  return From(bytesRead, v)
}
