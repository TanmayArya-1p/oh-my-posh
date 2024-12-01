package upgrade

import (
	"fmt"
	"os"

	"github.com/jandedobbeleer/oh-my-posh/src/build"
	"github.com/jandedobbeleer/oh-my-posh/src/cache"
	"github.com/jandedobbeleer/oh-my-posh/src/log"
	"github.com/jandedobbeleer/oh-my-posh/src/runtime/http"
)

const (
	RELEASEURL = "https://api.github.com/repos/jandedobbeleer/oh-my-posh/releases/latest"
	CACHEKEY   = "upgrade_check"

	upgradeNotice = `
A new release of Oh My Posh is available: %s → %s
To upgrade, run: 'oh-my-posh upgrade%s'

To enable automated upgrades, run: 'oh-my-posh enable upgrade'.
`
)

// Returns the upgrade notice if a new version is available
// that should be displayed to the user.
//
// The upgrade check is only performed every other week.
func (cfg *Config) Notice() (string, bool) {
	// never validate when we install using the Windows Store
	if os.Getenv("POSH_INSTALLER") == "ws" {
		log.Debug("skipping upgrade check because we are using the Windows Store")
		return "", false
	}

	// do not check when last validation was < 1 week ago
	if _, OK := cfg.Cache.Get(CACHEKEY); OK && !cfg.Force {
		return "", false
	}

	if !http.IsConnected() {
		return "", false
	}

	latest, err := cfg.Latest()
	if err != nil {
		return "", false
	}

	cfg.Cache.Set(CACHEKEY, latest, cache.ONEWEEK)

	version := fmt.Sprintf("v%s", build.Version)
	if latest == version {
		return "", false
	}

	var forceUpdate string
	if IsMajorUpgrade(version, latest) {
		forceUpdate = " --force"
	}

	return fmt.Sprintf(upgradeNotice, version, latest, forceUpdate), true
}
