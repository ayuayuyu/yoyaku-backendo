package main

import (
	"log"
	"os"

	"yoyaku/auth"
	"yoyaku/db"
	"yoyaku/handler"
	"yoyaku/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

func main() {

	//データベースとの接続
	sqlDB, err := utils.NewDBConnection()
	if err != nil {
		log.Fatalf("データベースに接続できませんでした: %v", err)
	}
	defer sqlDB.Close()
	// sql.DBからsqlcのクエリオブジェクトを生成
	queries := db.New(sqlDB)

	// 1. OAuth設定の初期化
	if err := auth.Setup(); err != nil {
		log.Fatalf("OAuth設定の初期化に失敗しました: %v", err)
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatalf("環境変数 SECRET_KEY が設定されていません")
	}

	// セッション情報を保存するためのストア (キーは秘密の値にしてください)
	var store = sessions.NewCookieStore([]byte(secretKey))

	// Ginのルーティング
	r := gin.Default()

	// CORS設定
	// フロントエンドのURL(http://localhost:3000)からのリクエストを許可
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true // Cookieの送受信を許可
	r.Use(cors.New(config))

	// セッションミドルウェア
	// 全てのリクエストでセッションストアを利用可能にする
	r.Use(func(c *gin.Context) {
		c.Set("session_store", store)
		c.Next()
	})

	// r.GET("/", func(c *gin.Context) {
	// 	// GORMの`db.Create`の代わりに、sqlcが生成したメソッドを使います
	// 	// 例として、authorsテーブルに新しいレコードを作成します
	// 	createdAuthor, err := queries.CreateAuthor(context.Background(), db.CreateAuthorParams{
	// 		Name: "Gin Framework",
	// 		Bio:  sql.NullString{String: "A web framework written in Go.", Valid: true},
	// 	})
	// 	if err != nil {
	// 		// エラーが発生した場合は、サーバーエラーを返します
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create author"})
	// 		return
	// 	}

	// 	// 成功した場合は、作成されたデータを含むJSONを返します
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message":        "Hello Database, author created!",
	// 		"created_author": createdAuthor,
	// 	})
	// })

	// 3. ルーティングの設定
	r.GET("/login", handler.HandleGoogleLogin)
	r.GET("/callback", func(c *gin.Context) {
		handler.HandleGoogleCallback(c, queries)
	})

	// フロントエンドがユーザー情報を確認するためのAPIエンドポイント
	api := r.Group("/api")
	{
		api.GET("/me", handler.HandleGetMe)
		api.POST("/logout", handler.HandleLogout)
	}

	// 4. サーバーの起動
	log.Println("Started server on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
