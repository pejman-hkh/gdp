# GDP: GoLang DOM Parser

GDP is a lightweight HTML DOM parser for Go, designed for parsing HTML documents using CSS selectors.

## Installation

To install GDP, use go get command:
```
go get github.com/pejman-hkh/gdp
```

## Usage
 ```
package main

import (
	"fmt"

	"github.com/pejman-hkh/gdp/gdp"
)

func main() {
	// Example HTML content
	htmlContent := `
        <html>
            <body>
                <div id="content">
                    <h1>Title</h1>
                    <p>Paragraph content</p>
                </div>
            </body>
        </html>
    `

	// Parse the HTML content
	doc := gdp.Default(htmlContent)

	// Example: Extract text from <h1> tag
	title := doc.Find("#content h1").Eq(0).Text()
	fmt.Println("Title:", title)

	// Example: Extract text from <p> tag
	paragraph := doc.Find("#content p").Eq(0).Text()
	fmt.Println("Paragraph:", paragraph)
}

```

https://go.dev/play/p/HALIPVQef_B

## API Reference
https://pkg.go.dev/github.com/pejman-hkh/gdp/gdp

## Example

This example demonstrates basic usage of GDP to parse HTML and query elements using CSS selectors.

## Contributing

Contributions are welcome! Fork the repository, make your changes, and submit a pull request.
