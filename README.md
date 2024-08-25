# StructMerge

`StructMerge` is a Go package that allows you to merge two structs of the same type based on configurable options.

## Features

- **Merge two structs**: Combine fields from two structs of the same type.
- **Configurable merging**: Choose to include or exclude specific fields, exclude empty fields, or overwrite empty fields in the destination struct.
- **Nested struct support**: Automatically merge nested structs.
- **Defines an interface *Merger* for merging custom data structs and fields**

## Installation

To install the package, run:

```bash
go get github.com/abiiranathan/structmerge
```

## Usage

### Import the package

```go
import "github.com/abiiranathan/structmerge"
```

### Merge Two Structs

The primary function is `Merge`, which allows you to merge the fields of two structs based on a set of options.

#### Example Structs

```go
type Address struct {
    Street  string
    City    string
    Country string
}

type Person struct {
    Name    string
    Age     int
    Address Address
    Active  bool
    Score   int
}
```

#### Example Usage

```go
package main

import (
    "fmt"
    "github.com/abiiranathan/structmerge"
)

func main() {
    person1 := Person{
        Name: "Alice",
        Age:  25,
        Address: Address{
            Street:  "123 Old St",
            City:    "Old City",
            Country: "",
        },
        Active: true,
        Score:  5,
    }

    person2 := Person{
        Name: "Alice",
        Age:  25,
        Address: Address{
            Street:  "456 New St",
            City:    "New City",
            Country: "New Country",
        },
        Active: true,
        Score:  10,
    }

    // Configure the merge operation
    cfg := structmerge.Config{
        Option: structmerge.OverwriteEmpty,
    }

    err := structmerge.Merge(&person1, person2, cfg)
    if err != nil {
        fmt.Println("Error merging structs:", err)
        return
    }

    fmt.Printf("Merged Person: %+v\n", person1)
}
```

#### Output

```
Merged Person: {Name:Alice Age:25 Address:{Street:456 New St City:New City Country:New Country} Active:true Score:5}
```

### Options

The `Merge` function accepts several options to control the behavior of the merge:

- **`IncludeAll`**: Includes all fields from the source struct in the merge.
- **`ExcludeEmpty`**: Excludes empty fields from the source struct when merging.
- **`OverwriteEmpty`**: Overwrites empty fields in the destination struct with non-empty fields from the source struct.

#### Example: Exclude Empty Fields

```go
cfg := structmerge.Config{
    Option: structmerge.ExcludeEmpty,
}

err := structmerge.Merge(&person1, person2, cfg)
if err != nil {
    fmt.Println("Error merging structs:", err)
    return
}

fmt.Printf("Merged Person with ExcludeEmpty: %+v\n", person1)
```

### Include and Exclude Specific Fields

You can specify which fields to include or exclude during the merge:

#### Example: Include Specific Fields

```go
cfg := structmerge.Config{
    Option:  structmerge.IncludeAll,
    Include: []string{"Address.Street", "Address.Country"},
}

err := structmerge.Merge(&person1, person2, cfg)
if err != nil {
    fmt.Println("Error merging structs:", err)
    return
}

fmt.Printf("Merged Person with Include: %+v\n", person1)
```

#### Example: Exclude Specific Fields

```go
cfg := structmerge.Config{
    Option:  structmerge.IncludeAll,
    Exclude: []string{"Address.City"},
}

err := structmerge.Merge(&person1, person2, cfg)
if err != nil {
    fmt.Println("Error merging structs:", err)
    return
}

fmt.Printf("Merged Person with Exclude: %+v\n", person1)
```

## merging custom struct types

You can implement the Merger interface to handle complex types on struct level or
field-level.

For example:
```go
type Date time.Time

func (d *Date) Merge(src reflect.Value) error {
	srcDate, ok := src.Interface().(Date)
	if !ok {
		return ErrInvalidSource
	}
	*d = srcDate
	return nil
}
```

## Error Handling

The `Merge` function will return an error in the following cases:

- **`ErrInvalidDestination`**: The destination parameter is not a pointer to a struct.
- **`ErrInvalidSource`**: The source parameter is not a struct.
- **`ErrTypeMismatch`**: The source and destination structs are not of the same type.

You can check these errors as follows:

```go
if err := structmerge.Merge(&person1, person2, cfg); err != nil {
    if errors.Is(err, structmerge.ErrInvalidDestination) {
        fmt.Println("Invalid destination struct")
    } else if errors.Is(err, structmerge.ErrInvalidSource) {
        fmt.Println("Invalid source struct")
    } else if errors.Is(err, structmerge.ErrTypeMismatch) {
        fmt.Println("Source and destination types do not match")
    } else {
        fmt.Println("An unexpected error occurred:", err)
    }
}
```

## Contributing

Feel free to fork the repository and submit pull requests with improvements or bug fixes. Please ensure that any new code is covered by tests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
