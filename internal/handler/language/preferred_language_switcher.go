package language

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// PreferredLanguageSwitcher handles preferred language switching
type PreferredLanguageSwitcher struct {
	configService *service.ConfigService
	featureService *service.FeatureService
	i18nService   *service.I18nService
	cookieService *service.CookieService
	logger        *service.Logger
}

// NewPreferredLanguageSwitcher creates a new instance of PreferredLanguageSwitcher
func NewPreferredLanguageSwitcher(
	configService *service.ConfigService,
	featureService *service.FeatureService,
	i18nService *service.I18nService,
	cookieService *service.CookieService,
	logger *service.Logger,
) *PreferredLanguageSwitcher {
	return &PreferredLanguageSwitcher{
		configService:  configService,
		featureService: featureService,
		i18nService:    i18nService,
		cookieService:  cookieService,
		logger:         logger,
	}
}

// InitPreferredLanguage initializes the preferred language
func (p *PreferredLanguageSwitcher) InitPreferredLanguage(c *gin.Context) error {
	// Check if preferred language cookie is disabled
	disabled, err := p.featureService.IsEnabled(c, "disable_preferred_language_cookie")
	if err != nil {
		return err
	}

	if disabled {
		return nil
	}

	// Get preferred language
	preferredLanguage, err := p.GetPreferredLanguage(c)
	if err != nil {
		return err
	}

	// Set preferred language cookie
	return p.cookieService.Set(c, "preferred_language", preferredLanguage)
}

// GetPreferredLanguage gets the preferred language
func (p *PreferredLanguageSwitcher) GetPreferredLanguage(c *gin.Context) (string, error) {
	// Get preferred language from cookie
	cookieLanguage, err := p.cookieService.Get(c, "preferred_language")
	if err == nil && cookieLanguage != "" {
		// Check if cookie language is available
		if p.isAvailableLocale(cookieLanguage) {
			return cookieLanguage, nil
		}
	}

	// Get language from params
	paramLanguage := p.getLanguageFromParams(c)
	if paramLanguage != "" && p.isAvailableLocale(paramLanguage) {
		return paramLanguage, nil
	}

	// Get browser languages
	browserLanguages := p.getBrowserLanguages(c)
	for _, lang := range browserLanguages {
		if p.isAvailableLocale(lang) {
			return lang, nil
		}
	}

	// Get default preferred language
	return p.getDefaultPreferredLanguage(c)
}

// SelectableLanguage selects a language from the available options
func (p *PreferredLanguageSwitcher) SelectableLanguage(languageOptions []string) string {
	// Get ordered selectable locales codes
	orderedSelectableLocalesCodes, err := p.getOrderedSelectableLocalesCodes(c)
	if err != nil {
		return ""
	}

	// Find first language that is in ordered selectable locales codes
	for _, lang := range languageOptions {
		for _, code := range orderedSelectableLocalesCodes {
			if lang == code {
				return lang
			}
		}
	}

	return ""
}

// GetOrderedSelectableLocalesCodes gets the ordered selectable locales codes
func (p *PreferredLanguageSwitcher) getOrderedSelectableLocalesCodes(c *gin.Context) ([]string, error) {
	// Get ordered selectable locales
	orderedSelectableLocales, err := p.i18nService.GetOrderedSelectableLocales(c)
	if err != nil {
		return nil, err
	}

	// Extract codes
	codes := make([]string, len(orderedSelectableLocales))
	for i, locale := range orderedSelectableLocales {
		codes[i] = locale.Value
	}

	return codes, nil
}

// GetBrowserLanguages gets the browser languages
func (p *PreferredLanguageSwitcher) getBrowserLanguages(c *gin.Context) []string {
	// Get accept language header
	acceptLanguage := c.GetHeader("Accept-Language")
	if acceptLanguage == "" {
		return []string{}
	}

	// Format accept language header
	formattedAcceptLanguage := strings.ReplaceAll(acceptLanguage, "-", "_")

	// Split accept language header
	parts := strings.Split(formattedAcceptLanguage, ",")
	languages := make([]string, 0, len(parts))

	for _, part := range parts {
		// Remove quality value
		lang := strings.Split(part, ";")[0]
		lang = strings.TrimSpace(lang)

		// Skip empty languages
		if lang == "" {
			continue
		}

		languages = append(languages, lang)
	}

	return languages
}

// GetLanguageFromParams gets the language from params
func (p *PreferredLanguageSwitcher) getLanguageFromParams(c *gin.Context) string {
	// This is a placeholder for EE implementation
	return ""
}

// IsAvailableLocale checks if a locale is available
func (p *PreferredLanguageSwitcher) isAvailableLocale(locale string) bool {
	// Get available locales
	availableLocales, err := p.i18nService.GetAvailableLocales()
	if err != nil {
		return false
	}

	// Check if locale is available
	for _, availableLocale := range availableLocales {
		if locale == availableLocale {
			return true
		}
	}

	return false
}

// GetDefaultPreferredLanguage gets the default preferred language
func (p *PreferredLanguageSwitcher) getDefaultPreferredLanguage(c *gin.Context) (string, error) {
	// Get current settings
	settings, err := p.configService.GetCurrentSettings(c)
	if err != nil {
		return "", err
	}

	return settings.DefaultPreferredLanguage, nil
}
