package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	desc := flag.String("desc", "", "Description of the skill to generate")
	city := flag.String("city", "", "City for weather-like skills")
	url := flag.String("url", "https://api.example.com", "URL for API skills")
	path := flag.String("path", "/path/to/file.txt", "File path for file operations")
	content := flag.String("content", "hello world", "Content for file write")
	flag.Parse()

	if *desc == "" {
		fmt.Println("Usage: risor-gen --desc \"description\" [options]")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  risor-gen --desc \"fetch weather for city\" --city Riga")
		fmt.Println("  risor-gen --desc \"read file\" --path /tmp/test.txt")
		fmt.Println("  risor-gen --desc \"fetch API data\" --url https://api.example.com")
		fmt.Println("")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	script := generateScript(*desc, *city, *url, *path, *content)
	
	fmt.Println("==================================================")
	fmt.Println("Generated Risor Script:")
	fmt.Println("==================================================")
	fmt.Println(script)
	fmt.Println("==================================================")
}

func generateScript(desc, city, url, path, content string) string {
	desc = strings.ToLower(desc)

	// Weather template
	if strings.Contains(desc, "weather") || strings.Contains(desc, "temperature") {
		if city == "" {
			city = "Riga"
		}
		return fmt.Sprintf(`# Weather skill for %s
city = "%s"
result = http_get("wttr.in/" + city + "?format=j1")
print(result.body)`, city, city)
	}

	// HTTP GET template
	if strings.Contains(desc, "fetch") || strings.Contains(desc, "get") || strings.Contains(desc, "api") || strings.Contains(desc, "http") {
		if strings.Contains(desc, "json") || strings.Contains(desc, "parse") {
			return fmt.Sprintf(`# API fetch with JSON parsing
url = "%s"
result = http_get(url)
data = json_parse(result.body)
print(data)`, url)
		}
		return fmt.Sprintf(`# HTTP GET request
result = http_get("%s")
print(result.body)`, url)
	}

	// HTTP POST template
	if strings.Contains(desc, "post") && strings.Contains(desc, "api") {
		return fmt.Sprintf(`# HTTP POST request
url = "%s"
body = '{"key": "value"}'
result = http_post(url, body)
print(result.body)`, url)
	}

	// File read template
	if strings.Contains(desc, "file") && strings.Contains(desc, "read") {
		return fmt.Sprintf(`# File read skill
path = "%s"
content = file_read(path)
print(content)`, path)
	}

	// File write template
	if strings.Contains(desc, "file") && strings.Contains(desc, "write") {
		return fmt.Sprintf(`# File write skill
path = "%s"
content = """%s"""
result = file_write(path, content)
print(result)`, path, content)
	}

	// File exists template
	if strings.Contains(desc, "file") && strings.Contains(desc, "exists") {
		return fmt.Sprintf(`# File exists check
path = "%s"
exists = file_exists(path)
print(exists)`, path)
	}

	// JSON parse template
	if strings.Contains(desc, "json") && strings.Contains(desc, "parse") {
		return `# JSON parse skill
json_str = '{"key": "value"}'
data = json_parse(json_str)
print(data)`
	}

	// Environment variable template
	if strings.Contains(desc, "env") || strings.Contains(desc, "environment") {
		return `# Environment variable skill
key = "PATH"
value = env_get(key)
print(value)`
	}

	// Default template
	return `# Generated skill: ` + desc + `
# Add your Risor code here

# Example:
# result = http_get("https://api.example.com")
# print(result.body)
`
}
