package main

import (
  "github.com/stretchr/testify/require"
  "testing"
)

type TestCaseFormatString struct {
  Input string
  Params map[string]string
  Result string
}
var testCasesFormatString = []TestCaseFormatString{
  // invalid input
  {
    Input: "{}",
    Result: "{}",
  },
  {
    Input: "{",
    Result: "{",
  },
  {
    Input: "a",
    Result: "a",
  },
  // Normal cases and edge cases
  {
    Params: map[string]string{"name": "Hoai", "c": "PrimeData"},
    Input: "{c}",
    Result: "PrimeData",
  },
  {
    Params: map[string]string{"name": "Hoai", "company": "PrimeData"},
    Input: "Hi {name}, welcome to {company}!",
    Result: "Hi Hoai, welcome to PrimeData!",
  },
  {
    Params: map[string]string{"name": "Hoai", "company": "PrimeData"},
    Input: "Hi {name_notfound}, welcome to {company}!",
    Result: "Hi {name_notfound}, welcome to PrimeData!",
  },
  {
    Params: map[string]string{"name": "Hoai", "company": "PrimeData"},
    Input: "{name is nothing, nha {name}",
    Result: "{name is nothing, nha Hoai",
  },
  {
    Params: map[string]string{"name": "Hoai", "company": "PrimeData"},
    Input: "name} is nothing, nha {name}",
    Result: "name} is nothing, nha Hoai",
  },
  {
    Params: map[string]string{"name": "Hoai", "company": "PrimeData"},
    Input: "Hi {{name}",
    Result: "Hi {Hoai",
  },
  { // I found error on this test case.
    Params: map[string]string{"name": "Hoai", "company": "PrimeData"},
    Input: "Hi {{{name}",
    Result: "Hi {{Hoai",
  },
  {
    Params: map[string]string{"name": "Hoai", "company": "PrimeData"},
    Input: "Hi {{{name}",
    Result: "Hi {{Hoai",
  },
  // Complex string with double {}
  {
    Params: map[string]string{"name": "Hoai", "company": "PrimeData"},
    Input: "Hi {{{name}, there is {{company}} {{{company}}}, with {company}} {{company} {{{company}",
    Result: "Hi {{Hoai, there is {company} {{company}}, with PrimeData} {PrimeData {{PrimeData",
  },
}
func TestFormatString(t *testing.T) {
  for _, tc := range testCasesFormatString {
    result := FormatString(tc.Input, tc.Params)
    require.Equal(t, tc.Result, result, tc.Input)
  }
}

type TestCaseFormatStringWithJson struct {
  Input string
  Data string
  Result string
}
var jsonData = `
{
  "person": {
    "name": "Hoai",
    "age": 18,
    "hobbies": ["sleep", "eat"]
  },
  "company": {
    "name": "PrimeData",
    "Address": "12B Nguyen Huu Canh, Binh Thanh District, HCMC.",
    "Contact" : {
      "Phone": "(+84) 888 818 688",
      "Email": "info@primedata.ai"
    }
  }
}`

var formatJsonTestCases = []TestCaseFormatStringWithJson{
  {
    Input: "Hey {person.name}, with age {person.age} and hobbies {person.hobbies.0} and {person.hobbies.1}, welcome to {company.name}, {company.address}, contact information is phone: {company.contact.phone} and email: {company.contact.email}",
    Data: jsonData,
    Result: "Hey Hoai, with age 18 and hobbies sleep and eat, welcome to PrimeData, 12B Nguyen Huu Canh, Binh Thanh District, HCMC., contact information is phone: (+84) 888 818 688 and email: info@primedata.ai",
  },
}
func TestFormatJson(t *testing.T) {
  for _, tc := range formatJsonTestCases {
    result, err := FormatStringWithJson(tc.Input, tc.Data)
    require.Nil(t, err, tc)
    require.Equal(t, tc.Result, result, tc)
  }
}