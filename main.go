package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type SearchResult struct {
	Response struct {
		Docs []Doc `json:"docs"`
	} `json:"response"`
}

type Doc struct {
	G             string `json:"g"`
	A             string `json:"a"`
	V             string `json:"v"`
	LatestVersion string `json:"latestVersion"`
}

type Config struct {
	Query      string
	Rows       int
	Format     string
	IsPrecise  bool
	GroupId    string
	ArtifactId string
}

func main() {
	flag.Usage = func() {
		fmt.Println("Maven 仓库搜索工具")
		fmt.Println("\n用法:")
		fmt.Println("  mvnreps <query> [format]")
		fmt.Println("\n示例:")
		fmt.Println("  mvnreps okhttp                (模糊搜索)")
		fmt.Println("  mvnreps g:a                   (查看库的最近5个版本)")
		fmt.Println("  mvnreps g:a,10 gradle         (查看最近10个版本并以gradle格式输出)")
		fmt.Println("  mvnreps g:a,10 gradle.kts     (查看最近10个版本并以Kotlin DSL格式输出)")
		fmt.Println("\n格式选项:")
		fmt.Println("  default   - g:a:v (默认)")
		fmt.Println("  pom       - Maven XML 格式")
		fmt.Println("  gradle    - Gradle 格式")
		fmt.Println("  gradle.kts - Kotlin DSL 格式")
		fmt.Println("  format    - 表格格式")
		os.Exit(0)
	}

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("❌ 错误: 请提供搜索关键词或 G:A 坐标")
		flag.Usage()
	}

	config := parseConfig(args[0], getFormat(args))
	runSearch(config)
}

func getFormat(args []string) string {
	if len(args) > 1 {
		return strings.ToLower(args[1])
	}
	return "default"
}

func parseConfig(firstParam, formatParam string) Config {
	config := Config{
		Format: formatParam,
		Rows:   5,
	}

	if strings.Contains(firstParam, ":") {
		config.IsPrecise = true
		parts := strings.Split(firstParam, ",")
		coords := parts[0]
		if len(parts) > 1 {
			if rows, err := strconv.Atoi(parts[1]); err == nil {
				config.Rows = rows
			}
		}

		coordParts := strings.Split(coords, ":")
		config.GroupId = strings.TrimSpace(coordParts[0])
		config.ArtifactId = strings.TrimSpace(coordParts[1])
		config.Query = fmt.Sprintf(`g:"%s" AND a:"%s"`, config.GroupId, config.ArtifactId)
	} else {
		parts := strings.Split(firstParam, ",")
		config.Query = strings.TrimSpace(parts[0])
		if len(parts) > 1 {
			if rows, err := strconv.Atoi(parts[1]); err == nil {
				config.Rows = rows
			}
		}
	}

	return config
}

func runSearch(config Config) {
	baseURL := "https://search.maven.org/solrsearch/select"
	params := url.Values{}
	params.Set("q", config.Query)
	params.Set("rows", strconv.Itoa(config.Rows))
	params.Set("wt", "json")

	if config.IsPrecise {
		params.Set("core", "gav")
		params.Set("sort", "timestamp desc")
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Printf("❌ 发生错误: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ HTTP 错误: %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ 读取响应失败: %v\n", err)
		return
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("❌ 解析 JSON 失败: %v\n", err)
		return
	}

	if len(result.Response.Docs) == 0 {
		fmt.Println("\n⚠️ 未找到匹配的结果。")
		return
	}

	modeText := "库搜索"
	if config.IsPrecise {
		modeText = "历史版本"
	}
	queryText := config.Query
	if config.IsPrecise {
		queryText = fmt.Sprintf("%s:%s", config.GroupId, config.ArtifactId)
	}
	fmt.Printf("\n🔍 %s结果 (\"%s\", 前 %d 条):\n\n", modeText, queryText, config.Rows)

	for i, doc := range result.Response.Docs {
		g := doc.G
		a := doc.A
		v := doc.V
		if !config.IsPrecise && doc.LatestVersion != "" {
			v = doc.LatestVersion
		}

		switch config.Format {
		case "pom":
			fmt.Printf(`<dependency>
  <groupId>%s</groupId>
  <artifactId>%s</artifactId>
  <version>%s</version>
</dependency>
`, g, a, v)
		case "gradle":
			fmt.Printf("implementation '%s:%s:%s'\n", g, a, v)
		case "gradle.kts":
			fmt.Printf("implementation(\"%s:%s:%s\")\n", g, a, v)
		case "format":
			if i == 0 {
				fmt.Printf("%-40s | %-30s | %s\n", "GroupId", "ArtifactId", "Version")
				fmt.Println(strings.Repeat("-", 85))
			}
			fmt.Printf("%-40s | %-30s | %s\n", g, a, v)
		default:
			fmt.Printf("%s:%s:%s\n", g, a, v)
		}
	}

	fmt.Println("\n--------------------------------------------------")
}
