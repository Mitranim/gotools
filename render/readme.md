# Work In Progress

-------------------------------------------------------------------------------

## Description
## Features
## Contents
## Installation
## API Reference

## Type gotchas

The type used for data is `map[string]interface{}`. With a map like this, you can assign any values, and can safely reference them in templates, both by uppercase and lowercase keys.

The first gotcha is that we don't use pointers for data maps. Golang maps reference their underlying data structures with a pointer. Passing a map by pointer is basically equivalent to passing it by value, and incurs inline type conversions when accessing fields.

The second, more important gotcha is how to read and modify interface values programmatically. When reading a value into a variable, you need to hint its type with a type assertion:

```golang
data := map[string]interface{}{}

title := data["Title"].(string)
```

If the value is missing or has a different type, the variable will be set to a [zero value](https://golang.org/ref/spec#The_zero_value) for that type, like if you allocated it with `var <varname> <type>`.

```golang
data := map[string]interface{}{
  "Title": 123,
}

title := data["Title"].(string)

title == ""  // true
```

If we used map pointers, we'd have to deference them before accessing fields:

```golang
data := &map[string]interface{}{}

ready := (*data)["Ready"].(bool)
```

When modifying a value, put it through a variable:

```golang
data := map[string]interface{}{}
data["cache"] = map[string]bool{}

// Trying to do:
//   data["cache"]["first"] = true
// will cause a compile error, because the compiler doesn't know
// the implied cache type.

cache := data["cache"].(map[string]bool)
cache["first"] = true
data["cache"] = cache
```

You can also do an inline type assertion in an assignment.

```golang
data := map[string]interface{}{}
data["cache"] = map[string]bool{}
data["cache"].(map[string]bool)["first"] = true

// To check if the cache is nil, do this:
if value, _ := data["cache"].(map[string]bool); value == nil {
  data["cache"] = map[string]bool{}
}
```
