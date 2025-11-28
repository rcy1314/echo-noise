// setting.go
package dto

type SettingDto struct {
    AllowRegistration bool `json:"allowRegistration"`
    FrontendSettings  struct {
        SiteTitle           string   `json:"siteTitle"`
        SubtitleText        string   `json:"subtitleText"`
        AvatarURL           string   `json:"avatarURL"`
        Username            string   `json:"username"`
        Description         string   `json:"description"`
        Backgrounds         []string `json:"backgrounds"`
        CardFooterTitle     string   `json:"cardFooterTitle"`
        CardFooterLink      string   `json:"cardFooterLink"`
        PageFooterHTML      string   `json:"pageFooterHTML"`
        RSSTitle            string   `json:"rssTitle"`
        RSSDescription      string   `json:"rssDescription"`
        RSSAuthorName       string   `json:"rssAuthorName"`
        RSSFaviconURL       string   `json:"rssFaviconURL"`
        WalineServerURL     string   `json:"walineServerURL"`
        EnableGithubCard    *bool    `json:"enableGithubCard"`
        PwaEnabled          *bool    `json:"pwaEnabled"`
        PwaTitle            *string  `json:"pwaTitle"`
        PwaDescription      *string  `json:"pwaDescription"`
        PwaIconURL          *string  `json:"pwaIconURL"`
        DefaultContentTheme *string  `json:"defaultContentTheme"`
        AnnouncementText    string   `json:"announcementText"`
        AnnouncementEnabled *bool    `json:"announcementEnabled"`
    } `json:"frontendSettings"`
    SmtpEnabled    *bool   `json:"smtpEnabled"`
    SmtpDriver     *string `json:"smtpDriver"`
    SmtpHost       *string `json:"smtpHost"`
    SmtpPort       *int    `json:"smtpPort"`
    SmtpUser       *string `json:"smtpUser"`
    SmtpPass       *string `json:"smtpPass"`
    SmtpFrom       *string `json:"smtpFrom"`
    SmtpEncryption *string `json:"smtpEncryption"`
    SmtpTLS        *bool   `json:"smtpTLS"`
}
