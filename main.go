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

	"github.com/google/uuid"
	"github.com/risor-io/risor"
	"github.com/risor-io/risor/object"
)

func main() {
	prettyOutput := flag.Bool("pretty", false, "Pretty print JSON output")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Usage: risor-runner [--pretty] <script>")
		os.Exit(1)
	}

	script := flag.Arg(0)

	ctx := context.Background()

	result, err := risor.Eval(ctx, script,
		// HTTP
		risor.WithGlobal("http_get", wrapFunc(httpGet)),
		risor.WithGlobal("http_post", wrapFunc(httpPost)),
		risor.WithGlobal("http_put", wrapFunc(httpPut)),
		risor.WithGlobal("http_delete", wrapFunc(httpDelete)),
		risor.WithGlobal("http_headers", wrapFunc(httpHeaders)),
		// File
		risor.WithGlobal("file_read", wrapFunc(fileRead)),
		risor.WithGlobal("file_write", wrapFunc(fileWrite)),
		risor.WithGlobal("file_exists", wrapFunc(fileExists)),
		risor.WithGlobal("file_delete", wrapFunc(fileDelete)),
		risor.WithGlobal("file_list", wrapFunc(fileList)),
		// Exec & Env
		risor.WithGlobal("exec_cmd", wrapFunc(execCmd)),
		risor.WithGlobal("env_get", wrapFunc(envGet)),
		risor.WithGlobal("env_vars", wrapFunc(envVars)),
		// JSON
		risor.WithGlobal("json_parse", wrapFunc(jsonParse)),
		risor.WithGlobal("json_stringify", wrapFunc(jsonStringify)),
		risor.WithGlobal("json_to_yaml", wrapFunc(jsonToYaml)),
		// String
		risor.WithGlobal("split", wrapFunc(split)),
		risor.WithGlobal("join", wrapFunc(join)),
		risor.WithGlobal("trim", wrapFunc(trim)),
		risor.WithGlobal("upper", wrapFunc(upper)),
		risor.WithGlobal("lower", wrapFunc(lower)),
		risor.WithGlobal("replace", wrapFunc(replace)),
		risor.WithGlobal("regex_match", wrapFunc(regexMatch)),
		risor.WithGlobal("regex_replace", wrapFunc(regexReplace)),
		risor.WithGlobal("contains", wrapFunc(contains)),
		risor.WithGlobal("starts_with", wrapFunc(startsWith)),
		risor.WithGlobal("ends_with", wrapFunc(endsWith)),
		// List
		risor.WithGlobal("first", wrapFunc(first)),
		risor.WithGlobal("last", wrapFunc(last)),
		risor.WithGlobal("reverse", wrapFunc(reverseList)),
		risor.WithGlobal("unique", wrapFunc(unique)),
		risor.WithGlobal("flatten", wrapFunc(flatten)),
		risor.WithGlobal("sort", wrapFunc(sortList)),
		// Math
		risor.WithGlobal("min", wrapFunc(minVal)),
		risor.WithGlobal("max", wrapFunc(maxVal)),
		risor.WithGlobal("sum", wrapFunc(sumVals)),
		risor.WithGlobal("avg", wrapFunc(avgVals)),
		risor.WithGlobal("round_val", wrapFunc(roundVal)),
		risor.WithGlobal("floor_val", wrapFunc(floorVal)),
		risor.WithGlobal("ceil_val", wrapFunc(ceilVal)),
		risor.WithGlobal("abs_val", wrapFunc(absVal)),
		// Time
		risor.WithGlobal("now", wrapFunc(now)),
		risor.WithGlobal("timestamp", wrapFunc(timestamp)),
		risor.WithGlobal("format_time", wrapFunc(formatTime)),
		risor.WithGlobal("parse_time", wrapFunc(parseTime)),
		// Crypto
		risor.WithGlobal("md5_hash", wrapFunc(md5Hash)),
		risor.WithGlobal("sha256_hash", wrapFunc(sha256Hash)),
		risor.WithGlobal("base64_encode", wrapFunc(base64Encode)),
		risor.WithGlobal("base64_decode", wrapFunc(base64Decode)),
		// System
		risor.WithGlobal("os_name", wrapFunc(osName)),
		risor.WithGlobal("hostname", wrapFunc(hostname)),
		risor.WithGlobal("env_var", wrapFunc(envVar)),
		// Random & ID
		risor.WithGlobal("uuid", wrapFunc(uuidGen)),
		risor.WithGlobal("random_int", wrapFunc(randomInt)),
		risor.WithGlobal("random_choice", wrapFunc(randomChoice)),
		// Encoding
		risor.WithGlobal("url_encode", wrapFunc(urlEncode)),
		risor.WithGlobal("url_decode", wrapFunc(urlDecode)),
		risor.WithGlobal("html_encode", wrapFunc(htmlEncode)),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "{\"status\": \"error\", \"error\": \"%v\"}\n", err)
		os.Exit(1)
	}

	if result != nil {
		// Always output JSON
		outputMap := map[string]interface{}{
			"result": toGoValue(result),
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

func wrapFunc(fn RisorFunc) object.Object {
	return object.NewBuiltin(fmt.Sprintf("fn_%p", fn), func(ctx context.Context, args ...object.Object) object.Object {
		goArgs := make([]interface{}, len(args))
		for i, a := range args {
			goArgs[i] = toGoValue(a)
		}
		result, err := fn(goArgs...)
		if err != nil {
			return object.NewError(fmt.Errorf("%v", err))
		}
		return toObject(result)
	})
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
	req, err := http.NewRequest("PUT", stringArg(args, 0, "url"), strings.NewReader(stringArg(args, 1, "body")))
	if err != nil {
		return nil, err
	}
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
	req, err := http.NewRequest("DELETE", stringArg(args, 0, "url"), nil)
	if err != nil {
		return nil, err
	}
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
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
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
	return rand.Intn(max-min) + min, nil
}

func randomChoice(args ...interface{}) (interface{}, error) {
	list := args[0].([]interface{})
	if len(list) == 0 {
		return nil, nil
	}
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

// Helpers
func toGoValue(result object.Object) interface{} {
	switch result.Type() {
	case object.STRING:
		return result.(*object.String).Value()
	case object.INT:
		return result.(*object.Int).Value()
	case object.FLOAT:
		return result.(*object.Float).Value()
	case object.BOOL:
		return result.(*object.Bool).Value()
	case object.LIST:
		list := result.(*object.List)
		items := list.Value()
		res := make([]interface{}, 0, len(items))
		for _, item := range items {
			res = append(res, toGoValue(item))
		}
		return res
	case object.MAP:
		m := result.(*object.Map)
		pairs := m.Value()
		res := make(map[string]interface{})
		for key, val := range pairs {
			res[key] = toGoValue(val)
		}
		return res
	case object.NIL:
		return nil
	default:
		return result.Inspect()
	}
}

func toObject(v interface{}) object.Object {
	switch x := v.(type) {
	case string:
		return object.NewString(x)
	case float64:
		return object.NewFloat(x)
	case bool:
		return object.NewBool(x)
	case nil:
		return &object.NilType{}
	case int:
		return object.NewInt(int64(x))
	case int64:
		return object.NewInt(x)
	case []interface{}:
		items := make([]object.Object, len(x))
		for i, item := range x {
			items[i] = toObject(item)
		}
		return object.NewList(items)
	case map[string]interface{}:
		m := make(map[string]object.Object)
		for k, val := range x {
			m[k] = toObject(val)
		}
		return object.NewMap(m)
	default:
		return object.NewString(fmt.Sprintf("%v", v))
	}
}
