package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"reflect"

	"github.com/deepnoodle-ai/risor/v2"
	"github.com/deepnoodle-ai/risor/v2/pkg/object"
	"github.com/google/uuid"
)

func main() {
	prettyOutput := flag.Bool("pretty", false, "Pretty print JSON output")
	profile := flag.String("profile", "core", "Execution profile: core, io, admin")
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: risor-runner [--pretty] <script>")
		os.Exit(1)
	}

	script := os.Args[1]

	ctx := context.Background()

	// Start from the standard library builtins, then explicitly add our functions
	// based on the selected security profile.
	env := risor.Builtins()

	// Core-safe builtins (always enabled)
	env["json_parse"] = wrapFunc(jsonParse)
	env["json_stringify"] = wrapFunc(jsonStringify)
	env["json_to_yaml"] = wrapFunc(jsonToYaml)

	env["split"] = wrapFunc(split)
	env["join"] = wrapFunc(join)
	env["trim"] = wrapFunc(trim)
	env["upper"] = wrapFunc(upper)
	env["lower"] = wrapFunc(lower)
	env["replace"] = wrapFunc(replace)
	env["regex_match"] = wrapFunc(regexMatch)
	env["regex_replace"] = wrapFunc(regexReplace)
	env["contains"] = wrapFunc(contains)
	env["starts_with"] = wrapFunc(startsWith)
	env["ends_with"] = wrapFunc(endsWith)

	env["first"] = wrapFunc(first)
	env["last"] = wrapFunc(last)
	env["reverse"] = wrapFunc(reverseList)
	env["unique"] = wrapFunc(unique)
	env["flatten"] = wrapFunc(flatten)
	env["sort"] = wrapFunc(sortList)

	env["min"] = wrapFunc(minVal)
	env["max"] = wrapFunc(maxVal)
	env["sum"] = wrapFunc(sumVals)
	env["avg"] = wrapFunc(avgVals)
	env["round_val"] = wrapFunc(roundVal)
	env["floor_val"] = wrapFunc(floorVal)
	env["ceil_val"] = wrapFunc(ceilVal)
	env["abs_val"] = wrapFunc(absVal)

	env["now"] = wrapFunc(now)
	env["timestamp"] = wrapFunc(timestamp)
	env["format_time"] = wrapFunc(formatTime)
	env["parse_time"] = wrapFunc(parseTime)

	env["md5_hash"] = wrapFunc(md5Hash)
	env["sha256_hash"] = wrapFunc(sha256Hash)
	env["base64_encode"] = wrapFunc(base64Encode)
	env["base64_decode"] = wrapFunc(base64Decode)

	env["uuid"] = wrapFunc(uuidGen)
	env["random_int"] = wrapFunc(randomInt)
	env["random_choice"] = wrapFunc(randomChoice)

	env["url_encode"] = wrapFunc(urlEncode)
	env["url_decode"] = wrapFunc(urlDecode)
	env["html_encode"] = wrapFunc(htmlEncode)

	// Guarded I/O: enabled for profiles "io" and "admin"
	if *profile == "io" || *profile == "admin" {
		env["http_get"] = wrapFunc(httpGet)
		env["http_post"] = wrapFunc(httpPost)
		env["http_put"] = wrapFunc(httpPut)
		env["http_delete"] = wrapFunc(httpDelete)
		env["http_headers"] = wrapFunc(httpHeaders)

		env["file_read"] = wrapFunc(fileRead)
		env["file_write"] = wrapFunc(fileWrite)
		env["file_exists"] = wrapFunc(fileExists)
		env["file_delete"] = wrapFunc(fileDelete)
		env["file_list"] = wrapFunc(fileList)
	}

	// Dangerous / admin-only capabilities
	if *profile == "admin" {
		env["exec_cmd"] = wrapFunc(execCmd)
		env["env_get"] = wrapFunc(envGet)
		env["env_vars"] = wrapFunc(envVars)
		env["os_name"] = wrapFunc(osName)
		env["hostname"] = wrapFunc(hostname)
		env["env_var"] = wrapFunc(envVar)
	}

	result, err := risor.Eval(ctx, script, risor.WithEnv(env))
	if err != nil {
		fmt.Fprintf(os.Stderr, "{\"status\": \"error\", \"error\": \"%v\"}\n", err)
		os.Exit(1)
	}

	if result != nil {
		// Always output JSON
		outputMap := map[string]interface{}{
			"result": result,
		}
		output, err := json.Marshal(outputMap)
		if err != nil {
			fmt.Fprintf(os.Stderr, "{\"status\": \"error\", \"error\": \"%v\"}\n", err)
			os.Exit(1)
		}
		
		if *prettyOutput {
			var prettyJSON bytes.Buffer
			json.Indent(&prettyJSON, output, "", "  ")
			fmt.Println(prettyJSON.String())
		} else {
			fmt.Println(string(output))
		}
	}
}

