package geoip

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/comment"
)

type Resolver struct {
	city     *geoip2.Reader
	asn      *geoip2.Reader
	cityPath string
	asnPath  string
	mu       sync.Mutex
}

func NewResolver(cityPath, asnPath string) (*Resolver, error) {
	resolver := &Resolver{cityPath: cityPath, asnPath: asnPath}
	if cityPath != "" {
		reader, err := geoip2.Open(cityPath)
		if err != nil {
			return nil, err
		}
		resolver.city = reader
	}
	if asnPath != "" {
		reader, err := geoip2.Open(asnPath)
		if err != nil {
			return nil, err
		}
		resolver.asn = reader
	}
	return resolver, nil
}

func NewLazyResolver(cityPath, asnPath string) *Resolver {
	return &Resolver{cityPath: cityPath, asnPath: asnPath}
}

func (r *Resolver) Resolve(ip string) string {
	if r == nil || r.city == nil {
		r.tryOpen()
		if r.city == nil {
			return ""
		}
	}
	addr, err := netip.ParseAddr(strings.TrimSpace(ip))
	if err != nil {
		return ""
	}
	record, err := r.city.City(addr)
	if err != nil || !record.HasData() {
		return ""
	}

	var parts []string
	countryName := pickName(record.Country.Names)
	if countryName == "" && record.Country.ISOCode == "CN" {
		countryName = "中国"
	}
	if countryName != "" {
		parts = append(parts, countryName)
	}
	if len(record.Subdivisions) > 0 {
		region := pickName(record.Subdivisions[0].Names)
		if region != "" && record.Country.ISOCode == "CN" {
			region = normalizeCNRegion(region)
		}
		if region != "" {
			parts = append(parts, region)
		}
	}
	city := pickName(record.City.Names)
	if city != "" && record.Country.ISOCode == "CN" {
		city = normalizeCNCity(city)
	}
	if city != "" {
		parts = append(parts, city)
	}

	isp := r.resolveASN(addr)
	if isp != "" {
		parts = append(parts, isp)
	}
	return strings.Join(parts, "")
}

var _ comment.GeoIPResolver = (*Resolver)(nil)

func (r *Resolver) tryOpen() {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.city == nil && r.cityPath != "" {
		reader, err := geoip2.Open(r.cityPath)
		if err == nil {
			r.city = reader
		}
	}
	if r.asn == nil && r.asnPath != "" {
		reader, err := geoip2.Open(r.asnPath)
		if err == nil {
			r.asn = reader
		}
	}
}

func EnsureDatabaseAsync(ctx context.Context, dbPath, url string, logger func(string, ...any)) {
	if dbPath == "" || url == "" {
		return
	}
	if _, err := os.Stat(dbPath); err == nil {
		return
	}
	if !markDownloadOnce(dbPath) {
		return
	}
	go func() {
		if err := downloadDatabase(ctx, dbPath, url); err != nil {
			if logger != nil {
				logger("geoip db download failed: %v", err)
			}
		}
	}()
}

func EnsureDatabasesAsync(ctx context.Context, cityPath, cityURL, asnPath, asnURL string, logger func(string, ...any)) {
	EnsureDatabaseAsync(ctx, cityPath, cityURL, logger)
	EnsureDatabaseAsync(ctx, asnPath, asnURL, logger)
}

func downloadDatabase(ctx context.Context, dbPath, url string) error {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("geoip download status: " + resp.Status)
	}

	tmpPath := dbPath + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer out.Close()

	var reader io.Reader = resp.Body
	if strings.HasSuffix(strings.ToLower(url), ".gz") {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
		defer gz.Close()
		reader = gz
	}

	if _, err := io.Copy(out, reader); err != nil {
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, dbPath)
}

var downloadGuard sync.Map

func markDownloadOnce(dbPath string) bool {
	_, loaded := downloadGuard.LoadOrStore(dbPath, struct{}{})
	return !loaded
}

func pickName(names geoip2.Names) string {
	if names.SimplifiedChinese != "" {
		return names.SimplifiedChinese
	}
	if names.English != "" {
		return names.English
	}
	return ""
}

func normalizeCNRegion(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	if strings.HasSuffix(name, "省") || strings.HasSuffix(name, "市") || strings.HasSuffix(name, "自治区") {
		return name
	}
	return name + "省"
}

func normalizeCNCity(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	if strings.HasSuffix(name, "市") {
		return name
	}
	return name + "市"
}

func (r *Resolver) resolveASN(addr netip.Addr) string {
	if r == nil {
		return ""
	}
	if r.asn == nil {
		r.tryOpen()
		if r.asn == nil {
			return ""
		}
	}
	record, err := r.asn.ASN(addr)
	if err != nil {
		return ""
	}
	return normalizeISP(strings.TrimSpace(record.AutonomousSystemOrganization))
}

func normalizeISP(org string) string {
	if org == "" {
		return ""
	}
	if containsNonASCII(org) {
		return org
	}
	lower := strings.ToLower(org)
	switch {
	case strings.Contains(lower, "china telecom") || strings.Contains(lower, "chinanet"):
		return "电信"
	case strings.Contains(lower, "china unicom"):
		return "联通"
	case strings.Contains(lower, "china mobile"):
		return "移动"
	case strings.Contains(lower, "china broadcasting") || strings.Contains(lower, "china broadnet"):
		return "广电"
	case strings.Contains(lower, "cernet"):
		return "教育网"
	case strings.Contains(lower, "cstnet"):
		return "科技网"
	default:
		return org
	}
}

func containsNonASCII(input string) bool {
	for _, r := range input {
		if r > 127 {
			return true
		}
	}
	return false
}
