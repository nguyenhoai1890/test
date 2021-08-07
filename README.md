### Run test case:
`go test ./...`

In main_test.go, there are 2 tests which coverts A and B question.

About func `TestFormatJson`, there is only one test case to check how the function handle with input data because the test cases for main func which runs for question A already covered all edge cases for the format data.  


### Explanation:

##### **A: Extend your solution to support an escaping syntax for the format_string function:**

###### First Init solution, pseudocode that based on current code's flow.

```
const (
 NotOpen = 0
 SingleOpen = 1
 DoubleOpen = 2
)

const (
    keyMarkOpen = "{"
    keyMarkClose = keyMarkClose
)


// abc {{key1}} {key1} -> abc {key1} val1
// abc {{{key1} deg -> abc {{val1 deg
// abc {{{{key1}} deg -> abc {{{key1} deg
// abc {{{{key1}a} deg -> abc {{{val1}a} deg
// abc -> abc
func formatString(template, parameters)
    i,j := 0
    currentOpen:= NotOpen
    key := ""
    result := ""
    Loop i < len(template) {
        char := template[i]
        if char == keyMarkOpen {
            switch currentOpen {
                case NotOpen:
                    currentOpen = SingleOpen
                    key = keyMarkOpen
                case SingleOpen:
                    if key == keyMarkOpen { // Case {{
                        currentOpen = DoubleOpen
                        key = "{{"
                    } else { // Case {abc{ -> add {abc into result and keep current status is SingleOpen
                        result += key
                        key = keyMarkOpen
                    }
                case DoubleOpen:
                    // Case {{{ -> add "{" (the 1st one) into result and keep current status is Double Open
                    if key = DoubleOpen {
                        result += "{"
                    } else { // Case {{abc{ -> add "{{abc" into result and change current status to SingleOpen.
                        result += key
                        key = keyMarkOpen
                        currentOpen = SingleOpen
                    }
                }
            }
        } else if char == keyMarkClose {
            switch currentOpen {
                case NotOpen: // Case abc{ -> add "{" into result
                    result += keyMarkClose
                case SingleOpen:
                    key += keyMarkClose
                    result += getKeyValue(key, params)
                    currentOpen = NotOpen
                case DoubleOpen:
                    key += keyMarkClose
                    if len(key) > 2 && key[len -2:] == "}}" { // Case {{abc}} or {{}}
                        result += getKeyValue(key, params)
                        currentOpen = NotOpen
                        key =""
                    }
            }
        } else {
            switch currentOpen {
                case NotOpen:
                    result += char
                case DoubleOpen: , {{ke..., {{
                    if len(key) > 0 && key[len-1] = keyMarkClose { // case: {{key}a -> {val1a
                        result += "{" + getKeyValue(key, params) + char
                        key = ""
                        currentOpen = NotOpen
                    } else { // case: {{key..
                        key += char
                    }
                case SingleOpen:
                    key += char
            }
        }
    }
    if len(key) > 0 { //Case when there is no close in the end. E.g: abc{abc -> abc{abc
        result += key
    }
}

```

###### Second thought
Since the current init logic makes the lines of code is huge because it needs to handle 3 stages of key: NotOpen, DoubleOpen and SingleOpen.
So I removed that flow out of my mind to rethink again the solution. As result, the second solution is just care about the value of key (first 2 chars of key) when current character is "}"

My second pseudocode, pls note that the code I implemented quite different with this, because:
* Found some bugs when I check with the test case in main_test.go. 
* Improve the flow which don't need to use recursion. 


```
func formatString(template, parameters) string {
    result = processString("", 0, template)
}

func pString(key string, index int, template) (index, string) {
    result = ""
    if index < size {
        if t[index] == "{" {
            if key Not == "{" && "{{" {
                result += key
            } else {
                if key = "{{" {
                    result += "{"
                    key = "{"
                }
            }
            v := ""
            v, index = pString(key + "{", index +1)
            result +=v
        } else if t[index] == "}" {
            key += "}"
            if t[index+1] == "}" {
                key += "}}"
                index ++
            }
            result += getVal(key)
            key = ""
            result += pString(key, index, template)
        } else {
            key += t[index]
        }
    }
    return result + key
}
```

##### B: Support “grouping” or “containing” relationship data
Since the previous code flow is supporting well with any key format, as long as key is in correct value with params.
So Instead of modifying the flow of code, I think about format the current data into correct params having type: `map[string]string`

###### So the come up solution is:

* Consider that the passed data should be json, I unmarshalled the json into map with format: "<parent field name>.<sub field name>....": "value"
* Pass result params into the previous function.
