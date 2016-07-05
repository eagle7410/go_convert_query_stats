package json

import (
  "encoding/json"
  "io/ioutil"
  "strings"
  "os"
)

var content []byte

func addExtJson(s *string) {
  if strings.Index(*s, ".json") == -1 {
    *s += ".json"
  }
}

func To(v interface{}) ([]byte, error)  {
  return json.Marshal(v)
}

func ToFile(v interface{}, filePath string)  (error)   {

  json, err := To(v)

  addExtJson(&filePath)

  if err != nil {
    return err
  }

  return ioutil.WriteFile(filePath, json, 0644)

}

func clearContent(f *os.File) {
  content = []byte{}
  f.Close()
}

func ToFileSync(v interface{}, filePath string)  (error)   {
  var (
    err error
    f   *os.File
  )

  content, err = To(v)

  if err != nil {
    return err
  }

  addExtJson(&filePath)

  f, err = os.Create(filePath)

  defer clearContent(f)

  if err != nil {
    return err
  }

  _, err = f.Write(content)

  return err


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
