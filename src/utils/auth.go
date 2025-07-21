package utils

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// GetUserIDFromSession は、Ginのコンテキストからセッション情報を取得し、ユーザーIDを返します。
// 処理中にエラーが発生した場合は、自動的にクライアントにエラーレスポンスを返し、falseを返します。
// 成功した場合は、ユーザーID(uint64)とtrueを返します。
func GetUserIDFromSession(c *gin.Context) (uint64, bool) {
	// c.MustGetからセッションストアを取得
	store, ok := c.MustGet("session_store").(*sessions.CookieStore)
	if !ok {
		log.Println("セッションストアの取得に失敗しました")
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "サーバー内部エラーが発生しました"})
		c.Abort() // 処理を中断
		return 0, false
	}

	session, _ := store.Get(c.Request, "session-name")

	// セッションからuser_idを取得
	userIDStr, ok := session.Values["user_id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "ログイン情報が見つかりません"})
		c.Abort()
		return 0, false
	}

	// 文字列のIDをint64に変換
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "ユーザーIDが不正です"})
		c.Abort()
		return 0, false
	}

	// ユーザーIDが負でないことを確認
	if userID < 0 {
		log.Println("不正なユーザーID（負の値）:", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "不正なユーザーIDです"})
		c.Abort()
		return 0, false
	}

	// uint64にキャストして返す
	return uint64(userID), true
}
