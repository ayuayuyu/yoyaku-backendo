package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"yoyaku/auth"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// Googleから返ってくるユーザー情報の構造体
type GoogleUserInfo struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`        // フルネーム
	GivenName  string `json:"given_name"`  // 名
	FamilyName string `json:"family_name"` // 姓
	Picture    string `json:"picture"`     // プロフィール画像のURL
}

// (変更) HandleGoogleCallbackはセッションに情報を保存し、フロントエンドにリダイレクトする
func HandleGoogleCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	// ユーザー情報を取得
	content, err := auth.GetUserInfo(state, code)
	if err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/login?error=true")
		return
	}

	//デバッグ用：Googleから返ってきたJSONをそのまま出力
	fmt.Println("Response from Google:", string(content))
	// ユーザー情報をパース
	var userInfo GoogleUserInfo
	if err := json.Unmarshal(content, &userInfo); err != nil {
		log.Println("JSON Unmarshal error:", err)
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/login?error=true")
		return
	}

	// セッションを取得
	store := c.MustGet("session_store").(*sessions.CookieStore)
	session, _ := store.Get(c.Request, "session-name")

	// セッションにユーザー情報を保存
	session.Values["user_id"] = userInfo.ID
	session.Values["user_email"] = userInfo.Email
	session.Values["user_name"] = userInfo.Name
	session.Values["user_picture"] = userInfo.Picture

	fmt.Printf("user_id: %s, user_email: %s, user_name: %s, user_picture: %s", session.Values["user_id"], session.Values["user_email"], session.Values["user_name"], session.Values["user_picture"])
	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Println("Session save error:", err)
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/login?error=true")
		return
	}

	// ログイン成功後、フロントエンドのトップページにリダイレクト
	c.Redirect(http.StatusPermanentRedirect, "http://localhost:3000/")
}

// (新規) HandleGetMeは現在のユーザー情報を返すAPI
func HandleGetMe(c *gin.Context) {
	store := c.MustGet("session_store").(*sessions.CookieStore)
	session, _ := store.Get(c.Request, "session-name")

	// セッションにユーザー情報がなければ未認証エラー
	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// ユーザー情報をJSONで返す
	c.JSON(http.StatusOK, gin.H{
		"id":      userID,
		"email":   session.Values["user_email"],
		"name":    session.Values["user_name"],
		"picture": session.Values["user_picture"],
	})
}

// (新規) HandleLogoutはセッションを破棄してログアウトさせる
func HandleLogout(c *gin.Context) {
	store := c.MustGet("session_store").(*sessions.CookieStore)
	session, _ := store.Get(c.Request, "session-name")

	// セッション情報をクリア
	session.Values["user_id"] = ""
	session.Options.MaxAge = -1 // Cookieを即時削除

	if err := session.Save(c.Request, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// (元のHandleGoogleLoginは変更なしでOK)
func HandleGoogleLogin(c *gin.Context) {
	url := auth.GoogleOauthConfig.AuthCodeURL(auth.OauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}
