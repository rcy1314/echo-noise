package dto

type UserInfoDto struct {
    Username     string `json:"username"`
    Password     string `json:"password"`
    AvatarURL    string `json:"avatar_url"`
    SiteName     string `json:"siteName"`
    Theme        string `json:"theme"`
    AllowSignUp  bool   `json:"allowSignUp"`
    IsAdmin      bool   `json:"is_admin"`
}

// 标签相关 DTO
type TagDto struct {
    Name  string `json:"name"`
    Count int    `json:"count"`
}

// 图片相关 DTO
type ImageDto struct {
    ID        uint      `json:"id"`
    URL       string    `json:"url"`
    CreatedAt string    `json:"created_at"`
}

// 标签查询响应
type TagsResponse struct {
    Tags []TagDto `json:"tags"`
}

// 图片查询响应
type ImagesResponse struct {
    Images []ImageDto `json:"images"`
}