type RisorFunc func(args ...interface{}) (interface{}, error)

var defaultRegistry = risor.NewTypeRegistry().Build()

func wrapFunc(fn RisorFunc) object.Object {
	return object.NewGoFunc(reflect.ValueOf(fn), "builtin", defaultRegistry)
}

func stringArg(args []interface{}, i int, name string) string {
	if i >= len(args) {
		return ""
	}
	if s, ok := args[i].(string); ok {
		return s
	}
	return fmt.Sprintf("%v", args[i])
}

func toFloat(v interface{}) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int64:
		return float64(x)
	case int:
		return float64(x)
	default:
		return 0
	}
}

// HTTP Functions
func httpGet(args ...interface{}) (interface{}, error) {
	resp, err := http.Get(stringArg(args, 0, "url"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return map[string]interface{}{"status": resp.StatusCode, "body": string(body), "headers": resp.Header}, nil
}

func httpPost(args ...interface{}) (interface{}, error) {
	resp, err := http.Post(stringArg(args, 0, "url"), "text/plain", strings.NewReader(stringArg(args, 1, "body")))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	return map[string]interface{}{"status": resp.StatusCode, "body": string(respBody)}, nil
}

func httpPut(args ...interface{}) (interface{}, error) {
	req, _ := http.NewRequest("PUT", stringArg(args, 0, "url"), strings.NewReader(stringArg(args, 1, "body")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	return map[string]interface{}{"status": resp.StatusCode, "body": string(respBody)}, nil
}

func httpDelete(args ...interface{}) (interface{}, error) {
	req, _ := http.NewRequest("DELETE", stringArg(args, 0, "url"), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	return map[string]interface{}{"status": resp.StatusCode, "body": string(respBody)}, nil
}

func httpHeaders(args ...interface{}) (interface{}, error) {
	resp, err := http.Head(stringArg(args, 0, "url"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return map[string]interface{}{"headers": resp.Header}, nil
}

// File Functions
func fileRead(args ...interface{}) (interface{}, error) {
	data, err := os.ReadFile(stringArg(args, 0, "path"))
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func fileWrite(args ...interface{}) (interface{}, error) {
	err := os.WriteFile(stringArg(args, 0, "path"), []byte(stringArg(args, 1, "content")), 0644)
	return err == nil, err
}

func fileExists(args ...interface{}) (interface{}, error) {
	_, err := os.Stat(stringArg(args, 0, "path"))
	return err == nil, nil
}

func fileDelete(args ...interface{}) (interface{}, error) {
	err := os.Remove(stringArg(args, 0, "path"))
	return err == nil, err
}

func fileList(args ...interface{}) (interface{}, error) {
	dir := "."
	if len(args) > 0 {
		dir = args[0].(string)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(entries))
	for _, e := range entries {
		result = append(result, e.Name())
	}
	return result, nil
}

// Exec & Env
func execCmd(args ...interface{}) (interface{}, error) {
	cmd := stringArg(args, 0, "cmd")
	var cmdArgs []string
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			cmdArgs = append(cmdArgs, args[i].(string))
		}
	}
	out, err := exec.Command(cmd, cmdArgs...).CombinedOutput()
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	return map[string]interface{}{"output": string(out), "error": errStr}, nil
}

func envGet(args ...interface{}) (interface{}, error) {
	return os.Getenv(stringArg(args, 0, "key")), nil
}

func envVars(args ...interface{}) (interface{}, error) {
	return os.Environ(), nil
}

func envVar(args ...interface{}) (interface{}, error) {
	key := stringArg(args, 0, "key")
	val, exists := os.LookupEnv(key)
	return map[string]interface{}{"value": val, "exists": exists}, nil
}

// JSON Functions
func jsonParse(args ...interface{}) (interface{}, error) {
	var v interface{}
	err := json.Unmarshal([]byte(stringArg(args, 0, "json")), &v)
	return v, err
}

func jsonStringify(args ...interface{}) (interface{}, error) {
	data, err := json.Marshal(args[0])
	return string(data), err
}

func jsonToYaml(args ...interface{}) (interface{}, error) {
	data, err := json.Marshal(args[0])
	if err != nil {
		return nil, err
	}
	var obj map[string]interface{}
	json.Unmarshal(data, &obj)
	return toYaml(obj, 0), nil
}

func toYaml(v interface{}, indent int) string {
	prefix := strings.Repeat("  ", indent)
	switch x := v.(type) {
	case map[string]interface{}:
		var lines []string
		for k, val := range x {
			lines = append(lines, fmt.Sprintf("%s%s: %s", prefix, k, toYaml(val, indent+1)))
		}
		return strings.Join(lines, "\n")
	case []interface{}:
		var lines []string
		for _, item := range x {
			lines = append(lines, fmt.Sprintf("%s- %s", prefix, toYaml(item, indent+1)))
		}
		return strings.Join(lines, "\n")
	case string:
		return fmt.Sprintf("\"%s\"", x)
	default:
		return fmt.Sprintf("%v", x)
	}
}

// String Functions
func split(args ...interface{}) (interface{}, error) {
	return strings.Split(stringArg(args, 0, "string"), stringArg(args, 1, "separator")), nil
}

func join(args ...interface{}) (interface{}, error) {
	list := args[0].([]interface{})
	sep := ","
	if len(args) > 1 {
		sep = args[1].(string)
	}
	strs := make([]string, len(list))
	for i, v := range list {
		strs[i] = fmt.Sprintf("%v", v)
	}
	return strings.Join(strs, sep), nil
}

func trim(args ...interface{}) (interface{}, error) {
	return strings.TrimSpace(stringArg(args, 0, "string")), nil
}

func upper(args ...interface{}) (interface{}, error) {
	return strings.ToUpper(stringArg(args, 0, "string")), nil
}

func lower(args ...interface{}) (interface{}, error) {
	return strings.ToLower(stringArg(args, 0, "string")), nil
}

func replace(args ...interface{}) (interface{}, error) {
	return strings.ReplaceAll(stringArg(args, 0, "string"), stringArg(args, 1, "old"), stringArg(args, 2, "new")), nil
}

func regexMatch(args ...interface{}) (interface{}, error) {
	return regexp.MatchString(stringArg(args, 1, "pattern"), stringArg(args, 0, "string"))
}

func regexReplace(args ...interface{}) (interface{}, error) {
	re := regexp.MustCompile(stringArg(args, 1, "pattern"))
	return re.ReplaceAllString(stringArg(args, 0, "string"), stringArg(args, 2, "replacement")), nil
}

func contains(args ...interface{}) (interface{}, error) {
	return strings.Contains(stringArg(args, 0, "string"), stringArg(args, 1, "substr")), nil
}

func startsWith(args ...interface{}) (interface{}, error) {
	return strings.HasPrefix(stringArg(args, 0, "string"), stringArg(args, 1, "prefix")), nil
}

func endsWith(args ...interface{}) (interface{}, error) {
	return strings.HasSuffix(stringArg(args, 0, "string"), stringArg(args, 1, "suffix")), nil
}

// List Functions
func first(args ...interface{}) (interface{}, error) {
	list := args[0].([]interface{})
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func last(args ...interface{}) (interface{}, error) {
	list := args[0].([]interface{})
	if len(list) == 0 {
		return nil, nil
	}
	return list[len(list)-1], nil
}

func reverseList(args ...interface{}) (interface{}, error) {
	list := args[0].([]interface{})
	result := make([]interface{}, len(list))
	for i, v := range list {
		result[len(list)-1-i] = v
	}
	return result, nil
}

func unique(args ...interface{}) (interface{}, error) {
	list := args[0].([]interface{})
	seen := make(map[string]bool)
	var result []interface{}
	for _, v := range list {
		key := fmt.Sprintf("%v", v)
		if !seen[key] {
			seen[key] = true
			result = append(result, v)
		}
	}
	return result, nil
}

func flatten(args ...interface{}) (interface{}, error) {
	return flatten2(args[0].([]interface{})), nil
}

func flatten2(args []interface{}) []interface{} {
	var result []interface{}
	for _, v := range args {
		if arr, ok := v.([]interface{}); ok {
			result = append(result, flatten2(arr)...)
		} else {
			result = append(result, v)
		}
	}
	return result
}

func sortList(args ...interface{}) (interface{}, error) {
	list := args[0].([]interface{})
	sorted := make([]interface{}, len(list))
	copy(sorted, list)
	sort.Slice(sorted, func(i, j int) bool {
		return fmt.Sprintf("%v", sorted[i]) < fmt.Sprintf("%v", sorted[j])
	})
	return sorted, nil
}

// Math Functions
func minVal(args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("min requires at least 2 arguments")
	}
	min := toFloat(args[0])
	for i := 1; i < len(args); i++ {
		if toFloat(args[i]) < min {
			min = toFloat(args[i])
		}
	}
	return min, nil
}

func maxVal(args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("max requires at least 2 arguments")
	}
	max := toFloat(args[0])
	for i := 1; i < len(args); i++ {
		if toFloat(args[i]) > max {
			max = toFloat(args[i])
		}
	}
	return max, nil
}

func sumVals(args ...interface{}) (interface{}, error) {
	list := args[0].([]interface{})
	var sum float64
	for _, v := range list {
		sum += toFloat(v)
	}
	return sum, nil
}

func avgVals(args ...interface{}) (interface{}, error) {
	list := args[0].([]interface{})
	if len(list) == 0 {
		return 0.0, nil
	}
	var sum float64
	for _, v := range list {
		sum += toFloat(v)
	}
	return sum / float64(len(list)), nil
}

func roundVal(args ...interface{}) (interface{}, error) {
	return math.Round(toFloat(args[0])), nil
}

func floorVal(args ...interface{}) (interface{}, error) {
	return math.Floor(toFloat(args[0])), nil
}

func ceilVal(args ...interface{}) (interface{}, error) {
	return math.Ceil(toFloat(args[0])), nil
}

func absVal(args ...interface{}) (interface{}, error) {
	return math.Abs(toFloat(args[0])), nil
}

// Time Functions
func now(args ...interface{}) (interface{}, error) {
	return time.Now().Format(time.RFC3339), nil
}

func timestamp(args ...interface{}) (interface{}, error) {
	return time.Now().Unix(), nil
}

func formatTime(args ...interface{}) (interface{}, error) {
	format := "2006-01-02 15:04:05"
	if len(args) > 1 {
		format = args[1].(string)
	}
	return time.Unix(int64(toFloat(args[0])), 0).Format(format), nil
}

func parseTime(args ...interface{}) (interface{}, error) {
	format := "2006-01-02T15:04:05Z07:00"
	if len(args) > 1 {
		format = args[1].(string)
	}
	t, err := time.Parse(format, stringArg(args, 0, "time"))
	if err != nil {
		return nil, err
	}
	return t.Unix(), nil
}

// Crypto Functions
func md5Hash(args ...interface{}) (interface{}, error) {
	return fmt.Sprintf("%x", md5.Sum([]byte(stringArg(args, 0, "data")))), nil
}

func sha256Hash(args ...interface{}) (interface{}, error) {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(stringArg(args, 0, "data")))), nil
}

func base64Encode(args ...interface{}) (interface{}, error) {
	return base64.StdEncoding.EncodeToString([]byte(stringArg(args, 0, "data"))), nil
}

func base64Decode(args ...interface{}) (interface{}, error) {
	decoded, err := base64.StdEncoding.DecodeString(stringArg(args, 0, "data"))
	return string(decoded), err
}

// System Functions
func osName(args ...interface{}) (interface{}, error) {
	return runtime.GOOS, nil
}

func hostname(args ...interface{}) (interface{}, error) {
	return os.Hostname()
}

// Random & ID
func uuidGen(args ...interface{}) (interface{}, error) {
	return uuid.New().String(), nil
}

func randomInt(args ...interface{}) (interface{}, error) {
	min, max := 0, 100
	if len(args) >= 1 {
		min = int(toFloat(args[0]))
	}
	if len(args) >= 2 {
		max = int(toFloat(args[1]))
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min, nil
}

func randomChoice(args ...interface{}) (interface{}, error) {
	list := args[0].([]interface{})
	if len(list) == 0 {
		return nil, nil
	}
	rand.Seed(time.Now().UnixNano())
	return list[rand.Intn(len(list))], nil
}

// Encoding
func urlEncode(args ...interface{}) (interface{}, error) {
	return url.QueryEscape(stringArg(args, 0, "string")), nil
}

func urlDecode(args ...interface{}) (interface{}, error) {
	return url.QueryUnescape(stringArg(args, 0, "string"))
}

func htmlEncode(args ...interface{}) (interface{}, error) {
	s := stringArg(args, 0, "string")
	var sb strings.Builder
	for _, r := range s {
		switch r {
		case '&':
			sb.WriteString("&amp;")
		case '<':
			sb.WriteString("&lt;")
		case '>':
			sb.WriteString("&gt;")
		case '"':
			sb.WriteString("&quot;")
		case '\'':
			sb.WriteString("&#39;")
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String(), nil
}
