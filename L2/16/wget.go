package main

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	defaultDepth   = 2
	defaultWorkers = 8
	defaultTimeout = time.Second * 30
)

type Task struct {
	URL   string
	Depth int
}

type Set struct {
	mu  sync.Mutex
	set map[string]struct{}
}

func NewSet() *Set {
	return &Set{set: make(map[string]struct{})}
}

func (v *Set) Contains(u string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	_, ok := v.set[u]
	return ok
}

func (v *Set) Add(u string) {
	v.mu.Lock()
	v.set[u] = struct{}{}
	v.mu.Unlock()
}

type Downloader struct {
	opts    Options
	tasks   chan Task
	visited *Set
	wg      sync.WaitGroup
	client  *http.Client
}

func NewDownloader(opts Options) *Downloader {
	return &Downloader{
		opts:    opts,
		tasks:   make(chan Task, opts.workers*10),
		visited: NewSet(),
		client: &http.Client{
			Timeout: opts.timeout,
		},
	}
}

func (dl *Downloader) Start() {
	ts := time.Now()

	slog.Info(
		"start downloader",
		slog.Group("opts",
			slog.String("baseURL", dl.opts.baseURL),
			slog.Int("maxDepth", dl.opts.maxDepth),
			slog.Int("workers", dl.opts.workers),
			slog.Duration("timeout", dl.opts.timeout),
		),
	)

	dl.wg.Add(1)
	go func() {
		dl.tasks <- Task{URL: dl.opts.baseURL, Depth: 0}
	}()

	for range dl.opts.workers {
		go dl.worker()
	}

	dl.wg.Wait()
	close(dl.tasks)

	slog.Info(
		"stop downloader",
		slog.Duration("elapsed", time.Since(ts)),
	)
}

func (dl *Downloader) worker() {
	for task := range dl.tasks {
		dl.processTask(task)
	}
}

func (dl *Downloader) processTask(task Task) {
	defer dl.wg.Done()

	clean := dl.cleanLink(task.URL)
	if dl.visited.Contains(clean) {
		return
	}

	dl.visited.Add(clean)

	slog.Info(
		"started task",
		"url", task.URL,
		"depth", task.Depth,
	)

	resp, err := dl.client.Get(task.URL)
	if err != nil {
		slog.Error(
			"failed to send http request",
			"url", task.URL,
			"err", err,
		)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(
			"failed to read page content",
			"url", task.URL,
			"err", err,
		)
		return
	}

	var (
		urls []string

		ct      = resp.Header.Get("Content-Type")
		updated = body
	)

	if strings.HasPrefix(ct, "text/html") {
		baseURL, _ := url.Parse(task.URL)

		processed, err := dl.processHTML(baseURL, body)
		if err != nil {
			slog.Error(
				"failed to process page",
				"url", task.URL,
				"err", err,
			)
		} else {
			updated = processed.Body
			urls = processed.Links
		}
	}

	dl.saveFile(task.URL, updated)

	if task.Depth < dl.opts.maxDepth {
		for _, url := range urls {
			dl.wg.Add(1) // Увеличиваем счетчик синхронно
			go func() {
				dl.tasks <- Task{URL: url, Depth: task.Depth + 1}
			}()
		}
	}
}

func (dl *Downloader) cleanLink(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}

	u.Fragment = ""

	return u.String()
}

type ProcessedHTML struct {
	Body  []byte
	Links []string
}

func (dl *Downloader) processHTML(baseURL *url.URL, body []byte) (*ProcessedHTML, error) {
	r := bytes.NewReader(body)

	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var (
		links []string
		walk  func(*html.Node)
		buf   bytes.Buffer
	)

	walk = func(n *html.Node) {
		if n == nil {
			return
		}

		var attrName string
		switch n.DataAtom {
		case atom.A, atom.Link, atom.Area:
			attrName = "href"
		case atom.Img, atom.Script, atom.Iframe:
			attrName = "src"
		case atom.Form:
			attrName = "action"
		}

		if attrName != "" {
			for i, a := range n.Attr {
				if a.Key == attrName && a.Val != "" {
					abs, err := url.Parse(a.Val)
					if err == nil {
						resolved := baseURL.ResolveReference(abs)

						if resolved.Scheme == "http" || resolved.Scheme == "https" {
							links = append(links, dl.cleanLink(resolved.String()))
						}
					}

					rel := dl.toRelativeLink(baseURL, a.Val)
					if rel != "" {
						n.Attr[i].Val = rel
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(root)

	err = html.Render(&buf, root)
	if err != nil {
		return nil, err
	}

	return &ProcessedHTML{
		Body:  buf.Bytes(),
		Links: links,
	}, nil
}

func (dl *Downloader) toRelativeLink(base *url.URL, link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return link
	}

	resolved := base.ResolveReference(u)
	if resolved.Scheme != "http" && resolved.Scheme != "https" {
		return link
	}

	target := dl.toLocalPath(resolved)
	dir := filepath.Dir(dl.toLocalPath(base))

	rel, err := filepath.Rel(dir, target)
	if err != nil {
		return link
	}

	return filepath.ToSlash(rel)
}

func (dl *Downloader) toLocalPath(u *url.URL) string {
	p := u.EscapedPath()

	if p == "" || p == "/" {
		p = "/index.html"
	}

	filename := filepath.Join(u.Host, p)
	if filepath.Ext(filename) == "" {
		filename += ".html"
	}

	return filename
}

func (dl *Downloader) saveFile(URL string, body []byte) {
	u, err := url.Parse(URL)
	if err != nil {
		return
	}

	filename := dl.toLocalPath(u)
	dir := filepath.Dir(filename)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		slog.Error(
			"failed to create dir",
			"dir", dir,
			"err", err,
		)
		return
	}

	err = os.WriteFile(filename, body, 0644)
	if err != nil {
		slog.Error(
			"failed to save file",
			"dir", dir,
			"file", filename,
			"err", err,
		)
		return
	}

	slog.Info(
		"saved file",
		"url", URL,
		"file", filename,
	)
}

type Options struct {
	baseURL  string
	workers  int
	maxDepth int
	timeout  time.Duration
}

func parseOptions() (Options, error) {
	var opts Options

	flag.IntVar(&opts.maxDepth, "depth", defaultDepth, "Max depth of recursion")
	flag.IntVar(&opts.workers, "workers", defaultWorkers, "Number of parallel workers")
	flag.DurationVar(&opts.timeout, "timeout", defaultTimeout, "Timeout for HTTP request")

	flag.Parse()

	_, err := url.Parse(flag.Arg(0))
	if err != nil {
		return Options{}, errors.New("invalid base URL")
	}

	opts.baseURL = flag.Arg(0)

	return opts, nil
}

func main() {
	opts, err := parseOptions()
	if err != nil {
		log.Fatalf("failed to parse options: %v\n", err)
	}

	dl := NewDownloader(opts)

	dl.Start()
}
