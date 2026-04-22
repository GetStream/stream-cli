package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"golang.org/x/term"
)

// UploadToS3 uploads a local file to the given S3 presigned URL via HTTP PUT.
// When stderr is a TTY, a progress bar is rendered so large uploads are not opaque.
func UploadToS3(ctx context.Context, filename, url string) error {
	data, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer data.Close()

	stat, err := data.Stat()
	if err != nil {
		return err
	}

	size := stat.Size()
	body := io.Reader(data)
	var pr *progressReader
	if term.IsTerminal(int(os.Stderr.Fd())) {
		pr = newProgressReader(data, size, filepath.Base(filename), os.Stderr)
		body = pr
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = size

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if pr != nil {
			pr.abort()
		}
		return err
	}
	defer resp.Body.Close()

	if pr != nil {
		pr.finish()
	}
	return nil
}

// progressReader wraps an io.Reader and renders a single-line progress bar
// (carriage-return overwritten) to the given writer. Updates are throttled.
type progressReader struct {
	r         io.Reader
	total     int64
	read      int64
	name      string
	out       io.Writer
	start     time.Time
	lastPrint time.Time
	barWidth  int
	done      int32
}

func newProgressReader(r io.Reader, total int64, name string, out io.Writer) *progressReader {
	return &progressReader{
		r:        r,
		total:    total,
		name:     name,
		out:      out,
		start:    time.Now(),
		barWidth: 30,
	}
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.r.Read(b)
	if n > 0 {
		atomic.AddInt64(&p.read, int64(n))
		p.maybeRender(false)
	}
	return n, err
}

func (p *progressReader) maybeRender(force bool) {
	now := time.Now()
	if !force && now.Sub(p.lastPrint) < 100*time.Millisecond {
		return
	}
	p.lastPrint = now
	p.render()
}

func (p *progressReader) render() {
	read := atomic.LoadInt64(&p.read)
	elapsed := time.Since(p.start).Seconds()
	var rate float64
	if elapsed > 0 {
		rate = float64(read) / elapsed
	}

	var pct float64
	var bar string
	var eta string
	if p.total > 0 {
		pct = float64(read) / float64(p.total) * 100
		filled := min(int(float64(p.barWidth)*float64(read)/float64(p.total)), p.barWidth)
		b := make([]byte, p.barWidth)
		for i := 0; i < p.barWidth; i++ {
			if i < filled {
				b[i] = '='
			} else {
				b[i] = ' '
			}
		}
		bar = string(b)
		if rate > 0 && read < p.total {
			remaining := time.Duration(float64(p.total-read)/rate) * time.Second
			eta = " ETA " + formatDuration(remaining)
		}
		fmt.Fprintf(p.out, "\rUploading %s [%s] %5.1f%% %s/%s %s/s%s",
			p.name, bar, pct, humanBytes(read), humanBytes(p.total), humanBytes(int64(rate)), eta)
	} else {
		fmt.Fprintf(p.out, "\rUploading %s %s %s/s", p.name, humanBytes(read), humanBytes(int64(rate)))
	}
}

func (p *progressReader) finish() {
	if !atomic.CompareAndSwapInt32(&p.done, 0, 1) {
		return
	}
	p.render()
	fmt.Fprintln(p.out)
}

func (p *progressReader) abort() {
	if !atomic.CompareAndSwapInt32(&p.done, 0, 1) {
		return
	}
	fmt.Fprintln(p.out)
}

func humanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for n/div >= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return "<1s"
	}
	d = d.Round(time.Second)
	h := int(d / time.Hour)
	m := int((d % time.Hour) / time.Minute)
	s := int((d % time.Minute) / time.Second)
	if h > 0 {
		return fmt.Sprintf("%dh%02dm%02ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%02ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
