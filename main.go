package main

import (
  "encoding/json"
  "fmt"
  "reflect"
  "strconv"
  "strings"
)

func main()  {
  fmt.Println("Please run test funcs to check the result.")
}


// FormatString format value string with mapped key in params.
// Mapping key is case insensitive
func FormatString(value string, params map[string]string) string {
  // value should have at least 3 chars to be able to have key.
  if len(value) <= 2 {
    return value
  }

  // To avoid mistake when passing Upper or lower case, the params' keys need to be convert to lower case all.
  for k, v := range params {
    delete(params, k)
    params[strings.ToLower(k)] = v
  }

  // The result string can be use strings.Builder to process a large value string.
  var result, key string
  index := 0
  size := len(value)

  for index < size {
    currentChar := fmt.Sprintf("%c", value[index])
    if currentChar == "{" {
      // if key already had value. E.g: {abc{, {{abc{, -> key need to be reset before continue processing
      if key != "{" && key != "{{" {
        result += key
        key = ""
      } else {
        // this covers case that has "{{{" -> result need to be added one 1 and continue with key's value "{{"
        // this logic mentioned in the last testcase.
        if key == "{{" {
          result += "{"
          key = "{"
        }
      }
      key += "{"
    } else if currentChar == "}" && strings.Index(key, "{") == 0 { // There are only 3 cases.E.g abc}, {{abc}, and {abc}. We only care about 2 last cases.
      keyValue := ""
      key += "}"
      // As above logic when we detect "{". The value of key can only be started with "{" or "{{".
      // So we only need to check if this key is {{ }} or not.
      if strings.Index(key, "{{") == 0  {
        if index < size - 1 && fmt.Sprintf("%c", value[index+1]) == "}" { // key has format {{..}}
          index ++
          result += key[1:] //We cut down 1 { } and put value again into result.
          key = ""
        } else { // cover key with format {{....}.
          keyValue = key[1: len(key) -1]
          result += "{"
          keyValue = key[2: len(key) -1]
        }
      } else {
        keyValue = key[1: len(key) -1]
      }
      if keyValue != "" {
        if v, found := params[strings.ToLower(keyValue)]; found { // if the key is in params, we will replace with value in params, if not we will use keyValue.
          result += v
        } else {
          result = fmt.Sprintf("%s{%s}", result, keyValue)
        }
      }
      key = ""
    } else {
      key += currentChar
    }
    index++
  }
  // There are some cases that key has value but not completed as a fully pattern in the end of value.
  return result + key
}

func FormatStringWithJson(value string, jsonData string) (string, error) {
  data, err := convertJsonToMapData(jsonData)
  if err !=nil {
    return "", err
  }
  return FormatString(value, data), nil
}

// convertJsonToMapData convert json into map with key value is <field name>.<sub field name>
// the key will be converted into lower case
func convertJsonToMapData(jsonData string) (mapData map[string]string, err error) {
  jsonMap := map[string]interface{}{}
  if err = json.Unmarshal([]byte(jsonData), &jsonMap); err != nil {
    return nil, err
  }
  mapData = exposeJsonData(jsonMap)
  return mapData, nil
}
func exposeJsonData(mapData map[string]interface{}) map[string]string {
  params := map[string]string{}
  for k, v := range mapData {
    r := reflect.TypeOf(v)
    kind := r.Kind()
    fmt.Println(kind)
    fmt.Println(k)
    switch reflect.TypeOf(v).Kind() {
    case reflect.Slice:
      for index, sliceValue := range v.([]interface{}) {
        params[strings.ToLower(fmt.Sprintf("%s.%s", k, strconv.Itoa(index)))] = fmt.Sprint(sliceValue)
      }
    case reflect.Map: //so there is a map fields again, we will expose the value agains to get sub map Params
      nParams := exposeJsonData(v.(map[string]interface{}))
      for nK, nV := range nParams {
        params[strings.ToLower(fmt.Sprintf("%s.%s", k, nK))] = nV
      }
    default:
      params[k] = fmt.Sprint(v)
    }
  }
  return params
}
