package traefik

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/freqkflag/joeyctl/internal/run"
)

const defaultDynamicDir = "/etc/traefik/dynamic"

type CatchallMode string

const (
	ModeWildcard CatchallMode = "wildcard"
	ModeAnyhost  CatchallMode = "anyhost"
)

type CatchallOptions struct {
	Domain     string
	DynamicDir string
	EntryPoint string
	BackendURL string
	Priority   int
	Mode       CatchallMode
}

func WriteCatchall(r run.Runner, o CatchallOptions) (string, error) {
	o = normalizeOptions(o)

	if err := validateOptions(o); err != nil {
		return "", err
	}

	filename, payload, err := buildYAML(o)
	if err != nil {
		return "", err
	}

	if r.DryRun {
		return fmt.Sprintf("[dry-run] would write %s\n%s\n[dry-run] would reload traefik", filename, payload), nil
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return "", fmt.Errorf("ensure dynamic dir: %w", err)
	}

	tmp := filename + ".tmp"
	if err := os.WriteFile(tmp, []byte(payload), 0o600); err != nil {
		return "", fmt.Errorf("write temp file: %w", err)
	}

	if err := os.Rename(tmp, filename); err != nil {
		return "", fmt.Errorf("activate config: %w", err)
	}

	if err := reloadTraefik(r); err != nil {
		return "", err
	}

	return "Wrote: " + filename, nil
}

func RemoveCatchalls(r run.Runner, dynamicDir, domain string) (string, error) {
	dir := dynamicDir
	if dir == "" {
		dir = defaultDynamicDir
	}
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return "", fmt.Errorf("domain is required")
	}

	paths := []string{
		filepath.Join(dir, "catchall-anyhost.yml"),
		filepath.Join(dir, fmt.Sprintf("catchall-%s.yml", domain)),
	}

	if r.DryRun {
		return fmt.Sprintf("[dry-run] would remove:\n%s\n[dry-run] would reload traefik", strings.Join(paths, "\n")), nil
	}

	for _, p := range paths {
		err := os.Remove(p)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("remove %s: %w", p, err)
		}
	}

	if err := reloadTraefik(r); err != nil {
		return "", err
	}

	return "Removed catchall definitions (if present).", nil
}

func normalizeOptions(o CatchallOptions) CatchallOptions {
	if o.DynamicDir == "" {
		o.DynamicDir = defaultDynamicDir
	}
	if o.EntryPoint == "" {
		o.EntryPoint = "websecure"
	}
	if o.Priority == 0 {
		o.Priority = 1
	}
	o.Domain = strings.TrimSpace(o.Domain)
	o.BackendURL = strings.TrimSpace(o.BackendURL)
	return o
}

func validateOptions(o CatchallOptions) error {
	if o.BackendURL == "" {
		return fmt.Errorf("backend URL is required")
	}
	if !strings.HasPrefix(o.BackendURL, "http://") && !strings.HasPrefix(o.BackendURL, "https://") {
		return fmt.Errorf("backend must start with http:// or https://")
	}

	switch o.Mode {
	case ModeWildcard:
		if o.Domain == "" {
			return fmt.Errorf("domain is required for wildcard catchall")
		}
	case ModeAnyhost:
	default:
		return fmt.Errorf("unknown catchall mode: %s", o.Mode)
	}

	return nil
}

func buildYAML(o CatchallOptions) (string, string, error) {
	switch o.Mode {
	case ModeWildcard:
		router := routerName(o.Domain)
		service := router + "-svc"
		filename := filepath.Join(o.DynamicDir, fmt.Sprintf("catchall-%s.yml", o.Domain))
		rule := fmt.Sprintf("Host(`%s`) || HostRegexp(`{subdomain:.+}.%s`)", o.Domain, o.Domain)
		content := fmt.Sprintf(`http:
  routers:
    %s:
      rule: %q
      entryPoints:
        - %s
      service: %s
      priority: %d
      tls: {}
  services:
    %s:
      loadBalancer:
        servers:
          - url: "%s"
`, router, rule, o.EntryPoint, service, o.Priority, service, o.BackendURL)
		return filename, content, nil
	case ModeAnyhost:
		filename := filepath.Join(o.DynamicDir, "catchall-anyhost.yml")
		rule := "HostRegexp(`{host:.+}`)"
		content := fmt.Sprintf(`http:
  routers:
    catchall-anyhost:
      rule: %q
      entryPoints:
        - %s
      service: catchall-anyhost-svc
      priority: 0
      tls: {}
  services:
    catchall-anyhost-svc:
      loadBalancer:
        servers:
          - url: "%s"
`, rule, o.EntryPoint, o.BackendURL)
		return filename, content, nil
	default:
		return "", "", fmt.Errorf("unsupported mode: %s", o.Mode)
	}
}

func routerName(domain string) string {
	domain = strings.ToLower(domain)
	return "catchall-" + strings.ReplaceAll(domain, ".", "-")
}

func reloadTraefik(r run.Runner) error {
	if r.DryRun {
		return nil
	}

	if _, _, err := r.Cmd("systemctl", "reload", "traefik"); err == nil {
		return nil
	} else {
		if _, _, restartErr := r.Cmd("systemctl", "restart", "traefik"); restartErr != nil {
			return fmt.Errorf("systemctl reload traefik failed: %v; restart also failed: %w", err, restartErr)
		}
	}

	return nil
}
